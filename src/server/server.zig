const Data = @import("../data.zig");

const std = @import("std");
const Server = std.net.Server;
const Address = std.net.Address;
const Stream = std.net.Stream;
const Connection = std.net.Server.Connection;
const Thread = std.Thread;
const http = std.http;

data: Data = undefined,

const Self = @This();

pub const Handler = *const fn (self: *Self, req: *http.Server.Request, allocator: std.mem.Allocator) anyerror!void;

pub const GET_ROUTES = std.StaticStringMap(Handler).initComptime(.{
    .{ "/budget/", budgetGet },
});

pub const POST_ROUTES = std.StaticStringMap(Handler).initComptime(.{
    .{ "/balance/", balancePost },
});

pub fn run(ip: [4]u8, port: u16) !void {
    const address: Address = Address.initIp4(ip, port);
    var listener = try address.listen(.{ .reuse_address = true });
    std.log.info("Listening on {d}.{d}.{d}.{d}\n", .{ ip[0], ip[1], ip[2], ip[3] });
    defer listener.deinit();
    const allocator = std.heap.page_allocator;
    while (true) {
        var conn = try listener.accept();
        const requestThread: Thread = try Thread.spawn(.{ .allocator = allocator }, handleRequest, .{&conn});
        requestThread.detach();
    }
}

fn handleRequest(conn: *Connection) !void {
    defer conn.stream.close();
    var addressBuf: [64]u8 = undefined;
    var addressWriter = std.io.Writer.fixed(&addressBuf);
    try conn.address.in.format(&addressWriter);
    std.debug.print("Request from {s}\n", .{addressWriter.buffered()});

    var recBuf: [4096]u8 = undefined;
    var sendBuf: [4096]u8 = undefined;
    var reader = conn.stream.reader(recBuf[0..]);
    var writer = conn.stream.writer(sendBuf[0..]);
    var server = http.Server.init(reader.interface(), &writer.interface);

    var req = try server.receiveHead();

    const method = req.head.method;
    // const length = req.head.content_length orelse 0;
    const target = splitTarget(req.head.target);

    var self: Self = undefined;

    self.data = Data.init("log") catch {
        std.log.err("Failed to initialize data", .{});
        try notFoundHandler(&self, &req, std.heap.page_allocator);
        return;
    };

    if (method == http.Method.GET) {
        const handler = GET_ROUTES.get(target.path) orelse notFoundHandler;
        try handler(&self, &req, std.heap.page_allocator);
    } else if (method == http.Method.POST) {
        const handler = POST_ROUTES.get(target.path) orelse notFoundHandler;
        try handler(&self, &req, std.heap.page_allocator);
    } else {
        try notFoundHandler(&self, &req, std.heap.page_allocator);
    }
    self.data.write("log") catch {
        std.log.err("Failed to write data to log", .{});
    };
}

fn splitTarget(target: []const u8) struct { path: []const u8, query: []const u8 } {
    const partIndex = std.mem.indexOf(u8, target, "?") orelse target.len;
    return .{ .path = target[0..partIndex], .query = target[partIndex..] };
}

fn budgetGet(self: *Self, req: *http.Server.Request, _: std.mem.Allocator) anyerror!void {
    // Placeholder implementation for GET /budget/
    const budget = self.data.read();
    var responseBuf: [1048]u8 = undefined;
    const response = std.fmt.bufPrint(&responseBuf, "Budget is {d}\n", .{budget}) catch {
        try req.respond("Internal Server Error", .{ .status = http.Status.internal_server_error });
        return;
    };
    try req.respond(response, .{ .status = http.Status.ok });
}

fn balancePost(self: *Self, req: *http.Server.Request, _: std.mem.Allocator) anyerror!void {
    var bodyBuffer: [4096]u8 = undefined;
    var bodyReader = req.readerExpectNone(bodyBuffer[0..]);

    var fixed = std.Io.Writer.fixed(&bodyBuffer);
    const n = try bodyReader.stream(&fixed, .limited(bodyBuffer.len));

    if (n == 0) {
        try req.respond("Empty body", .{ .status = http.Status.bad_request });
        return;
    }

    const newBalance = std.fmt.parseFloat(f32, bodyBuffer[0..n]) catch {
        try req.respond("Invalid balance value", .{ .status = http.Status.bad_request });
        return;
    };

    _ = self.data.enter(newBalance) catch {
        std.log.err("Failed to update balance", .{});
        try req.respond("Internal Server Error", .{ .status = http.Status.internal_server_error });
        return;
    };
    return budgetGet(self, req, std.heap.page_allocator);
}

fn notFoundHandler(_: *Self, req: *http.Server.Request, _: std.mem.Allocator) anyerror!void {
    try req.respond("404 Not Found", .{ .status = http.Status.not_found });
}

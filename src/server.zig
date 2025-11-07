const std = @import("std");
const Server = std.net.Server;
const Address = std.net.Address;
const Stream = std.net.Stream;
const Connection = std.net.Server.Connection;
const Thread = std.Thread;
const http = std.http;

const Self = @This();

pub const Handler = *const fn (req: *http.Request, allocator: std.mem.Allocator) anyerror!void;

pub const GET_ROUTES = std.StaticStringMap(Handler).initComptime(.{
    .{"/budget/", budgetGet},
});

pub const POST_ROUTES = std.StaticStringMap(Handler).initComptime(.{
    .{"/balance/", balancePost},
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
    const length = req.head.content_length orelse 0;
    const target = req.head.target;

    var bodyBuffer: [4096]u8 = undefined;
    var bodyReader = req.readerExpectNone(bodyBuffer[0..]);

    var fixed = std.Io.Writer.fixed(&bodyBuffer);
    const n = try bodyReader.stream(&fixed, .limited(bodyBuffer.len));

    std.debug.print("Received request: '{s}' with body length: {d}\n", .{bodyBuffer[0..n], n});

    try req.respond(bodyBuffer[0..n], .{ .status = http.Status.ok, .keep_alive = true, .extra_headers = &.{.{ .name = "content-type", .value = "text/plain; charset=utf-8" }} });
}

fn splitTarget(target: []const u8) struct {path: []const u8, query: []const u8} {
    const parts = std.mem.split(target, "?");
    if (parts.len < 2) {
        return error.InvalidTarget;
    }
    return parts[1];
}

fn budgetGet(req: *http.Request, _: std.mem.Allocator) anyerror!void {
    // Placeholder implementation for GET /budget/
    try req.respond("Budget GET response", .{ .status = http.Status.ok });
}

fn balancePost(req: *http.Request, _: std.mem.Allocator) anyerror!void {
    // Placeholder implementation for POST /balance/
    try req.respond("Balance POST response", .{ .status = http.Status.ok });
}

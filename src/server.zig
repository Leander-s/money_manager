const std = @import("std");
const Server = std.net.Server;
const Address = std.net.Address;
const Stream = std.net.Stream;
const Connection = std.net.Server.Connection;
const Thread = std.Thread;

const Self = @This();

pub fn run(ip: [4]u8, port: u16) !void {
    const address: Address = Address.initIp4(ip, port);
    var server = try address.listen(.{ .reuse_address = true });
    std.log.info("Listening on {d}.{d}.{d}.{d}\n", .{ ip[0], ip[1], ip[2], ip[3] });
    defer server.deinit();
    const allocator = std.heap.page_allocator;
    while (true) {
        var conn = try server.accept();
        const clientThread: Thread = try Thread.spawn(.{ .allocator = allocator }, handleConnection, .{&conn});
        clientThread.detach();
    }
}

fn handleConnection(conn: *Connection) !void {
    defer conn.stream.close();
    var addressBuf:[64]u8 = undefined;
    var addressWriter = std.io.Writer.fixed(&addressBuf);
    try conn.address.in.format(&addressWriter);
    std.debug.print("Client connected: {s}\n", .{addressWriter.buffered()});

    var recBuf: [4096]u8 = undefined;
    var sendBuf: [4096]u8 = undefined;
    while (true) {
        const n = try conn.stream.read(&recBuf);
        std.log.info("Received: {s}\n", .{recBuf[0..n]});
        if (n == 0) break;
        const prefix = "Received: ";
        @memcpy(sendBuf[0..prefix.len], prefix);
        @memcpy(sendBuf[prefix.len..prefix.len+n], recBuf[0..n]);
        const msgLength = prefix.len + n;

        try conn.stream.writeAll(sendBuf[0..msgLength]);
    }
}

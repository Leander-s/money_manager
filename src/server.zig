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
    std.log.info("Client connected: {any}", .{conn.address});

    var recBuf: [4096]u8 = undefined;
    var sendBuf: [4096]u8 = undefined;
    @memset(&sendBuf, 0);
    @memset(&recBuf, 0);
    while (true) {
        const n = try conn.stream.read(&recBuf);
        if (n == 0) break;
        const prefix = "Received: ";
        @memcpy(sendBuf[0..prefix.len], prefix);
        var index: usize = 0;
        while (index < recBuf.len) {
            const sendBufIndex = index + sendBuf.len;
            if (sendBufIndex >= 4096) {
                break;
            }

            sendBuf[sendBufIndex] = recBuf[index];
            index += 1;
        }

        try conn.stream.writeAll(&sendBuf);
    }
}

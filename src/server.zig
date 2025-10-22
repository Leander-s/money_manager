const std = @import("std");
const Server = std.net.Server;
const Address = std.net.Address;
const Stream = std.net.Stream;
const Connection = std.net.Server.Connection;
const Thread = std.Thread;

pub fn run(ip: [4]u8, port: u16) !void {
    const address: Address = Address.initIp4(ip, port);
    const server = try address.listen(.{});
    const allocator = std.heap.page_allocator;
    while (true) {
        const conn = try server.accept();
        const clientThread: Thread = Thread.spawn(.{ .allocator = allocator, .stack_size = 4096 }, handleConnection, &conn);
    }
}

pub fn handleConnection(conn: *Connection) void {}

const std = @import("std");

file: std.fs.File,
buffer: [1024]u8,

const Self = @This();

pub fn init(path: []const u8) !Self {
    var self: Self = undefined;
    self.buffer = undefined;
    self.file = std.fs.cwd().openFile(path, .{ .mode = .read_write }) catch try std.fs.cwd().createFile(path, .{});
    return self;
}

pub fn update(self: *Self, new_amount: f32) !void {
    try self.file.writer(&self.buffer).interface.print("{d}", .{new_amount});
}

pub fn read(self: *Self) !f32 {
    const readBuffer: []u8 = undefined;
    const n = try self.file.reader(&self.buffer).read(&readBuffer);
    if (n == 0) return 0;
    const number = try std.fmt.parseFloat(f32, readBuffer);
    return number;
}

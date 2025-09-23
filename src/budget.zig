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

pub fn destroy(self: *Self) void {
    self.file.close();
}

pub fn update(self: *Self, new_amount: f32) !void {
    var writer: std.fs.File.Writer = .init(self.file, &self.buffer);
    try writer.interface.print("{d}", .{new_amount});
    try writer.interface.flush();
}

pub fn read(self: *Self) !f32 {
    const n = try self.file.read(&self.buffer);
    const buffer = self.buffer[0..self.buffer.len];
    std.debug.print("{s}\n", .{buffer});
    if (n == 0) return 0;
    const number = try std.fmt.parseFloat(f32, buffer);
    return number;
}

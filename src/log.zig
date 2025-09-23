const std = @import("std");

file: std.fs.File,
buffer: [1024]u8,

const Self = @This();

pub fn init(path: []const u8) !Self {
    var self: Self = undefined;
    self.buffer = undefined;
    self.file = std.fs.cwd().openFile(path, .{ .mode = .read_only }) catch try std.fs.cwd().createFile(path, .{ .read = true });
    return self;
}

pub fn update(self: *Self, new_amount: f32) !void {
    const tempFile = try std.fs.cwd().createFile("log_temp", .{});

    var buffer: [1024]u8 = undefined;

    var writer = tempFile.writer(&buffer);
    try writer.interface.print("{d},", .{new_amount});
    const stat = try tempFile.stat();
    try tempFile.seekTo(stat.size);
    _ = try self.file.read(&self.buffer);
    std.debug.print("{s}\n", .{self.buffer});
    _ = try writer.interface.writeAll(&self.buffer);

    try writer.interface.flush();

    self.file.close();
    self.file = tempFile;
    try std.fs.cwd().deleteFile("log");
    try std.fs.cwd().rename("log_temp", "log");
}

pub fn getLastNumber(self: *Self) !f32 {
    var reader = self.file.reader(&self.buffer);
    const numberString = reader.interface.takeDelimiterExclusive(',') catch {
        return 0;
    };
    const number = try std.fmt.parseFloat(f32, numberString);
    return number;
}

pub fn destroy(self: *Self) void {
    self.file.close();
}

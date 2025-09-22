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
    const tempFile = try std.fs.cwd().createFile("log_temp", .{});
    defer tempFile.close();

    var buffer: [1024]u8 = undefined;

    var writer = tempFile.writer(&buffer).interface;
    try writer.print("{d},", .{new_amount});
    while (true) {
        const n = try self.file.read(&buffer);
        if (n == 0) break;
        try tempFile.writeAll(&buffer);
    }

    try std.fs.cwd().deleteFile("log");
    try std.fs.cwd().rename("log_tmp", "log");
}

pub fn getLastNumber(self: *Self) !f32 {
    var buffer: [1024]u8 = undefined;
    const n = try self.file.read(&buffer);
    if (n == 0) return 0;
    var i: usize = 0;
    while (buffer[i] != ',') {
        i += 1;
    }
    const numberString = buffer[0..i];
    const number = try std.fmt.parseFloat(f32, numberString);
    return number;
}

const std = @import("std");

pub fn read() []const u8 {
    const file = try std.fs.cwd().openFile("log", .{ .mode = .read_only });
    defer file.close();

    const tempFile = try std.fs.cwd().createFile("log_temp", .{});
    defer tempFile.close();

    const buffer: [1024]u8 = undefined;
    var n = try file.read(&buffer);
    var i = 0;
    while (buffer[i] != ',') {
        i += 1;
    }
    const numberString = buffer[0..i];
    return numberString;
}

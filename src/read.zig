const std = @import("std");

pub fn read() !f32 {
    const file = try std.fs.cwd().openFile("log", .{ .mode = .read_only });
    defer file.close();

    var buffer: [1024]u8 = undefined;
    var reader: std.fs.File.Reader = .init(file, &buffer);
    const numberString = reader.interface.buffered();
    const number = try std.fmt.parseFloat(f32, numberString);
    return number;
}

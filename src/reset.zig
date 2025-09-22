const std = @import("std");

pub fn reset() !f32 {
    try std.fs.cwd().deleteFile("data");
    try std.fs.cwd().deleteFile("log");
    return 0;
}

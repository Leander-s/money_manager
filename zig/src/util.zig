const std = @import("std");
const expect = std.testing.expect;

pub fn contains(str: []const u8, val: []const u8) ?usize {
    if (str.len < val.len) return null;

    var index: usize = 0;
    var i: usize = 0;
    var j: usize = 0;
    while (i < str.len) {
        if (val[j] == str[i]) {
            i += 1;
            j += 1;
            if (j == val.len) return index;
            continue;
        }
        if (j > 0) i -= 1;
        j = 0;
        i += 1;
        index = i;
    }
    return null;
}

/// Makes a path to the file location and returns the opened file
/// File returned needs to be closed by caller
pub fn openFileAbsoluteMakePath(path: []const u8) !std.fs.File {
    var base = try std.fs.openDirAbsolute("/", .{});
    defer base.close();

    // This should get the parent directory of base in orelse
    try base.makePath(std.fs.path.dirname(path) orelse "/");

    const file = std.fs.createFileAbsolute(path, .{}) catch |e| {
        std.debug.print("Failed to create config file at save.\n", .{});
        return e;
    };
    return file;
}

test "string contains" {
    const testStr = "Hello World";
    const sucTest = "lo Wo";
    const failTest = "Leander";

    try expect(contains(testStr, sucTest).? == 3);
    try expect(contains(testStr, failTest) == null);
}

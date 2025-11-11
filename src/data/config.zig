const std = @import("std");
const fs = std.fs;

const contains = @import("util").contains;

ratio: f32,
changed: bool,

const Self = @This();

pub fn load(path: []const u8) !Self {
    var self = defaultConfig();
    if (path.len == 0) {
        return self;
    }

    const file = std.fs.openFileAbsolute(path, .{ .mode = .read_only }) catch {
        self.changed = true;
        return self;
    };
    defer file.close();

    var buffer: [4096]u8 = undefined;
    var reader = file.reader(&buffer);

    while (true) {
        const line = reader.interface.takeDelimiterExclusive('\n') catch {
            break;
        };

        // skip \n
        try reader.seekBy(1);

        const value = findValue(line) catch {
            std.debug.print("Failed to find value in config line, using default.\n", .{});
            continue;
        };

        if (contains(line, "ratio") != null) {
            self.ratio = parseRatio(value) catch {
                continue;
            };
        }
    }
    return self;
}

fn defaultConfig() Self {
    return Self{
        .ratio = 0.5,
        .changed = false,
    };
}

pub fn save(self: *Self, path: []const u8) !void {
    var base = try std.fs.openDirAbsolute("/", .{});
    defer base.close();

    try base.makePath(std.fs.path.dirname(path) orelse "/");

    var file = std.fs.createFileAbsolute(path, .{}) catch {
        std.debug.print("Failed to create config file at save.\n", .{});
        return;
    };
    defer file.close();

    var buffer: [1024]u8 = undefined;
    var writer = file.writer(&buffer);

    // save ratio
    try writer.interface.print("ratio={d}\n", .{self.ratio});

    try writer.seekTo(0);
    try writer.interface.flush();
}

fn findValue(line: []const u8) ![]const u8 {
    const eqIndex = contains(line, "=") orelse return error.InvalidFormat;
    var value = line[eqIndex + 1 ..];
    var index: usize = 0;
    while (value[index] == ' ' or value[index] == '\t') : (index += 1) {}
    return value[index..];
}

fn parseRatio(line: []const u8) !f32 {
    const parsedRatio = std.fmt.parseFloat(f32, line) catch {
        std.debug.print("Failed to parse ratio from config, using default.\n", .{});
        return error.InvalidFormat;
    };
    return parsedRatio;
}

pub fn updateRatio(self: *Self, newRatio: f32) void {
    self.ratio = newRatio;
    self.changed = true;
}

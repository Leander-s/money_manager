const std = @import("std");
const fs = std.fs;
const expect = std.testing.expect;
const AutoHashMap = std.AutoHashMap;

const contains = @import("util").contains;
const openFileAbsoluteMakePath = @import("util").openFileAbsoluteMakePath;

const Parser = *const fn (self: *Self, value: []const u8) anyerror!void;

const configKey = enum {
    ratio,
};

pub const configKeyMap = std.StaticStringMap(configKey).initComptime(.{
    .{ "ratio", .ratio},
});


pub const ConfigEntry = struct {
    key: configKey,
    value: []const u8,

    pub fn parseEntry(line: []const u8) !ConfigEntry {
        const eqIndex = contains(line, "=") orelse return error.InvalidFormat;
        const value = line[eqIndex + 1 ..];
        const key = line[0..eqIndex];
        var result: ConfigEntry = undefined;
        const keyString = std.mem.trim(u8, key, " \t\n");
        result.key = configKeyMap.get(keyString) orelse return error.InvalidFormat;
        result.value = std.mem.trim(u8, value, " \t\n");
        return result;
    }
};

ratio: f32,
changed: bool,

const Self = @This();

pub fn load(path: []const u8) !Self {
    // Start with default config
    var self = defaultConfig();
    if (path.len == 0) {
        return self;
    }

    // open config file in path
    const file = std.fs.openFileAbsolute(path, .{ .mode = .read_only }) catch {
        self.changed = true;
        return self;
    };
    defer file.close();

    try self.parseConfigFile(&file);

    return self;
}

fn parseConfigFile(self: *Self, file: *const std.fs.File) !void {
    var buffer: [4096]u8 = undefined;
    var reader = file.reader(&buffer);
    const allocator = std.heap.page_allocator;
    var configValueMap = AutoHashMap(configKey, Parser).init(allocator);
    try configValueMap.put(.ratio, parseRatio);

    while (true) {
        const line = reader.interface.takeDelimiterExclusive('\n') catch {
            break;
        };

        // skip \n
        try reader.seekBy(1);

        const entry = ConfigEntry.parseEntry(line) catch {
            std.debug.print("Failed to find value in config line, using default.\n", .{});
            continue;
        };

        const parser = configValueMap.get(entry.key) orelse continue;
        try parser(self, entry.value);
    }
}

fn defaultConfig() Self {
    return Self{
        .ratio = 0.5,
        .changed = false,
    };
}

pub fn save(self: *Self, path: []const u8) !void {
    var file = try openFileAbsoluteMakePath(path);
    defer file.close();

    var buffer: [1024]u8 = undefined;
    var writer = file.writer(&buffer);

    // save ratio
    try writer.interface.print("ratio={d}\n", .{self.ratio});

    try writer.seekTo(0);
    try writer.interface.flush();
}

fn parseRatio(self: *Self, line: []const u8) !void {
    const parsedRatio = std.fmt.parseFloat(f32, line) catch {
        std.debug.print("Failed to parse ratio from config, using default.\n", .{});
        return error.ParseError;
    };
    self.ratio = parsedRatio;
}

pub fn updateEntry(self: *Self, configEntry: *const ConfigEntry) !void {
    switch (configEntry.key) {
        .ratio => {
            const value = std.fmt.parseFloat(f32, configEntry.value) catch return error.InvalidValue;
            self.ratio = value;
        }
    }
    self.changed = true;
}

test "find value in config parser" {
    const entry = try ConfigEntry.parseEntry(" ratio = 0.1\n");
    try expect(std.mem.startsWith(u8, entry.value, "0.1"));
    try expect(entry.key == .ratio);
}

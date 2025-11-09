const std = @import("std");
const contains = @import("util").contains;
const LogEntry = @import("logentry.zig");
const Config = @import("config.zig");

const configLoc = "/.config/money_manager/config";

budget: f32,
allocator: std.mem.Allocator,
entries: std.ArrayList(LogEntry),
config: Config,
configPath: []const u8,

const Self = @This();

pub fn init(fileName: []const u8) !Self {
    var self: Self = undefined;
    self.allocator = std.heap.page_allocator;
    self.entries = std.ArrayList(LogEntry).empty;

    const homeDir = try getHomeDir(self.allocator);
    defer self.allocator.free(homeDir);
    self.configPath = std.fmt.allocPrint(self.allocator, "{s}{s}", .{homeDir, configLoc}) catch {
        std.debug.print("Failed to construct config path.\n", .{});
        return error.ConfigPathError;
    };

    self.config = try Config.load(self.configPath);

    const path = try self.getPathToFile(fileName);
    try self.parseFile(path);

    return self;
}

fn initDefault(self: *Self) !void {
    self.budget = 0;
    self.entries = std.ArrayList(LogEntry).empty;
}

fn parseFile(self: *Self, path: []const u8) !void {
    const file = std.fs.cwd().openFile(path, .{ .mode = .read_only }) catch {
        _ = try std.fs.cwd().createFile(path, .{});
        try self.initDefault();
        return;
    };

    var buffer: [1024]u8 = undefined;
    var reader = file.reader(&buffer);
    while (true) {
        // read next line
        const line = reader.interface.takeDelimiterExclusive('\n') catch {
            break;
        };

        // skip \n
        try reader.seekBy(1);

        // Ignore header line
        if (contains(line, "budget") != null)
            continue;

        // parsing line in current log version
        if (contains(line, ",") != null) {
            try self.entries.append(self.allocator, try LogEntry.parse(line));
            continue;
        }

        // parsing old version
        const newEntry = self.parseOldEntry(line) catch {
            // budget line
            continue;
        };
        try self.entries.append(self.allocator, newEntry);
    }

    // if there is a budget -> assign it for quick read op
    if (self.entries.items.len > 0)
        self.budget = self.entries.items[0].budget;
}

fn parseOldEntry(self: *Self, line: []const u8) !LogEntry {
    var budgetString: []const u8 = undefined;
    var balanceString: []const u8 = undefined;
    var timestampString: []const u8 = undefined;

    if (contains(line, ":")) |index| {
        budgetString = "0";
        balanceString = line[0..index];
        timestampString = line[index + 1 .. line.len];
    } else {
        self.budget = try std.fmt.parseFloat(f32, line);
        return error.NoLine;
    }

    // parsing numbers
    const budget = std.fmt.parseFloat(f32, budgetString) catch {
        return error.WrongFormat;
    };
    const balance = std.fmt.parseFloat(f32, balanceString) catch {
        return error.WrongFormat;
    };
    const timestamp = std.fmt.parseInt(i64, timestampString, 10) catch {
        return error.WrongFormat;
    };

    return .{.budget = budget, .balance = balance, .timestamp = timestamp, .ratio = 0.5};
}

fn getPathToFile(self: *Self, fileName: []const u8) ![]const u8 {
    const exe_dir = try std.fs.selfExeDirPathAlloc(self.allocator);
    defer self.allocator.free(exe_dir);

    var pathList = std.ArrayList(u8).empty;
    defer pathList.deinit(self.allocator);

    try pathList.appendSlice(self.allocator, exe_dir);
    try pathList.appendSlice(self.allocator, "/");
    try pathList.appendSlice(self.allocator, fileName);

    const path = try pathList.toOwnedSlice(self.allocator);

    return path;
}

pub fn write(self: *Self, fileName: []const u8) !void {
    const path = try self.getPathToFile(fileName);
    var file = try std.fs.cwd().createFile(path, .{});

    var buffer: [1024]u8 = undefined;
    var writer = file.writer(&buffer);
    try LogEntry.writeHeader(&writer);

    var index: usize = 0;
    while (index < self.entries.items.len) {
        var entry = self.entries.items[index];
        try entry.write(&writer);
        index += 1;
    }
    try writer.seekTo(0);
    try writer.interface.flush();

    if (self.config.changed) {
        try self.config.save(self.configPath);
    }
}

pub fn read(self: *Self) f32 {
    return self.budget;
}

pub fn enter(self: *Self, number: f32) !f32 {
    var lastEntry: ?*LogEntry = null;
    if (self.entries.items.len > 0) {
        lastEntry = &self.entries.items[0];
    }
    const newEntry = LogEntry.init(lastEntry, number, self.config.ratio);
    try self.entries.insert(self.allocator, 0, newEntry);
    self.budget = newEntry.budget;
    return self.budget;
}

pub fn reset(self: *Self) f32 {
    self.entries.clearAndFree(self.allocator);
    self.budget = 0;
    return 0;
}

pub fn destroy(self: *Self) void {
    self.entries.clearAndFree(self.allocator);
}

pub fn recalculateBudgets(self: *Self) f32{
    var index: usize = self.entries.items.len - 1;
    while (true) {
        var entry = self.entries.items[index];
        const ratio = entry.ratio;
        var previousEntry: ?*LogEntry = null;
        if (index < self.entries.items.len - 1) previousEntry = &self.entries.items[index + 1];
        self.entries.items[index] = entry.recalculate(previousEntry, ratio);
        if (index == 0) break;
        index -= 1;
    }
    self.budget = self.entries.items[0].budget;
    return self.budget;
}

fn getHomeDir(allocator: std.mem.Allocator) ![]const u8 {
    var env = std.process.getEnvMap(allocator) catch return error.EnvVarError;
    defer env.deinit();
    if (env.get("HOME")) |home| {
        std.debug.print("Home dir: {s}\n", .{home});
        return allocator.dupe(u8, home) catch return error.HomePathAllocError;
    }
    if (env.get("USERPROFILE")) |home| {
        return allocator.dupe(u8, home) catch return error.HomePathAllocError;
    }
    return error.HomeNotFound;
}

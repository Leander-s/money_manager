const std = @import("std");
const contains = @import("util").contains;
const LogEntry = @import("logentry.zig");
pub const Config = @import("config.zig");

const expect = std.testing.expect;
const expectEqual = std.testing.expectEqual;

const configLoc = "/.config/money_manager/config";

budget: f32,
balance: f32,
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
    self.configPath = prependHomeDir(self.allocator, configLoc) catch {
        std.debug.print("Failed to construct config path.\n", .{});
        return error.ConfigPathError;
    };

    self.config = try Config.load(self.configPath);

    const path = try self.getPathToFile(fileName);
    try self.parseFile(path);

    return self;
}

fn initDefault(self: *Self) !void {
    self.balance = 0;
    self.budget = 0;
    self.entries = std.ArrayList(LogEntry).empty;
}

fn parseFile(self: *Self, path: []const u8) !void {
    const file = std.fs.cwd().openFile(path, .{ .mode = .read_only }) catch {
        const file = try std.fs.cwd().createFile(path, .{});
        defer file.close();
        try self.initDefault();
        return;
    };
    defer file.close();

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
    if (self.entries.items.len > 0) {
        self.budget = self.entries.items[0].budget;
        self.balance = self.entries.items[0].balance;
    }
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

    return .{ .budget = budget, .balance = balance, .timestamp = timestamp, .ratio = 0.5 };
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
    defer file.close();

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

pub fn lastBalance(self: *Self) f32 {
    return self.balance;
}

pub fn currentBudget(self: *Self) f32 {
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
    self.balance = number;
    return self.budget;
}

pub fn reset(self: *Self) f32 {
    self.entries.clearAndFree(self.allocator);
    self.budget = 0;
    self.balance = 0;
    return 0;
}

pub fn destroy(self: *Self) void {
    self.allocator.free(self.configPath);
    self.entries.clearAndFree(self.allocator);
}

pub fn recalculateBudgets(self: *Self) f32 {
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
    self.balance = self.entries.items[0].balance;
    return self.budget;
}

fn getHomeDir(allocator: std.mem.Allocator) ![]const u8 {
    var env = std.process.getEnvMap(allocator) catch return error.EnvVarError;
    defer env.deinit();
    if (env.get("HOME")) |home| {
        return allocator.dupe(u8, home) catch return error.HomePathAllocError;
    }
    if (env.get("USERPROFILE")) |home| {
        return allocator.dupe(u8, home) catch return error.HomePathAllocError;
    }
    return error.HomeNotFound;
}

fn prependHomeDir(allocator: std.mem.Allocator, path: []const u8) ![]const u8 {
    const homeDir = try getHomeDir(allocator);
    defer allocator.free(homeDir);
    const result = std.fmt.allocPrint(allocator, "{s}{s}", .{ homeDir, path }) catch {
        std.debug.print("Failed to prepend home dir path.\n", .{});
        return error.PrependHomeDirError;
    };
    return result;
}

test "config test" {
    const testLoc = "/.config/money_manager/test_config";
    const testAlloc = std.heap.page_allocator;
    const testPath = try prependHomeDir(testAlloc, testLoc);
    defer testAlloc.free(testPath);
    var config = try Config.load(testPath);
    try expectEqual(true, config.changed);
    try expectEqual(0.5, config.ratio);
    try config.updateEntry(&.{ .key = .ratio, .value = "0.3" });
    try config.save(testPath);
    var newConfig = try Config.load(testPath);
    try expectEqual(false, newConfig.changed);
    try expectEqual(0.3, newConfig.ratio);
    try newConfig.updateEntry(&.{ .key = .ratio, .value = "0.5" });
    try expectEqual(true, newConfig.changed);
    try expectEqual(0.5, newConfig.ratio);
    try newConfig.save(testPath);
    const lastConfig = try Config.load(testPath);
    try expectEqual(0.5, lastConfig.ratio);
    try expectEqual(false, lastConfig.changed);
    try std.fs.deleteFileAbsolute(testPath);
}

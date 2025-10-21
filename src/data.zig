const std = @import("std");
const contains = @import("util.zig").contains;
const LogEntry = @import("logentry.zig");

budget: f32,
allocator: std.mem.Allocator,
entries: std.ArrayList(LogEntry),

const Self = @This();

pub fn init(fileName: []const u8) !Self {
    var self: Self = undefined;
    self.allocator = std.heap.page_allocator;
    self.entries = std.ArrayList(LogEntry).empty;

    const path = try self.getPathToFile(fileName);
    try self.parseFile(path);

    return self;
}

fn initDefault() Self {
    return .{ .allocator = std.heap.page_allocator, .budget = 0, .entries = std.ArrayList(LogEntry).empty };
}

fn parseFile(self: *Self, path: []const u8) !void {
    const file = std.fs.cwd().openFile(path, .{ .mode = .read_only }) catch {
        _ = try std.fs.cwd().createFile(path, .{});
        self.* = initDefault();
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
}

pub fn read(self: *Self) f32 {
    return self.budget;
}

pub fn enter(self: *Self, number: f32) !f32 {
    var lastEntry: ?*LogEntry = null;
    if (self.entries.items.len > 0) {
        lastEntry = &self.entries.items[0];
    }
    const newEntry = LogEntry.init(lastEntry, number, 0.5);
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
        var previousEntry: ?*LogEntry = null;
        if (index < self.entries.items.len - 1) previousEntry = &self.entries.items[index + 1];
        self.entries.items[index] = entry.recalculate(previousEntry, 0.5);
        if (index == 0) break;
        index -= 1;
    }
    self.budget = self.entries.items[0].budget;
    return self.budget;
}

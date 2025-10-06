const std = @import("std");
const contains = @import("util.zig").contains;

budget: f32,
allocator: std.mem.Allocator,
history: std.ArrayList(f32),
timestamps: std.ArrayList(i64),

const Self = @This();

pub fn init(fileName: []const u8) !Self {
    var self: Self = undefined;
    self.allocator = std.heap.page_allocator;
    self.history = std.ArrayList(f32).empty;
    self.timestamps = std.ArrayList(i64).empty;

    const path = try self.getPathToFile(fileName);

    const file = std.fs.cwd().openFile(path, .{ .mode = .read_only }) catch {
        _ = try std.fs.cwd().createFile(path, .{});
        return initDefault();
    };

    var buffer: [1024]u8 = undefined;
    var reader = file.reader(&buffer);
    const budgetString = reader.interface.takeDelimiterExclusive('\n') catch {
        return error.WrongFormat;
    };
    self.budget = try std.fmt.parseFloat(f32, budgetString);
    while (true) {
        const historyString = reader.interface.takeDelimiterExclusive('\n') catch {
            break;
        };
        var historyNumberString: []const u8 = undefined;
        var historyTimestampString: []const u8 = undefined;

        if (contains(historyString, ":")) |index| {
            historyNumberString = historyString[0..index];
            historyTimestampString = historyString[index + 1 .. historyString.len];
        } else {
            historyNumberString = historyString;
            historyTimestampString = "0";
        }
        const historyNumber = std.fmt.parseFloat(f32, historyNumberString) catch {
            return error.WrongFormat;
        };
        const historyTimestamp = std.fmt.parseInt(i64, historyTimestampString, 10) catch {
            return error.WrongFormat;
        };
        try self.history.append(self.allocator, historyNumber);
        try self.timestamps.append(self.allocator, historyTimestamp);
    }
    return self;
}

fn initDefault() Self {
    return .{ .allocator = std.heap.page_allocator, .budget = 0, .history = std.ArrayList(f32).empty, .timestamps = std.ArrayList(i64).empty };
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
    try writer.interface.print("{d}\n", .{self.budget});
    var stat = try file.stat();
    try writer.seekTo(stat.size);
    var historyIndex: usize = 0;
    while (historyIndex < self.history.items.len) {
        try writer.interface.print("{d}:{d}\n", .{self.history.items[historyIndex], self.timestamps.items[historyIndex]});
        stat = try file.stat();
        try writer.seekTo(stat.size);
        historyIndex += 1;
    }
    try writer.seekTo(0);
    try writer.interface.flush();
}

pub fn read(self: *Self) f32 {
    return self.budget;
}

pub fn enter(self: *Self, number: f32) !f32 {
    const timestamp = std.time.timestamp();
    try self.history.insert(self.allocator, 0, number);
    try self.timestamps.insert(self.allocator, 0, timestamp);
    self.recalculateBudget();
    return self.budget;
}

pub fn reset(self: *Self) f32 {
    self.history.clearAndFree(self.allocator);
    self.timestamps.clearAndFree(self.allocator);
    self.budget = 0;
    return 0;
}

pub fn destroy(self: *Self) void {
    self.history.clearAndFree(self.allocator);
}

fn recalculateBudget(self: *Self) void {
    if (self.history.items.len < 2) {
        self.budget = 0;
        return;
    }

    const before = self.history.items[1];
    const after = self.history.items[0];
    const diff = after - before;

    if (diff < 0) {
        // Spendings get taken from the budget
        self.budget += diff;
    } else {
        // Save half the income
        self.budget += diff / 2;
    }
    self.budget *= 100;
    self.budget = @round(self.budget);
    self.budget /= 100;
}

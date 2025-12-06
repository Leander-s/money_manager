const std = @import("std");
const time = std.time;
const expect = std.testing.expect;

const DateTime = @import("datetime.zig");
const contains = @import("util").contains;

budget: f32,
balance: f32,
timestamp: DateTime,
ratio: f32,

const Self = @This();

pub fn init(optLastEntry: ?*Self, newBalance: f32, currentRatio: f32) Self {
    const timestamp = DateTime.now();
    var self: Self = .{ .budget = 0, .balance = newBalance, .timestamp = timestamp, .ratio = currentRatio };
    if (optLastEntry == null) {
        return self;
    }
    const lastEntry = optLastEntry.?;
    self.budget = lastEntry.budget;

    const diff = calculateBudgetDiff(lastEntry.balance, self.balance, currentRatio);
    self.budget += diff;

    // round the budget to 2 places
    self.budget *= 100;
    self.budget = @round(self.budget);
    self.budget /= 100;
    return self;
}

// This function is meant to be used to recalculate an entry. It essentially reinitializes the
// entry with the current ratio but gives it its old timestamp.
// This way you can calculate a budget with a different previous entry or more commonly, a
// different ratio
pub fn recalculate(self: *Self, optLastEntry: ?*Self, currentRatio: f32) Self {
    var entry = init(optLastEntry, self.balance, currentRatio);
    entry.timestamp = self.timestamp;
    return entry;
}

pub fn writeHeader(writer: *std.fs.File.Writer) !void {
    const stat = try writer.file.stat();
    try writer.interface.print("budget,balance,timestamp,ratio,\n", .{});
    try writer.seekTo(stat.size);
}

pub fn write(self: *Self, writer: *std.fs.File.Writer) !void {
    const stat = try writer.file.stat();
    try writer.interface.print("{d},{d},{s},{d},\n", .{ self.budget, self.balance, try self.timestamp.ISO8601(), self.ratio });
    try writer.seekTo(stat.size);
}

// parses a line of the log
pub fn parse(line: []const u8) !Self {
    var leftover = line;
    var data: [4][]const u8 = undefined;
    var index: usize = 0;

    // Insert all items read into an array using , as delimiter
    while (contains(leftover, ",")) |pos| {
        data[index] = leftover[0..pos];
        leftover = leftover[pos + 1 ..];
        index += 1;
    }

    // Make sure data is filled completely
    while (index < 4) {
        data[index] = "0";
        index += 1;
    }

    // Parse all strings and store them in the struct
    var self: Self = undefined;
    self.budget = std.fmt.parseFloat(f32, data[0]) catch {
        return error.WrongFormat;
    };
    self.balance = std.fmt.parseFloat(f32, data[1]) catch {
        return error.WrongFormat;
    };
    self.timestamp = try parseTimeStamp(data[2]);
    self.ratio = std.fmt.parseFloat(f32, data[3]) catch {
        return error.WrongFormat;
    };
    return self;
}

fn parseTimeStamp(timeString: []const u8) !DateTime {
    if (!std.mem.endsWith(u8, timeString, "Z")) {
        const timeStamp = std.fmt.parseInt(i64, timeString, 10) catch {
            return error.WrongFormat;
        };
        return DateTime.init(timeStamp);
    }

    return DateTime.parseISO8601(timeString) catch {
        return error.WrongFormat;
    };
}

fn calculateBudgetDiff(oldBalance: f32, newBalance: f32, ratio: f32) f32 {
    const diff = newBalance - oldBalance;

    if (diff < 0) {
        // Spendings get taken from the budget
        return diff;
    } else {
        // The new budget is the income times the current ratio
        return diff * ratio;
    }
}

test "100 income" {
    const testBalanceOld = 100;
    const testBalanceNew = 200;

    var diff = calculateBudgetDiff(testBalanceOld, testBalanceNew, 0.5);
    try expect(diff == 50);

    diff = calculateBudgetDiff(testBalanceOld, testBalanceNew, 0.1);
    try expect(diff == 10);

    diff = calculateBudgetDiff(testBalanceOld, testBalanceNew, 0.9);
    try expect(diff == 90);
}

test "100 spending" {
    const testBalanceOld = 200;
    const testBalanceNew = 100;

    var diff = calculateBudgetDiff(testBalanceOld, testBalanceNew, 0.5);
    try expect(diff == -100);

    diff = calculateBudgetDiff(testBalanceOld, testBalanceNew, 0.1);
    try expect(diff == -100);

    diff = calculateBudgetDiff(testBalanceOld, testBalanceNew, 0.9);
    try expect(diff == -100);
}

const std = @import("std");
const time = std.time;

const contains = @import("util.zig").contains;

budget: f32,
balance: f32,
timestamp: i64,
ratio: f32,

const Self = @This();

pub fn init(optLastEntry: ?*Self, newBalance: f32, currentRatio: f32) Self {
    const timestamp = time.timestamp();
    var self: Self = .{ .budget = 0, .balance = newBalance, .timestamp = timestamp, .ratio = currentRatio };
    if (optLastEntry == null) {
        return self;
    }
    const lastEntry = optLastEntry.?;
    self.budget = lastEntry.budget;

    const diff = self.balance - lastEntry.balance;

    if (diff < 0) {
        // Spendings get taken from the budget
        self.budget += diff;
    } else {
        // Save half the income
        self.budget += diff * currentRatio;
    }

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
    try writer.interface.print("{d},{d},{d},{d},\n", .{ self.budget, self.balance, self.timestamp, self.ratio });
    try writer.seekTo(stat.size);
}

pub fn parse(line: []const u8) !Self {
    var leftover = line;
    var data: [4][]const u8 = undefined;
    var index: usize = 0;

    // Fülle alle gelesenen Dinge in ein array mit , als delimiter
    while (contains(leftover, ",")) |pos| {
        data[index] = leftover[0..pos];
        leftover = leftover[pos + 1..];
        index += 1;
    }

    // Sicher gehen, dass data komplett gefüllt ist
    while (index < 4) {
        data[index] = "0";
        index += 1;
    }

    // Alle strings parsen und zuweisen
    var self: Self = undefined;
    self.budget = std.fmt.parseFloat(f32, data[0]) catch {
        return error.WrongFormat;
    };
    self.balance = std.fmt.parseFloat(f32, data[1]) catch {
        return error.WrongFormat;
    };
    self.timestamp = std.fmt.parseInt(i64, data[2], 10) catch {
        return error.WrongFormat;
    };
    self.ratio = std.fmt.parseFloat(f32, data[3]) catch {
        return error.WrongFormat;
    };
    return self;
}

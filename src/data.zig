const std = @import("std");

budget: f32,
allocator: std.mem.Allocator,
history: std.ArrayList(f32),

const Self = @This();

pub fn init(path: []const u8) !Self {
    var self: Self = undefined;
    self.allocator = std.heap.page_allocator;
    self.history = std.ArrayList(f32).empty;

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
        const historyNumber = std.fmt.parseFloat(f32, historyString) catch {
            return error.WrongFormat;
        };
        try self.history.append(self.allocator, historyNumber);
    }
    return self;
}

fn initDefault() Self {
    return .{ .allocator = std.heap.page_allocator, .budget = 0, .history = std.ArrayList(f32).empty };
}

pub fn write(self: *Self, path: []const u8) !void {
    var file = std.fs.cwd().openFile(path, .{ .mode = .write_only }) catch try std.fs.cwd().createFile(path, .{});

    var buffer: [1024]u8 = undefined;
    var writer = file.writer(&buffer);
    try writer.interface.print("{d}\n", .{self.budget});
    var stat = try file.stat();
    try writer.seekTo(stat.size);
    var historyIndex: usize = 0;
    while (historyIndex < self.history.items.len) {
        try writer.interface.print("{d}\n", .{self.history.items[historyIndex]});
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
    try self.history.insert(self.allocator, 0, number);
    self.recalculateBudget();
    return self.budget;
}

pub fn reset(self: *Self) f32 {
    self.history.clearAndFree(self.allocator);
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
        self.budget += diff/2;
    }
}

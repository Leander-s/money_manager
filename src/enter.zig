const std = @import("std");
const Log = @import("log.zig");
const Budget = @import("budget.zig");

pub fn enter(given_arg: ?[]const u8) !f32 {
    const arg = given_arg orelse return error.InvalidArgument;
    const number = try std.fmt.parseFloat(f32, arg);

    var budget = try Budget.init("data");
    defer budget.destroy();

    var log = try Log.init("log");
    defer log.destroy();

    const last_number = try log.getLastNumber();
    try log.update(number);

    const diff = number - last_number;
    const additional_available_money = diff/2;

    const last_available_money = try budget.read();
    std.debug.print("{}\n", .{last_available_money});

    const new_available = last_available_money + additional_available_money; 

    try budget.update(new_available);

    return new_available;
}

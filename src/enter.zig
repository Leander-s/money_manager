const std = @import("std");
const Log = @import("log.zig");
const Budget = @import("budget.zig");

pub fn enter(given_arg: ?[]const u8) !f32 {
    const arg = given_arg orelse return error.InvalidArgument;
    const number = try std.fmt.parseFloat(f32, arg);

    var budget = try Budget.init("data");

    var log = try Log.init("log");

    const last_number = try log.getLastNumber();
    try log.update(number);

    const diff = number - last_number;
    const additional_available_money = diff/2;

    const last_available_money = try budget.read();

    const new_available = last_available_money + additional_available_money; 

    budget.update(new_available);

    return new_available;
}

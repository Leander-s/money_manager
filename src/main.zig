const std = @import("std");
const Writer = std.Io.Writer;
const Arg = @import("arg.zig").Arg;
const Command = @import("arg.zig").Command;
const Data = @import("data");
const Server = @import("server");
const StaticStringMap = std.StaticStringMap;

const Context = struct {
    stdout: *Writer,
    data: *Data,
    args: *const Arg,
};

const Handler = *const fn (*Context) anyerror!void;

const handlerMap = StaticStringMap(Handler).initComptime(.{
    .{ "enter", handleEnter },
    .{ "budget", handleBudget },
    .{ "balance", handleBalance },
    .{ "config", handleConfig },
    .{ "reset", handleReset },
    .{ "recalculate", handleRecalculate },
    .{ "host", handleHost },
});

pub fn main() !void {
    var stdout_buffer: [1024]u8 = undefined;
    var stdout_writer = std.fs.File.stdout().writer(&stdout_buffer);
    const stdout = &stdout_writer.interface;

    var args = std.process.args();
    // skipping the exe name
    _ = args.skip();

    const arg = Arg.parse(&args);

    var data: Data = try Data.init("log");

    const handler = handlerMap.get(arg.command) orelse handleInvalid;
    var context: Context = .{ .data = &data, .args = &arg, .stdout = stdout };
    handler(&context) catch {
        try stdout.print("No valid argument given\n", .{});
    };

    try data.write("log");
    data.destroy();

    try stdout.flush();
}

fn handleEnter(ctx: *Context) anyerror!void {
    const args = ctx.args;
    var data = ctx.data;
    const value = args.value orelse return error.InvalidArgument;
    const budget: f32 = try data.enter(value);
    try ctx.stdout.print("{d}\n", .{budget});
}

fn handleBudget(ctx: *Context) !void {
    var data = ctx.data;
    const budget: f32 = data.currentBudget();
    try ctx.stdout.print("{d}\n", .{budget});
}

fn handleBalance(ctx: *Context) !void {
    var data = ctx.data;
    const balance: f32 = data.lastBalance();
    try ctx.stdout.print("{d}\n", .{balance});
}

fn handleConfig(ctx: *Context) !void {
    const args = ctx.args;
    var data = ctx.data;
    const configEntry = args.configEntry orelse {
        try ctx.stdout.print("Not a valid config entry\n", .{});
        return error.InvalidArgument;
    };
    try data.config.updateEntry(&configEntry);
}

fn handleReset(ctx: *Context) !void {
    ctx.data.reset();
    try ctx.stdout.print("Data was reset\n", .{});
}

fn handleRecalculate(ctx: *Context) !void {
    const budget: f32 = ctx.data.recalculateBudgets();
    try ctx.stdout.print("{d}\n", .{budget});
}

fn handleHost(ctx: *Context) !void {
    const address: [4]u8 = .{ 127, 0, 0, 1 };
    const port = 8080;
    try ctx.stdout.print("Running server on {s}:{d}\n", .{ address, port });
    try Server.run(address, port);
}

fn handleInvalid(ctx: *Context) !void {
    try ctx.stdout.print("Invalid argument given\n", .{});
}

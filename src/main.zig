const std = @import("std");
const Arg = @import("arg.zig").Arg;
const Data = @import("data.zig");
const Server = @import("server.zig");

pub fn main() !void {
    var stdout_buffer: [1024]u8 = undefined;
    var stdout_writer = std.fs.File.stdout().writer(&stdout_buffer);
    const stdout = &stdout_writer.interface;

    var args = std.process.args();
    _ = args.skip();

    const arg = Arg.parse(args.next());
    const argString = args.next() orelse "0";
    const argNumber = std.fmt.parseFloat(f32, argString) catch {
        try stdout.print("Not a valid argument\n", .{});
        try stdout.flush();
        return;
    };

    var data: Data = try Data.init("log");

    const result: ?f32 = switch (arg) {
        .enter => try data.enter(argNumber),
        .read => data.read(),
        .reset => data.reset(),
        .recalculate => data.recalculateBudgets(),
        .runServer => blk: {
            try Server.run(.{ 127, 0, 0, 1 }, 8080);
            break :blk null;
        },
        .noArg, .unknown => null,
    };

    if (result) |value| {
        try data.write("log");
        data.destroy();
        try stdout.print("{d}\n", .{value});
    } else {
        try stdout.print("No valid argument given\n", .{});
    }

    try stdout.flush();
}

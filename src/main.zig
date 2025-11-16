const std = @import("std");
const Arg = @import("arg.zig").Arg;
const Data = @import("data");
const Server = @import("server");

pub fn main() !void {
    var stdout_buffer: [1024]u8 = undefined;
    var stdout_writer = std.fs.File.stdout().writer(&stdout_buffer);
    const stdout = &stdout_writer.interface;

    var args = std.process.args();
    // skipping the exe name
    _ = args.skip();

    const arg = Arg.parse(&args);

    var data: Data = try Data.init("log");

    const result: ?f32 = switch (arg.command) {
        .enter => if (arg.value) |value| try data.enter(value) else null,
        .config => blk: {
            if (arg.configEntry) |configEntry| {
                try data.config.updateEntry(&configEntry);
                break :blk data.read();
            } else {
                std.debug.print("Not a valid config entry\n", .{});
                return error.InvalidFormat;
            }
        },
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

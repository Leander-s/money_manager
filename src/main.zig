const std = @import("std");
const Arg = @import("arg.zig").Arg;
const enter = @import("enter.zig").enter;
const read = @import("read.zig").read;
const reset = @import("reset.zig").reset;

pub fn main() !void {
    var args = std.process.args();
    _ = args.skip();

    const arg = Arg.parse(args.next());

    const result: ?f32 = switch (arg) {
        .enter => try enter(args.next()),
        .read => try read(),
        .reset => try reset(),
        .noArg, .unknown => null,
    };

    var stdout_buffer: [1024]u8 = undefined;
    var stdout_writer = std.fs.File.stdout().writer(&stdout_buffer);
    const stdout = &stdout_writer.interface;

    if (result) |value| {
        try stdout.print("{d}\n", .{value});
    } else {
        try stdout.print("No valid argument given\n", .{});
    }

    try stdout.flush();
}

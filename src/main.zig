const std = @import("std");
const Arg = @import("arg.zig").Arg;

pub fn main() !void {
    var args = std.process.args();
    _ = args.skip();

    const arg = Arg.parse(args.next());

    const result = switch (arg) {
        .enter => "Enter amount",
        .read => "rich af",
        .reset => "resetting...",
        .noArg, .unknown => "no valid argument", 
    };

    var stdout_buffer: [1024]u8 = undefined;
    var stdout_writer = std.fs.File.stdout().writer(&stdout_buffer);
    const stdout = &stdout_writer.interface;

    try stdout.print("{s}\n", .{result});
    try stdout.flush();
}

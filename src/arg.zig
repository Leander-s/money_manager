const std = @import("std");

pub const Arg = enum {
    enter,
    read,
    reset,
    unknown,
    noArg,

    pub fn parse(str: ?[]const u8) Arg {
        const hashmap = comptime std.StaticStringMap(Arg).initComptime(.{
            .{"enter", .enter},
            .{"read", .read},
            .{"reset", .reset},
            .{"", .noArg},
            .{"unknown", .unknown},
        });

        return hashmap.get(str orelse "") orelse .unknown;
    }
};

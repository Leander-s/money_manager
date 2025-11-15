const std = @import("std");
const ArgIterator = std.process.ArgIterator;
const ConfigEntry = @import("data").Config.ConfigEntry;
const ConfigKeyMap = @import("data").Config.configKeyMap;

pub const Command = enum {
    enter,
    read,
    reset,
    recalculate,
    runServer,
    config,
    unknown,
    noArg,
};
pub const Arg = struct {
    command: Command,
    value: ?f32,
    configEntry: ?ConfigEntry,

    pub fn parse(args: ArgIterator) Arg {
        var self: Arg = .{ .command = .unknown, .value = null, .configEntry = null };
        const arg = args.next();
        const commandMap = comptime std.StaticStringMap(Arg).initComptime(.{
            .{ "enter", .enter },
            .{ "read", .read },
            .{ "reset", .reset },
            .{ "recalculate", .recalculate },
            .{ "run", .runServer },
            .{ "config", .config },
            .{ "", .noArg },
            .{ "unknown", .unknown },
        });

        self.command = commandMap.get(arg orelse "") orelse return self;
        switch (self.command) {
            .runServer, .unknown, .reset, .read, .recalculate => return self,
            .enter => {
                self.value = std.fmt.parseFloat(f32, arg) catch return self;
                return self;
            },
            .config => {
                const key = ConfigKeyMap.get(arg) orelse return self;
                const value = args.next();
                const configEntry: ConfigEntry = .{ .key = key, .value = value };
            },
        }
    }
};

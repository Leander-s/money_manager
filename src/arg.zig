const std = @import("std");
const ArgIterator = std.process.ArgIterator;
const ConfigEntry = @import("data").Config.ConfigEntry;
const ConfigKeyMap = @import("data").Config.configKeyMap;

pub const Arg = struct {
    command: []const u8,
    value: ?f32,
    configEntry: ?ConfigEntry,

    pub fn parse(args: *ArgIterator) Arg {
        var self: Arg = .{ .command = .unknown, .value = null, .configEntry = null };
        var arg = args.next();

        self.command = arg;
        arg = args.next();
        switch (self.command) {
            .enter => {
                const value = arg orelse return self;
                self.value = std.fmt.parseFloat(f32, value) catch return self;
                return self;
            },
            .config => {
                const entryStr = arg orelse return self;
                self.configEntry = ConfigEntry.parseEntry(entryStr) catch return self;
                return self;
            },
            else => return self,
        }
    }
};

const std = @import("std");
const time = std.time;

timestamp: i64,

const Self = @This();

pub fn now() Self {
    return .{ .timestamp = time.timestamp() };
}

pub fn init(timestamp: i64) Self {
    return .{ .timestamp = timestamp };
}

pub fn parseISO8601(timeString: []const u8) !Self {
    const yearString = timeString[0..4];
    const monthString = timeString[5..7];
    const dayString = timeString[8..10];
    const hourString = timeString[11..13];
    const minuteString = timeString[14..16];
    const secondString = timeString[17..19];

    const year = try std.fmt.parseInt(i64, yearString, 10);
    const month = try std.fmt.parseInt(i64, monthString, 10);
    const day = try std.fmt.parseInt(i64, dayString, 10);
    const hour = try std.fmt.parseInt(i64, hourString, 10);
    const minute = try std.fmt.parseInt(i64, minuteString, 10);
    const second = try std.fmt.parseInt(i64, secondString, 10);

    var timeStamp = second + minute * 60 + hour * 3_600;

    const totalDays = calculateDays(year, month, day);

    const secondsPerDay = 24 * 60 * 60;
    timeStamp += totalDays * secondsPerDay;

    return .{ .timestamp = timeStamp };
}

pub fn ISO8601(self: *const Self) ![]const u8 {
    const seconds = @mod(self.timestamp, 60);
    const mStamp = @divFloor(self.timestamp, 60);
    const minutes = @mod(mStamp, 60);
    const hStamp = @divFloor(mStamp, 60);
    const hours = @mod(hStamp, 24);
    const dStamp = @divFloor(hStamp, 24);

    const z = dStamp + 719_468;

    const era = if (z >= 0) @divFloor(z, 146_097) else @divFloor((z - 146_096), 146_097);
    const dayOfEra = z - era * 146_097;
    const yearOfEra = @divFloor((dayOfEra - @divFloor(dayOfEra, 1_460) + @divFloor(dayOfEra, 36_524) - @divFloor(dayOfEra, 146_096)), 365);
    var year = yearOfEra + era * 400;

    const dayOfYear = dayOfEra - (365 * yearOfEra + @divFloor(yearOfEra, 4) - @divFloor(yearOfEra, 100));

    const monthPrime = @divFloor((5 * dayOfYear + 2), 153);

    const day = dayOfYear - @divFloor((153 * monthPrime + 2), 5) + 1;
    var month = monthPrime + 3;
    if (month > 12) {
        month -= 1;
        year += 1;
    }

    var buffer: [20]u8 = undefined;

    const string = try std.fmt.bufPrint(buffer[0..], "{d:0>4}-{d:0>2}-{d:0>2}T{d:0>2}:{d:0>2}:{d:0>2}Z", .{ @as(u32, @intCast(year)), @as(u32, @intCast(month)), @as(u32, @intCast(day)), @as(u32, @intCast(hours)), @as(u32, @intCast(minutes)), @as(u32, @intCast(seconds)) });
    return string;
}

fn calculateDays(year: i64, month: i64, day: i64) i64 {
    var y = year;
    if (month < 3) y -= 1;

    const era = @divFloor(y, 400);
    const yearOfEra = y - era * 400;

    const monthPrime: i64 = if (month > 2) month - 3 else month + 9;
    const dayOfYear = @divFloor((153 * monthPrime + 2), 5) + day - 1;

    const dayOfEra = yearOfEra * 365 + @divFloor(yearOfEra, 4) - @divFloor(yearOfEra, 100) + dayOfYear;

    return era * 146_097 + dayOfEra - 719_468;
}

test "Time conversion test" {
    const testTime = Self.init(1_763_483_074);
    const timeString = try testTime.ISO8601();

    try std.testing.expectStringStartsWith(timeString, "2025-11-18T16:24:34Z");

    const otherWayTime = try Self.parseISO8601("2025-11-18T16:24:34Z");

    try std.testing.expectEqual(1_763_483_074, otherWayTime.timestamp);
}

const std = @import("std");
const Build = std.Build;
const Module = Build.Module;
const Step = Build.Step;
const Import = Module.Import;

pub fn build(b: *Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    const util_mod = b.addModule("util", .{
        .root_source_file = b.path("src/util.zig"),
        .target = target,
        .optimize = optimize,
    });

    const data_mod = b.addModule("data", .{
        .root_source_file = b.path("src/data/data.zig"),
        .target = target,
        .optimize = optimize,
        .imports = &.{.{ .name = "util", .module = util_mod }},
    });

    const server_mod = b.addModule("server", .{
        .root_source_file = b.path("src/server/server.zig"),
        .target = target,
        .optimize = optimize,
        .imports = &.{.{ .name = "data", .module = data_mod }},
    });

    const exe_mod = b.createModule(.{ .root_source_file = b.path("src/main.zig"), .target = target, .optimize = optimize, .imports = &.{
        .{ .name = "util", .module = util_mod },
        .{ .name = "data", .module = data_mod },
        .{ .name = "server", .module = server_mod },
    } });

    const exe = b.addExecutable(.{
        .name = "money",
        .root_module = exe_mod,
    });

    b.installArtifact(exe);
    const run_cmd = b.addRunArtifact(exe);
    run_cmd.step.dependOn(b.getInstallStep());
    if (b.args) |args| {
        run_cmd.addArgs(args);
    }
    const run_step = b.step("run", "Run the app");
    run_step.dependOn(&run_cmd.step);

    const exe_unit_tests = b.addTest(.{
        .root_module = exe_mod,
    });
    const run_exe_unit_tests = b.addRunArtifact(exe_unit_tests);
    const data_unit_tests = b.addTest(.{
        .root_module = data_mod,
    });
    const run_data_unit_tests = b.addRunArtifact(data_unit_tests);
    const test_step = b.step("test", "Run unit tests");

    test_step.dependOn(&run_exe_unit_tests.step);
    test_step.dependOn(&run_data_unit_tests.step);
}

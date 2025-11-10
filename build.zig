const std = @import("std");
const Build = std.Build;
const Module = Build.Module;
const Step = Build.Step;
const Import = Module.Import;

const ModuleData = struct {
    name: []const u8,
    path: []const u8,
    imports: []const []const u8,
    importAmount: usize,
};

// Every new module has to be added here. Make sure the order is correct. Modules can only import 
// other modules when they are above
const modulesToUse = [_]ModuleData{
    ModuleData{ .name = "util", .path = "src/util.zig", .imports = &.{}, .importAmount = 0 },
    ModuleData{ .name = "data", .path = "src/data/data.zig", .imports = &.{"util"}, .importAmount = 1 },
    ModuleData{ .name = "server", .path = "src/server/server.zig", .imports = &.{"data"}, .importAmount = 1 },
};

pub fn build(b: *Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    //######################### MODULES ######################################################//
    var modules: [modulesToUse.len]*Module = undefined;
    var imports: [modulesToUse.len]Import = undefined;

    for (0.., modulesToUse) |index, moduleData| {
        const new_mod = b.addModule(moduleData.name, .{
            .root_source_file = b.path(moduleData.path),
            .target = target,
            .optimize = optimize,
        });

        for (0..modules.len) |i| {
            const name = modulesToUse[i].name;
            for (moduleData.imports) |importName| {
                if (!std.mem.startsWith(u8, name, importName)) continue;
                const mod = modules[i];
                new_mod.addImport(importName, mod);
                break;
            }
        }
        modules[index] = new_mod;
        imports[index] = .{ .name = moduleData.name, .module = new_mod };
    }

    //######################### CREATING EXE ######################################################//
    const exe_mod = b.createModule(.{ .root_source_file = b.path("src/main.zig"), .target = target, .optimize = optimize, .imports = &imports });

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

    //######################### TESTS ######################################################//
    const exe_unit_tests = b.addTest(.{
        .root_module = exe_mod,
    });
    const run_exe_unit_tests = b.addRunArtifact(exe_unit_tests);
    const test_step = b.step("test", "Run unit tests");

    test_step.dependOn(&run_exe_unit_tests.step);

    // Adding all the module tests
    for (modules) |module| {
        const mod_tests = b.addTest(.{
            .root_module = module,
        });
        const run_mod_tests = b.addRunArtifact(mod_tests);
        test_step.dependOn(&run_mod_tests.step);
    }
}

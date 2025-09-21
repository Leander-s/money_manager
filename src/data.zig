budget: f32,
last_budget: f32,

const Self = @This();

pub fn init(budget: f32, last_budget: f32) Self {
    return .{ .budget = budget, .last_budget = last_budget };
}

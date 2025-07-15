-- Fitness Levels Table
CREATE TABLE fitness_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Fitness Goals Table
CREATE TABLE fitness_goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    category VARCHAR(50),
    icon_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Add fitness_level_id to users table
ALTER TABLE users ADD COLUMN fitness_level_id UUID REFERENCES fitness_levels(id) ON DELETE SET NULL;

-- User Fitness Goals Junction Table
CREATE TABLE user_fitness_goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    fitness_goal_id UUID NOT NULL,
    priority INTEGER DEFAULT 0,
    target_date TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (fitness_goal_id) REFERENCES fitness_goals(id) ON DELETE CASCADE,
    UNIQUE(user_id, fitness_goal_id)
);

-- Create indexes
CREATE INDEX idx_fitness_levels_name ON fitness_levels(name);
CREATE INDEX idx_fitness_levels_sort_order ON fitness_levels(sort_order);
CREATE INDEX idx_fitness_goals_name ON fitness_goals(name);
CREATE INDEX idx_fitness_goals_category ON fitness_goals(category);
CREATE INDEX idx_users_fitness_level_id ON users(fitness_level_id);
CREATE INDEX idx_user_fitness_goals_user_id ON user_fitness_goals(user_id);
CREATE INDEX idx_user_fitness_goals_fitness_goal_id ON user_fitness_goals(fitness_goal_id);
CREATE INDEX idx_user_fitness_goals_priority ON user_fitness_goals(priority);

-- Insert default fitness levels
INSERT INTO fitness_levels (name, description, sort_order) VALUES
    ('Beginner', 'New to fitness or returning after a long break', 1),
    ('Intermediate', 'Regular exercise experience with good form knowledge', 2),
    ('Advanced', 'Experienced athlete with years of consistent training', 3),
    ('Elite', 'Competitive athlete or professional level fitness', 4);

-- Insert default fitness goals
INSERT INTO fitness_goals (name, description, category, icon_name) VALUES
    ('Weight Loss', 'Reduce body weight and body fat percentage', 'body_composition', 'scale'),
    ('Muscle Gain', 'Build lean muscle mass and increase strength', 'body_composition', 'dumbbell'),
    ('Endurance', 'Improve cardiovascular fitness and stamina', 'performance', 'running'),
    ('Strength', 'Increase maximum strength and power output', 'performance', 'weight'),
    ('Flexibility', 'Improve range of motion and mobility', 'wellness', 'stretch'),
    ('General Fitness', 'Overall health and wellness improvement', 'wellness', 'heart'),
    ('Athletic Performance', 'Sport-specific performance enhancement', 'performance', 'trophy'),
    ('Rehabilitation', 'Recover from injury or medical condition', 'wellness', 'medical'),
    ('Body Recomposition', 'Simultaneous fat loss and muscle gain', 'body_composition', 'transform'),
    ('Stress Relief', 'Mental health and stress management through exercise', 'wellness', 'mindfulness');
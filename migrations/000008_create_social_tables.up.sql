-- Shared Workouts Table
CREATE TABLE shared_workouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_id UUID NOT NULL,
    shared_by_id UUID NOT NULL,
    shared_with_id UUID NOT NULL,
    permission VARCHAR(20) DEFAULT 'view',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE,
    FOREIGN KEY (shared_by_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (shared_with_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Workout Comments Table
CREATE TABLE workout_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shared_workout_id UUID NOT NULL,
    user_id UUID NOT NULL,
    parent_id UUID,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (shared_workout_id) REFERENCES shared_workouts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES workout_comments(id) ON DELETE CASCADE
);

-- Workout Comment Reactions Table
CREATE TABLE workout_comment_reactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    comment_id UUID NOT NULL,
    user_id UUID NOT NULL,
    reaction VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (comment_id) REFERENCES workout_comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_shared_workouts_workout_id ON shared_workouts(workout_id);
CREATE INDEX idx_shared_workouts_shared_by_id ON shared_workouts(shared_by_id);
CREATE INDEX idx_shared_workouts_shared_with_id ON shared_workouts(shared_with_id);
CREATE INDEX idx_shared_workouts_permission ON shared_workouts(permission);
CREATE INDEX idx_workout_comments_shared_workout_id ON workout_comments(shared_workout_id);
CREATE INDEX idx_workout_comments_user_id ON workout_comments(user_id);
CREATE INDEX idx_workout_comments_parent_id ON workout_comments(parent_id);
CREATE INDEX idx_workout_comment_reactions_comment_id ON workout_comment_reactions(comment_id);
CREATE INDEX idx_workout_comment_reactions_user_id ON workout_comment_reactions(user_id);
CREATE INDEX idx_workout_comment_reactions_reaction ON workout_comment_reactions(reaction);

-- Unique constraint to prevent duplicate reactions from the same user on the same comment
CREATE UNIQUE INDEX idx_workout_comment_reactions_user_comment_unique ON workout_comment_reactions(comment_id, user_id);
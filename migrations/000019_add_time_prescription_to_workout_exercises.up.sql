ALTER TABLE workout_exercises
  ADD COLUMN IF NOT EXISTS prescription varchar(20) NOT NULL DEFAULT 'reps',
  ADD COLUMN IF NOT EXISTS target_duration_sec integer NOT NULL DEFAULT 0;

ALTER TABLE workout_exercises
    ADD CONSTRAINT chk_we_prescription_valid
        CHECK (prescription IN ('reps','time'));

-- Revert back to TEXT-based user_id with FK to cognito_sub
DROP TABLE IF EXISTS reading_goal CASCADE;
DROP TABLE IF EXISTS reading_day CASCADE;

-- Recreate reading_day with user_id as TEXT
CREATE TABLE IF NOT EXISTS reading_day (
  user_id       TEXT NOT NULL,
  reading_date  DATE NOT NULL,
  pages_total   INTEGER NOT NULL,
  streak_days   INTEGER NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT reading_day_pk PRIMARY KEY (user_id, reading_date),
  CONSTRAINT reading_day_pages_total_chk CHECK (pages_total > 0),
  CONSTRAINT reading_day_streak_days_chk CHECK (streak_days >= 0),
  CONSTRAINT reading_day_user_fk FOREIGN KEY (user_id) REFERENCES users(cognito_sub) ON DELETE CASCADE
);

CREATE INDEX idx_reading_day_user_date ON reading_day(user_id, reading_date DESC);

-- Recreate reading_goal with user_id as TEXT
CREATE TABLE IF NOT EXISTS reading_goal (
  id          SERIAL PRIMARY KEY,
  user_id     TEXT NOT NULL,
  daily_pages INTEGER NOT NULL CHECK (daily_pages > 0 AND daily_pages <= 5000),
  start_date  TIMESTAMPTZ NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT reading_goal_user_start_unique UNIQUE (user_id, start_date),
  CONSTRAINT reading_goal_user_fk FOREIGN KEY (user_id) REFERENCES users(cognito_sub) ON DELETE CASCADE
);

CREATE INDEX idx_reading_goal_user_start_date ON reading_goal(user_id, start_date DESC);

-- Recreate updated_at triggers
CREATE OR REPLACE FUNCTION rc_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_reading_goal_updated_at ON reading_goal;
CREATE TRIGGER trg_reading_goal_updated_at
BEFORE UPDATE ON reading_goal
FOR EACH ROW
EXECUTE FUNCTION rc_set_updated_at();

DROP TRIGGER IF EXISTS trg_reading_day_updated_at ON reading_day;
CREATE TRIGGER trg_reading_day_updated_at
BEFORE UPDATE ON reading_day
FOR EACH ROW
EXECUTE FUNCTION rc_set_updated_at();

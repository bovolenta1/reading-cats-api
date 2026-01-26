ALTER TABLE IF EXISTS reading_goal RENAME TO reading_goal_old;

CREATE TABLE IF NOT EXISTS reading_goal (
  id          SERIAL PRIMARY KEY,
  user_id     TEXT NOT NULL,
  daily_pages INTEGER NOT NULL CHECK (daily_pages > 0 AND daily_pages <= 5000),
  start_date  TIMESTAMPTZ NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT reading_goal_user_start_unique UNIQUE (user_id, start_date),
  CONSTRAINT reading_goal_fk FOREIGN KEY (user_id) REFERENCES users(cognito_sub) ON DELETE CASCADE
);

CREATE INDEX idx_reading_goal_user_start_date ON reading_goal(user_id, start_date DESC);

INSERT INTO reading_goal (user_id, daily_pages, start_date, created_at)
SELECT user_id, daily_pages, (created_at AT TIME ZONE 'UTC')::date::timestamptz, created_at
FROM reading_goal_old
ON CONFLICT (user_id, start_date) DO NOTHING;

DROP TABLE IF EXISTS reading_goal_old;


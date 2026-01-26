-- Rollback: revert reading_goal table to original schema
CREATE TABLE IF NOT EXISTS reading_goal_old (
  user_id     TEXT PRIMARY KEY,
  daily_pages INTEGER NOT NULL DEFAULT 5,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT reading_goal_old_daily_pages_chk CHECK (daily_pages > 0 AND daily_pages <= 5000)
);

-- Migrate data back (use the latest goal per user)
INSERT INTO reading_goal_old (user_id, daily_pages, created_at, updated_at)
SELECT DISTINCT ON (user_id) user_id, daily_pages, created_at, now()
FROM reading_goal
ORDER BY user_id, start_date DESC
ON CONFLICT (user_id) DO NOTHING;

-- Drop new table
DROP TABLE IF EXISTS reading_goal;

-- Rename back
ALTER TABLE IF EXISTS reading_goal_old RENAME TO reading_goal;


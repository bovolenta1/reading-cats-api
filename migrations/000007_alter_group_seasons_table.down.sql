-- Revert: add edit_window_minutes back and rename started_at to starts_at
ALTER TABLE group_seasons
  ADD COLUMN edit_window_minutes smallint NOT NULL DEFAULT 15;

ALTER TABLE group_seasons
  RENAME COLUMN started_at TO starts_at;

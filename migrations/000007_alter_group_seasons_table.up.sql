-- Rename starts_at to started_at and remove edit_window_minutes
ALTER TABLE group_seasons
  RENAME COLUMN starts_at TO started_at;

ALTER TABLE group_seasons
  DROP COLUMN edit_window_minutes;

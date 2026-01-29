-- Revert: remove UUID PK and restore old structure
-- Rename index back
ALTER INDEX IF EXISTS idx_user_checkins_user_date RENAME TO idx_reading_day_user_date;

-- Remove the unique constraint
ALTER TABLE user_checkins DROP CONSTRAINT IF EXISTS user_checkins_user_date_key;

-- Remove UUID id column
ALTER TABLE user_checkins DROP COLUMN id;

-- Rename columns back
ALTER TABLE user_checkins RENAME COLUMN local_date TO reading_date;

-- Recreate the old composite PK
ALTER TABLE user_checkins ADD CONSTRAINT reading_day_pk PRIMARY KEY (user_id, reading_date);

-- Rename table back
ALTER TABLE user_checkins RENAME TO reading_day;

-- Rename reading_day to user_checkins and add UUID PK
ALTER TABLE reading_day RENAME TO user_checkins;

-- Drop the old composite PK (user_id, reading_date)
ALTER TABLE user_checkins DROP CONSTRAINT reading_day_pk;

-- Rename columns to match new semantics
ALTER TABLE user_checkins RENAME COLUMN reading_date TO local_date;

-- Add new UUID PK column
ALTER TABLE user_checkins ADD COLUMN id uuid PRIMARY KEY DEFAULT gen_random_uuid();

-- Add back the unique constraint on (user_id, local_date)
ALTER TABLE user_checkins ADD CONSTRAINT user_checkins_user_date_key UNIQUE (user_id, local_date);

-- Rename the index to match new table name
ALTER INDEX idx_reading_day_user_date RENAME TO idx_user_checkins_user_date;

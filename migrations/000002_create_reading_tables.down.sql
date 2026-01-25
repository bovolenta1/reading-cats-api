DROP TRIGGER IF EXISTS trg_reading_day_updated_at ON reading_day;
DROP TRIGGER IF EXISTS trg_reading_goal_updated_at ON reading_goal;

DROP TABLE IF EXISTS reading_day;
DROP TABLE IF EXISTS reading_goal;

DROP FUNCTION IF EXISTS rc_set_updated_at();

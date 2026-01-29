-- Drop triggers
DROP TRIGGER IF EXISTS set_user_checkins_updated_at ON user_checkins;
DROP TRIGGER IF EXISTS set_group_seasons_updated_at ON group_seasons;
DROP TRIGGER IF EXISTS set_groups_updated_at ON groups;

-- Drop tables (in reverse dependency order)
DROP TABLE IF EXISTS group_checkins;
DROP TABLE IF EXISTS group_seasons;
DROP TABLE IF EXISTS group_members;
DROP TABLE IF EXISTS groups;

-- Drop trigger function
DROP FUNCTION IF EXISTS set_updated_at();

-- Drop enums
DROP TYPE IF EXISTS group_metric;
DROP TYPE IF EXISTS group_season_status;
DROP TYPE IF EXISTS group_member_role;
DROP TYPE IF EXISTS group_visibility;

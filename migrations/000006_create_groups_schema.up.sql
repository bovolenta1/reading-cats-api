-- Create enums
CREATE TYPE group_visibility AS ENUM ('INVITE_ONLY', 'PUBLIC_SOON', 'FOUNDERS');
CREATE TYPE group_member_role AS ENUM ('ADMIN', 'MEMBER');
CREATE TYPE group_season_status AS ENUM ('DRAFT', 'ACTIVE', 'ENDED');
CREATE TYPE group_metric AS ENUM ('CHECKINS_PER_DAY', 'PAGES_SOON', 'MINUTES_SOON');

-- Create trigger function for updated_at
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create groups table
CREATE TABLE groups (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name varchar(30) NOT NULL,
  icon_id varchar NOT NULL,
  visibility group_visibility NOT NULL DEFAULT 'INVITE_ONLY',
  created_by_user_id uuid NOT NULL REFERENCES users(id),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_groups_created_by_user_id ON groups(created_by_user_id);

-- Create trigger for groups updated_at
CREATE TRIGGER set_groups_updated_at
  BEFORE UPDATE ON groups
  FOR EACH ROW
  EXECUTE FUNCTION set_updated_at();

-- Create group_members table
CREATE TABLE group_members (
  group_id uuid NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role group_member_role NOT NULL DEFAULT 'MEMBER',
  joined_at timestamptz NOT NULL DEFAULT now(),
  left_at timestamptz NULL,
  is_active boolean NOT NULL DEFAULT true,
  PRIMARY KEY (group_id, user_id)
);

CREATE INDEX idx_group_members_user_id ON group_members(user_id);
CREATE INDEX idx_group_members_group_id_active ON group_members(group_id, is_active);

-- Create group_seasons table
CREATE TABLE group_seasons (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  group_id uuid NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
  status group_season_status NOT NULL DEFAULT 'DRAFT',
  starts_at timestamptz NULL,
  ends_at timestamptz NULL,
  timezone varchar NOT NULL DEFAULT 'America/Sao_Paulo',
  metric group_metric NOT NULL DEFAULT 'CHECKINS_PER_DAY',
  edit_window_minutes smallint NOT NULL DEFAULT 15,
  created_by_user_id uuid NOT NULL REFERENCES users(id),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_group_seasons_group_id ON group_seasons(group_id);
CREATE INDEX idx_group_seasons_group_id_status ON group_seasons(group_id, status);

-- Partial unique index: only one ACTIVE season per group
CREATE UNIQUE INDEX idx_group_seasons_one_active_per_group
  ON group_seasons(group_id)
  WHERE status = 'ACTIVE';

-- Create trigger for group_seasons updated_at
CREATE TRIGGER set_group_seasons_updated_at
  BEFORE UPDATE ON group_seasons
  FOR EACH ROW
  EXECUTE FUNCTION set_updated_at();

-- Create group_checkins table
CREATE TABLE group_checkins (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  group_id uuid NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
  season_id uuid NOT NULL REFERENCES group_seasons(id) ON DELETE CASCADE,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  user_checkin_id uuid NOT NULL REFERENCES user_checkins(id) ON DELETE CASCADE,
  local_date date NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (group_id, season_id, user_id, local_date)
);

-- Add trigger for user_checkins if it has updated_at (rename migration adds it)
CREATE TRIGGER set_user_checkins_updated_at
  BEFORE UPDATE ON user_checkins
  FOR EACH ROW
  EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_group_checkins_group_season_date
  ON group_checkins(group_id, season_id, local_date DESC);

CREATE INDEX idx_group_checkins_season_user
  ON group_checkins(season_id, user_id);

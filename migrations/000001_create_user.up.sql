CREATE TABLE IF NOT EXISTS users (
  id            UUID PRIMARY KEY,
  cognito_sub   TEXT NOT NULL UNIQUE,
  email         TEXT,
  display_name  TEXT,
  avatar_url    TEXT,
  profile_source TEXT NOT NULL DEFAULT 'idp', -- 'idp' | 'user'
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_cognito_sub ON users (cognito_sub);
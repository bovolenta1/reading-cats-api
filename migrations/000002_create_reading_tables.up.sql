-- Goal (meta diária)
CREATE TABLE IF NOT EXISTS reading_goal (
  user_id     TEXT PRIMARY KEY,
  daily_pages INTEGER NOT NULL DEFAULT 5,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT reading_goal_daily_pages_chk CHECK (daily_pages > 0 AND daily_pages <= 5000)
);

-- Agregado por dia (1 linha por user + data)
CREATE TABLE IF NOT EXISTS reading_day (
  user_id       TEXT NOT NULL,
  reading_date  DATE NOT NULL,
  pages_total   INTEGER NOT NULL,
  streak_days   INTEGER NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT reading_day_pk PRIMARY KEY (user_id, reading_date),
  CONSTRAINT reading_day_pages_total_chk CHECK (pages_total >= 0),
  CONSTRAINT reading_day_streak_days_chk CHECK (streak_days >= 0)
);

-- (se quiser permitir 0 em algum cenário futuro, pode remover)
ALTER TABLE reading_day
  ADD CONSTRAINT reading_day_pages_total_positive_chk CHECK (pages_total > 0);

-- Trigger genérico de updated_at (escopado pra não conflitar com outros)
CREATE OR REPLACE FUNCTION rc_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_reading_goal_updated_at ON reading_goal;
CREATE TRIGGER trg_reading_goal_updated_at
BEFORE UPDATE ON reading_goal
FOR EACH ROW
EXECUTE FUNCTION rc_set_updated_at();

DROP TRIGGER IF EXISTS trg_reading_day_updated_at ON reading_day;
CREATE TRIGGER trg_reading_day_updated_at
BEFORE UPDATE ON reading_day
FOR EACH ROW
EXECUTE FUNCTION rc_set_updated_at();

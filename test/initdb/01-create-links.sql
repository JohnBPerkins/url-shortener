CREATE TABLE IF NOT EXISTS public.links (
  code       TEXT PRIMARY KEY,
  url        TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '24 hours'
);

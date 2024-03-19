CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE snippets (
  id UUID PRIMARY KEY default uuid_generate_v4(),
  created_at timestamp with time zone not null default now(),
  value text default ''::text,
  name text default ''::text
);

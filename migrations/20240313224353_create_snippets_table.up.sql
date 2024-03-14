CREATE TABLE snippets (
  id UUID PRIMARY KEY,
  created_at timestamp with time zone not null default now(),
  value text default ''::text,
  name text default ''::text
);

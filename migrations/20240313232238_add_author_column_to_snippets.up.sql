ALTER TABLE snippets
  ADD COLUMN author UUID,
  ADD CONSTRAINT fk_snippets_author
  FOREIGN KEY (author) REFERENCES users(id)
  ON DELETE CASCADE;

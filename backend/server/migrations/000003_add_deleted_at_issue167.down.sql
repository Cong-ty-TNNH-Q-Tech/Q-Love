ALTER TABLE card_transactions DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE wall_of_shames DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE vibe_checks DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE card_steals DROP COLUMN IF EXISTS deleted_at;

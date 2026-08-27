DROP INDEX IF EXISTS idx_wallet_transactions_deleted_at;
ALTER TABLE wallet_transactions DROP COLUMN IF EXISTS deleted_at;

DROP INDEX IF EXISTS idx_user_wallets_deleted_at;
ALTER TABLE user_wallets DROP COLUMN IF EXISTS deleted_at;

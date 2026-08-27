ALTER TABLE user_wallets ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;
CREATE INDEX IF NOT EXISTS idx_user_wallets_deleted_at ON user_wallets(deleted_at);

ALTER TABLE wallet_transactions ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;
CREATE INDEX IF NOT EXISTS idx_wallet_transactions_deleted_at ON wallet_transactions(deleted_at);

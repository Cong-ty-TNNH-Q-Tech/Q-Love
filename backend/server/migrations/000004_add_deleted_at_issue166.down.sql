ALTER TABLE blind_auctions DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE auction_bids DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE ex_ratings DROP COLUMN IF EXISTS deleted_at;

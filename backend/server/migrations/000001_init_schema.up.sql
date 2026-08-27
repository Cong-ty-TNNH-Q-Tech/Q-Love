-- Enable PostGIS extension for spatial data
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Core & Account
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(50) NOT NULL,
    dob DATE NOT NULL,
    gender VARCHAR(10),
    location GEOMETRY(Point, 4326),
    level INT DEFAULT 1,
    is_shadowbanned BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_location ON users USING GIST(location);

CREATE TABLE user_wallets (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    balance NUMERIC(15,2) DEFAULT 0.00,
    hold_balance NUMERIC(15,2) DEFAULT 0.00
);

CREATE TABLE wallet_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    amount NUMERIC NOT NULL,
    type VARCHAR(50),
    reference_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE user_premiums (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    is_active BOOLEAN DEFAULT false,
    expires_at TIMESTAMP,
    free_cancel_left INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Engagement & Matchmaking
CREATE TABLE matches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user1_id UUID REFERENCES users(id),
    user2_id UUID REFERENCES users(id),
    streak_score INT DEFAULT 0,
    highest_streak_score INT DEFAULT 0,
    island_level INT DEFAULT 1,
    last_interaction_at TIMESTAMP,
    UNIQUE(user1_id, user2_id)
);

CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID REFERENCES matches(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id),
    type VARCHAR(20),
    content TEXT,
    blur_url TEXT,
    blur_level INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ex_ratings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reviewer_id UUID REFERENCES users(id),
    target_id UUID REFERENCES users(id),
    rating_score INT CHECK (rating_score >= 1 AND rating_score <= 5),
    tags TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- O2O & Gamification (Khế Ước, Bản Đồ, Tòa Án)
CREATE TABLE dating_contracts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_a_id UUID REFERENCES users(id),
    user_b_id UUID REFERENCES users(id),
    deposit_amount NUMERIC NOT NULL,
    status VARCHAR(20),
    cancelled_by_id UUID REFERENCES users(id),
    totp_secret VARCHAR(255),
    appointment_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE clans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE,
    leader_id UUID REFERENCES users(id) ON DELETE SET NULL,
    weekly_score INT DEFAULT 0,
    campfire_streak INT DEFAULT 0,
    daily_active_members INT DEFAULT 0,
    last_campfire_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE clan_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    clan_id UUID REFERENCES clans(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id),
    role VARCHAR(20) DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(clan_id, user_id)
);

CREATE TABLE landmarks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100),
    location GEOMETRY(Point, 4326),
    radius_meters INT DEFAULT 200,
    current_owner_clan_id UUID REFERENCES clans(id) ON DELETE SET NULL
);

CREATE INDEX idx_landmarks_location ON landmarks USING GIST(location);

CREATE TABLE court_cases (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plaintiff_id UUID REFERENCES users(id),
    defendant_id UUID REFERENCES users(id),
    reason VARCHAR(100),
    status VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE court_votes (
    case_id UUID REFERENCES court_cases(id) ON DELETE CASCADE,
    juror_id UUID REFERENCES users(id),
    vote VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (case_id, juror_id)
);

-- Kinh Tế Ảo (Chợ Thẻ Bài Profile)
CREATE TABLE card_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    current_price NUMERIC(15,2) DEFAULT 100,
    total_cards INT DEFAULT 1000,
    available_cards INT DEFAULT 1000,
    match_count_cached INT DEFAULT 0,
    locket_count_cached INT DEFAULT 0,
    clan_upvote_cached INT DEFAULT 0,
    court_penalty_cached INT DEFAULT 0
);

CREATE TABLE card_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    collector_id UUID REFERENCES users(id),
    target_user_id UUID REFERENCES users(id),
    type VARCHAR(10),
    quantity INT,
    price_at_transaction NUMERIC,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Notifications & Violations
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    type VARCHAR(50) NOT NULL,
    payload TEXT,
    status VARCHAR(20) DEFAULT 'sent',
    reference_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user_status ON notifications(user_id, status);

CREATE TABLE user_violations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    type VARCHAR(50) NOT NULL,
    reason TEXT,
    is_active BOOLEAN DEFAULT true,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_violations_active ON user_violations(user_id, type, expires_at);

-- Gamification Đột Phá
CREATE TABLE bounties (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    description TEXT,
    reward_amount INT,
    status VARCHAR(20) DEFAULT 'open',
    winner_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE blind_auctions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    target_user_id UUID REFERENCES users(id),
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    current_highest_bid INT DEFAULT 0,
    winner_id UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(20) DEFAULT 'active'
);

CREATE TABLE auction_bids (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    auction_id UUID REFERENCES blind_auctions(id) ON DELETE CASCADE,
    bidder_id UUID REFERENCES users(id),
    bid_amount INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE wall_of_shames (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    reason TEXT,
    tomatoes_thrown INT DEFAULT 0,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE vibe_checks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user1_id UUID REFERENCES users(id),
    user2_id UUID REFERENCES users(id),
    track_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE wingman_referrals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wingman_id UUID REFERENCES users(id),
    target1_id UUID REFERENCES users(id),
    target2_id UUID REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE card_steals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    attacker_id UUID REFERENCES users(id),
    defender_id UUID REFERENCES users(id),
    target_card_id UUID REFERENCES users(id),
    result VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

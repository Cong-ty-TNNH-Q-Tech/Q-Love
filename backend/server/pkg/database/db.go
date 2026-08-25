// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package database

import (
	"fmt"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/config"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	GormOpen     = gorm.Open
	PostgresOpen = postgres.Open
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	if cfg.DatabaseDSN == "" {
		return nil, fmt.Errorf("DatabaseDSN is required in configuration")
	}
	if cfg.DatabaseDSN == "skip" {
		logger.Log.Info("Skipping database connection for testing")
		return nil, nil
	}

	db, err := GormOpen(PostgresOpen(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		logger.Log.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Log.Error("Failed to get database connection pool", zap.Error(err))
		return nil, err
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Log.Info("Successfully connected to PostgreSQL database")

	// Create Spiritual Match SQL function
	err = db.Exec(`
		CREATE OR REPLACE FUNCTION calculate_spiritual_match_score(dob1 TIMESTAMP, dob2 TIMESTAMP) RETURNS INT AS $$
		DECLARE
			d1 INT := extract(day from dob1); m1 INT := extract(month from dob1);
			d2 INT := extract(day from dob2); m2 INT := extract(month from dob2);
			z1 VARCHAR; z2 VARCHAR;
			ele1 VARCHAR; ele2 VARCHAR;
			n1 INT := 0; n2 INT := 0;
			tmp_sum INT;
			date_str1 VARCHAR := to_char(dob1, 'DDMMYYYY');
			date_str2 VARCHAR := to_char(dob2, 'DDMMYYYY');
			score INT := 40;
			i INT;
		BEGIN
			IF (m1=3 AND d1>=21) OR (m1=4 AND d1<20) THEN z1 := 'Aries'; ele1 := 'Fire';
			ELSIF (m1=4 AND d1>=20) OR (m1=5 AND d1<21) THEN z1 := 'Taurus'; ele1 := 'Earth';
			ELSIF (m1=5 AND d1>=21) OR (m1=6 AND d1<22) THEN z1 := 'Gemini'; ele1 := 'Air';
			ELSIF (m1=6 AND d1>=22) OR (m1=7 AND d1<23) THEN z1 := 'Cancer'; ele1 := 'Water';
			ELSIF (m1=7 AND d1>=23) OR (m1=8 AND d1<23) THEN z1 := 'Leo'; ele1 := 'Fire';
			ELSIF (m1=8 AND d1>=23) OR (m1=9 AND d1<23) THEN z1 := 'Virgo'; ele1 := 'Earth';
			ELSIF (m1=9 AND d1>=23) OR (m1=10 AND d1<24) THEN z1 := 'Libra'; ele1 := 'Air';
			ELSIF (m1=10 AND d1>=24) OR (m1=11 AND d1<22) THEN z1 := 'Scorpio'; ele1 := 'Water';
			ELSIF (m1=11 AND d1>=22) OR (m1=12 AND d1<22) THEN z1 := 'Sagittarius'; ele1 := 'Fire';
			ELSIF (m1=12 AND d1>=22) OR (m1=1 AND d1<20) THEN z1 := 'Capricorn'; ele1 := 'Earth';
			ELSIF (m1=1 AND d1>=20) OR (m1=2 AND d1<19) THEN z1 := 'Aquarius'; ele1 := 'Air';
			ELSE z1 := 'Pisces'; ele1 := 'Water'; END IF;
			
			IF (m2=3 AND d2>=21) OR (m2=4 AND d2<20) THEN z2 := 'Aries'; ele2 := 'Fire';
			ELSIF (m2=4 AND d2>=20) OR (m2=5 AND d2<21) THEN z2 := 'Taurus'; ele2 := 'Earth';
			ELSIF (m2=5 AND d2>=21) OR (m2=6 AND d2<22) THEN z2 := 'Gemini'; ele2 := 'Air';
			ELSIF (m2=6 AND d2>=22) OR (m2=7 AND d2<23) THEN z2 := 'Cancer'; ele2 := 'Water';
			ELSIF (m2=7 AND d2>=23) OR (m2=8 AND d2<23) THEN z2 := 'Leo'; ele2 := 'Fire';
			ELSIF (m2=8 AND d2>=23) OR (m2=9 AND d2<23) THEN z2 := 'Virgo'; ele2 := 'Earth';
			ELSIF (m2=9 AND d2>=23) OR (m2=10 AND d2<24) THEN z2 := 'Libra'; ele2 := 'Air';
			ELSIF (m2=10 AND d2>=24) OR (m2=11 AND d2<22) THEN z2 := 'Scorpio'; ele2 := 'Water';
			ELSIF (m2=11 AND d2>=22) OR (m2=12 AND d2<22) THEN z2 := 'Sagittarius'; ele2 := 'Fire';
			ELSIF (m2=12 AND d2>=22) OR (m2=1 AND d2<20) THEN z2 := 'Capricorn'; ele2 := 'Earth';
			ELSIF (m2=1 AND d2>=20) OR (m2=2 AND d2<19) THEN z2 := 'Aquarius'; ele2 := 'Air';
			ELSE z2 := 'Pisces'; ele2 := 'Water'; END IF;
			
			IF z1 = z2 THEN score := score + 25;
			ELSIF ele1 = ele2 THEN score := score + 30;
			ELSIF (ele1='Fire' AND ele2='Air') OR (ele1='Air' AND ele2='Fire') THEN score := score + 20;
			ELSIF (ele1='Earth' AND ele2='Water') OR (ele1='Water' AND ele2='Earth') THEN score := score + 20;
			ELSE score := score + 5; END IF;
			
			FOR i IN 1..8 LOOP n1 := n1 + CAST(substring(date_str1 FROM i FOR 1) AS INT); END LOOP;
			WHILE n1 > 9 AND n1 != 11 AND n1 != 22 LOOP
				tmp_sum := 0;
				FOR i IN 1..length(CAST(n1 AS VARCHAR)) LOOP tmp_sum := tmp_sum + CAST(substring(CAST(n1 AS VARCHAR) FROM i FOR 1) AS INT); END LOOP;
				n1 := tmp_sum;
			END LOOP;
			
			FOR i IN 1..8 LOOP n2 := n2 + CAST(substring(date_str2 FROM i FOR 1) AS INT); END LOOP;
			WHILE n2 > 9 AND n2 != 11 AND n2 != 22 LOOP
				tmp_sum := 0;
				FOR i IN 1..length(CAST(n2 AS VARCHAR)) LOOP tmp_sum := tmp_sum + CAST(substring(CAST(n2 AS VARCHAR) FROM i FOR 1) AS INT); END LOOP;
				n2 := tmp_sum;
			END LOOP;
			
			IF n1 = n2 THEN score := score + 30;
			ELSIF (n1%2=0 AND n2%2=0) OR (n1%2!=0 AND n2%2!=0) THEN score := score + 20;
			ELSE score := score + 10; END IF;
			
			IF score > 100 THEN score := 100; END IF;
			
			RETURN score;
		END;
		$$ LANGUAGE plpgsql IMMUTABLE;
	`).Error
	if err != nil {
		logger.Log.Error("Failed to create calculate_spiritual_match_score function", zap.Error(err))
	}

	return db, nil
}

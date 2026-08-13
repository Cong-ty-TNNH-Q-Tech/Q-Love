// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package repository

import (
	"context"
	"errors"
	"testing"
)

func TestTransactionManager_WithTransaction(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	tm := NewTransactionManager(db)

	mock.ExpectBegin()
	mock.ExpectCommit()

	err = tm.WithTransaction(context.Background(), func(ctx context.Context) error {
		// Just return nil to trigger commit
		return nil
	})
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestTransactionManager_WithTransaction_Error(t *testing.T) {
	db, mock, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	tm := NewTransactionManager(db)

	mock.ExpectBegin()
	mock.ExpectRollback()

	err = tm.WithTransaction(context.Background(), func(ctx context.Context) error {
		return errors.New("custom error")
	})
	
	if err == nil || err.Error() != "custom error" {
		t.Errorf("Expected custom error, got %v", err)
	}
}

func TestGetDB(t *testing.T) {
	db, _, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	
	ctx := context.Background()
	resultDB := GetDB(ctx, db)
	if resultDB != db {
		t.Error("Expected original DB when no transaction in context")
	}
}

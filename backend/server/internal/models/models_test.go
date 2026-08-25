// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUser_BeforeCreate(t *testing.T) {
	u := &User{}
	err := u.BeforeCreate(&gorm.DB{})
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, u.ID)
}

func TestCourtCase_BeforeCreate(t *testing.T) {
	c := &CourtCase{}
	err := c.BeforeCreate(&gorm.DB{})
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, c.ID)
}

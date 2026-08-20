// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExRating_Tags(t *testing.T) {
	r := &ExRating{}
	
	// Test empty tags
	assert.Empty(t, r.GetTags())

	// Test SetTags and GetTags
	r.SetTags([]string{"#tốt", "#vui_vẻ"})
	assert.Equal(t, "#tốt,#vui_vẻ", r.TagsString)

	tags := r.GetTags()
	assert.Len(t, tags, 2)
	assert.Equal(t, "#tốt", tags[0])
	assert.Equal(t, "#vui_vẻ", tags[1])
}

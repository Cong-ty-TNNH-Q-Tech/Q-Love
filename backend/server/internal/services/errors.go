// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import "errors"

var (
	ErrInsufficientBalance   = errors.New("insufficient balance to throw a tomato")
	ErrReferralNotFound      = errors.New("referral not found")
	ErrUserNotInReferral     = errors.New("user is not part of this referral")
	ErrReferralNotPending    = errors.New("referral is no longer pending")
	ErrReferralExpired       = errors.New("referral link expired")
	ErrWingmanCannotReferSelf= errors.New("wingman cannot refer themselves")
)

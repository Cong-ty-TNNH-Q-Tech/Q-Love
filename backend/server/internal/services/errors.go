// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

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

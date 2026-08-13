// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package logger

import (
	"log"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger(env string, sentryDSN string) {
	var err error
	if env == "production" {
		Log, err = zap.NewProduction()
	} else {
		Log, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}

	if sentryDSN != "" {
		err = sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDSN,
			Environment:      env,
			TracesSampleRate: 1.0,
		})
		if err != nil {
			Log.Error("sentry.Init failed", zap.Error(err))
		} else {
			Log.Info("Sentry initialized successfully")
		}
	}
}

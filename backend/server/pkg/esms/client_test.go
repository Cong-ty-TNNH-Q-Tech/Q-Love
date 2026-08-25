// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package esms

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_SendOTP_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"CodeResult": "100",
			"CountRegenerate": 0,
			"SMSID": "123456"
		}`))
	}))
	defer mockServer.Close()

	c := &client{
		apiKey:    "test-key",
		secretKey: "test-secret",
		baseURL:   mockServer.URL,
		hc:        mockServer.Client(),
	}

	err := c.SendOTP(context.Background(), "0901234567", "123456")
	assert.NoError(t, err)
}

func TestClient_SendOTP_ErrorResponse(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"CodeResult": "99",
			"ErrorMessage": "Loi gui tin"
		}`))
	}))
	defer mockServer.Close()

	c := &client{
		apiKey:    "test-key",
		secretKey: "test-secret",
		baseURL:   mockServer.URL,
		hc:        mockServer.Client(),
	}

	err := c.SendOTP(context.Background(), "0901234567", "123456")
	assert.Error(t, err)
	assert.Equal(t, "ESMS error 99: Loi gui tin", err.Error())
}

func TestClient_SendOTP_HTTPError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	c := &client{
		apiKey:    "test-key",
		secretKey: "test-secret",
		baseURL:   mockServer.URL,
		hc:        mockServer.Client(),
	}

	err := c.SendOTP(context.Background(), "0901234567", "123456")
	assert.Error(t, err)
	assert.Equal(t, "ESMS returned HTTP 500", err.Error())
}

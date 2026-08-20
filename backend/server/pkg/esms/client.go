// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package esms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client interface {
	SendOTP(ctx context.Context, phone, otp string) error
}

type client struct {
	apiKey    string
	secretKey string
	baseURL   string
	hc        *http.Client
}

func NewClient(apiKey, secretKey string) Client {
	return &client{
		apiKey:    apiKey,
		secretKey: secretKey,
		baseURL:   "https://restapi.esms.vn/MainService.svc/json/SendMultipleMessage_V4_post_json",
		hc: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type ESMSRequest struct {
	APIKey    string `json:"ApiKey"`
	SecretKey string `json:"SecretKey"`
	Phone     string `json:"Phone"`
	Content   string `json:"Content"`
	SmsType   string `json:"SmsType"`
	Brandname string `json:"Brandname,omitempty"`
}

type ESMSResponse struct {
	CodeResult string `json:"CodeResult"`
	CountRegenerate int `json:"CountRegenerate"`
	SMSID      string `json:"SMSID"`
	ErrorMessage string `json:"ErrorMessage,omitempty"`
}

func (c *client) SendOTP(ctx context.Context, phone, otp string) error {
	reqBody := ESMSRequest{
		APIKey:    c.apiKey,
		SecretKey: c.secretKey,
		Phone:     phone,
		Content:   fmt.Sprintf("Q-Love: Ma OTP cua ban la %s. Vui long khong chia se ma nay cho bat ky ai. (Co hieu luc trong 120 giay)", otp),
		SmsType:   "2", // Type for CSKH / OTP
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ESMS returned HTTP %d", resp.StatusCode)
	}

	var resBody ESMSResponse
	if err := json.NewDecoder(resp.Body).Decode(&resBody); err != nil {
		return err
	}

	if resBody.CodeResult != "100" { // 100 means success according to ESMS docs
		return fmt.Errorf("ESMS error %s: %s", resBody.CodeResult, resBody.ErrorMessage)
	}

	return nil
}

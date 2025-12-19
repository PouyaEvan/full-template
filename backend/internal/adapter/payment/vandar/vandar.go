package vandar

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	VandarRequestURL = "https://ipg.vandar.io/api/v3/send"
	VandarVerifyURL  = "https://ipg.vandar.io/api/v3/verify"
	VandarStartURL   = "https://ipg.vandar.io/v3/"
)

type VandarAdapter struct {
	APIKey string
	Client *http.Client
}

func NewVandarAdapter(apiKey string) *VandarAdapter {
	return &VandarAdapter{
		APIKey: apiKey,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

type requestPayload struct {
	APIKey      string `json:"api_key"`
	Amount      int64  `json:"amount"`
	CallbackURL string `json:"callback_url"`
	Mobile      string `json:"mobile_number,omitempty"`
	Description string `json:"description,omitempty"`
}

type requestResponse struct {
	Status int      `json:"status"`
	Token  string   `json:"token"`
	Errors []string `json:"errors,omitempty"`
}

func (v *VandarAdapter) RequestPayment(ctx context.Context, amount int64, callbackURL string, description string, mobile string) (string, string, error) {
	payload := requestPayload{
		APIKey:      v.APIKey,
		Amount:      amount, // Vandar uses Rials usually, check docs if Toman
		CallbackURL: callbackURL,
		Mobile:      mobile,
		Description: description,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", VandarRequestURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := v.Client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result requestResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	if result.Status != 1 {
		return "", "", fmt.Errorf("vandar error: %v", result.Errors)
	}

	return VandarStartURL + result.Token, result.Token, nil
}

type verifyPayload struct {
	APIKey string `json:"api_key"`
	Token  string `json:"token"`
}

type verifyResponse struct {
	Status  int      `json:"status"`
	Amount  string   `json:"amount"`
	TransId string   `json:"transId"`
	Errors  []string `json:"errors,omitempty"`
}

func (v *VandarAdapter) VerifyPayment(ctx context.Context, token string, amount int64) (string, error) {
	payload := verifyPayload{
		APIKey: v.APIKey,
		Token:  token,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", VandarVerifyURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := v.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result verifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Status != 1 {
		return "", fmt.Errorf("verification failed: %v", result.Errors)
	}

	return result.TransId, nil
}

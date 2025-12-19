package zarinpal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	ZarinpalRequestURL = "https://api.zarinpal.com/pg/v4/payment/request.json"
	ZarinpalVerifyURL  = "https://api.zarinpal.com/pg/v4/payment/verify.json"
	ZarinpalStartURL   = "https://www.zarinpal.com/pg/StartPay/"
)

type ZarinpalAdapter struct {
	MerchantID string
	Client     *http.Client
}

func NewZarinpalAdapter(merchantID string) *ZarinpalAdapter {
	return &ZarinpalAdapter{
		MerchantID: merchantID,
		Client:     &http.Client{Timeout: 10 * time.Second},
	}
}

type requestPayload struct {
	MerchantID  string `json:"merchant_id"`
	Amount      int64  `json:"amount"`
	CallbackURL string `json:"callback_url"`
	Description string `json:"description"`
	Metadata    struct {
		Email string `json:"email,omitempty"`
	} `json:"metadata,omitempty"`
}

type requestResponse struct {
	Data struct {
		Code      int    `json:"code"`
		Message   string `json:"message"`
		Authority string `json:"authority"`
		FeeType   string `json:"fee_type"`
		Fee       int    `json:"fee"`
	} `json:"data"`
	Errors []interface{} `json:"errors"`
}

func (z *ZarinpalAdapter) RequestPayment(ctx context.Context, amount int64, callbackURL string, description string, email string) (string, string, error) {
	payload := requestPayload{
		MerchantID:  z.MerchantID,
		Amount:      amount,
		CallbackURL: callbackURL,
		Description: description,
	}
	payload.Metadata.Email = email

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", ZarinpalRequestURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.Client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result requestResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	if result.Data.Code != 100 {
		return "", "", fmt.Errorf("zarinpal error code: %d", result.Data.Code)
	}

	return ZarinpalStartURL + result.Data.Authority, result.Data.Authority, nil
}

type verifyPayload struct {
	MerchantID string `json:"merchant_id"`
	Amount     int64  `json:"amount"`
	Authority  string `json:"authority"`
}

type verifyResponse struct {
	Data struct {
		Code     int    `json:"code"`
		Message  string `json:"message"`
		CardHash string `json:"card_hash"`
		CardPan  string `json:"card_pan"`
		RefID    int    `json:"ref_id"`
		FeeType  string `json:"fee_type"`
		Fee      int    `json:"fee"`
	} `json:"data"`
	Errors []interface{} `json:"errors"`
}

func (z *ZarinpalAdapter) VerifyPayment(ctx context.Context, authority string, amount int64) (string, error) {
	payload := verifyPayload{
		MerchantID: z.MerchantID,
		Amount:     amount,
		Authority:  authority,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", ZarinpalVerifyURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result verifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// 100 = Success, 101 = Verified Already
	if result.Data.Code != 100 && result.Data.Code != 101 {
		return "", fmt.Errorf("verification failed code: %d", result.Data.Code)
	}

	return fmt.Sprintf("%d", result.Data.RefID), nil
}

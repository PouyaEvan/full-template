package senator

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

const SenatorAPIURL = "https://api.fast-creat.ir/sms"

type SenatorAdapter struct {
	APIKey     string
	TemplateID string
	Client     *http.Client
}

func NewSenatorAdapter() *SenatorAdapter {
	return &SenatorAdapter{
		APIKey:     os.Getenv("SENATOR_API_KEY"),
		TemplateID: os.Getenv("SENATOR_TEMPLATE_ID"),
		Client:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *SenatorAdapter) SendOTP(ctx context.Context, phoneNumber string, code string) error {
	// https://api.fast-creat.ir/sms?apikey=xxxxx&type=sms&code=xxxxx&phone=xxxxx&template=xxxxx

	params := url.Values{}
	params.Add("apikey", s.APIKey)
	params.Add("type", "sms")
	params.Add("code", code)
	params.Add("phone", phoneNumber)
	params.Add("template", s.TemplateID)

	reqURL := fmt.Sprintf("%s?%s", SenatorAPIURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return err
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("senator sms failed with status: %d", resp.StatusCode)
	}

	return nil
}

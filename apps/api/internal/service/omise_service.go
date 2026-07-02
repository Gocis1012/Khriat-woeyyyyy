package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const omiseAPIBase = "https://api.omise.co"

type omiseSource struct {
	ID string `json:"id"`
}

type omiseCharge struct {
	ID     string `json:"id"`
	Status string `json:"status"` // pending | successful | failed | expired
	Source struct {
		ScannableCode struct {
			Image struct {
				DownloadURI string `json:"download_uri"`
			} `json:"image"`
		} `json:"scannable_code"`
	} `json:"source"`
}

type omiseErrorResponse struct {
	Message string `json:"message"`
}

// OmiseQRCharge is the result of creating a PromptPay QR charge.
type OmiseQRCharge struct {
	ChargeID  string
	Status    string
	QRCodeURI string
}

type OmiseService struct {
	secretKey  string
	httpClient *http.Client
}

func NewOmiseService(secretKey string) *OmiseService {
	return &OmiseService{
		secretKey:  secretKey,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (s *OmiseService) doForm(ctx context.Context, path string, form url.Values, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, omiseAPIBase+path, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.secretKey, "")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var oerr omiseErrorResponse
		_ = json.Unmarshal(body, &oerr)
		if oerr.Message != "" {
			return fmt.Errorf("omise %s error: %s", path, oerr.Message)
		}
		return fmt.Errorf("omise %s error: status %d", path, resp.StatusCode)
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

// CreatePromptPayCharge creates an Omise Source + Charge for a PromptPay QR
// payment. amountTHB is in Thai Baht and is converted to satang (THB * 100)
// for the API call.
func (s *OmiseService) CreatePromptPayCharge(ctx context.Context, amountTHB float64) (*OmiseQRCharge, error) {
	amountSatang := int64(amountTHB * 100)

	sourceForm := url.Values{}
	sourceForm.Set("amount", strconv.FormatInt(amountSatang, 10))
	sourceForm.Set("currency", "thb")
	sourceForm.Set("type", "promptpay")

	var source omiseSource
	if err := s.doForm(ctx, "/sources", sourceForm, &source); err != nil {
		return nil, fmt.Errorf("create source: %w", err)
	}

	chargeForm := url.Values{}
	chargeForm.Set("amount", strconv.FormatInt(amountSatang, 10))
	chargeForm.Set("currency", "thb")
	chargeForm.Set("source", source.ID)

	var charge omiseCharge
	if err := s.doForm(ctx, "/charges", chargeForm, &charge); err != nil {
		return nil, fmt.Errorf("create charge: %w", err)
	}

	return &OmiseQRCharge{
		ChargeID:  charge.ID,
		Status:    charge.Status,
		QRCodeURI: charge.Source.ScannableCode.Image.DownloadURI,
	}, nil
}

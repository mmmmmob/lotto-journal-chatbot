package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type GLOValue struct {
	Round int    `json:"round"`
	Value string `json:"value"`
}

type GLOPrize struct {
	Price  string     `json:"price"`
	Number []GLOValue `json:"number"`
}

type GLOL6Data struct {
	First     GLOPrize `json:"first"`
	Second    GLOPrize `json:"second"`
	Third     GLOPrize `json:"third"`
	Fourth    GLOPrize `json:"fourth"`
	Fifth     GLOPrize `json:"fifth"`
	Last2     GLOPrize `json:"last2"`
	Last3f    GLOPrize `json:"last3f"`
	Last3b    GLOPrize `json:"last3b"`
	Near1     GLOPrize `json:"near1"`
}

type GLON3Data struct {
	Straight3 GLOPrize `json:"straight3"`
	Shuffle3  GLOPrize `json:"shuffle3"`
	Straight2 GLOPrize `json:"straight2"`
	Special   GLOPrize `json:"special"`
}

type GLOLotteryResponse struct {
	StatusMessage string `json:"statusMessage"`
	StatusCode    int    `json:"statusCode"`
	Status        bool   `json:"status"`
	Response      struct {
		Date string    `json:"date"`
		Data GLOL6Data `json:"data"`
		N3   GLON3Data `json:"n3"`
	} `json:"response"`
}

type GLOScheduleRequest struct {
	Year int    `json:"year"`
	Type string `json:"type"`
}

type PeriodResult struct {
	Date string `json:"date"`
}

type GLOScheduleResponse struct {
	StatusMessage string `json:"statusMessage"`
	StatusCode    int    `json:"statusCode"`
	Status        bool   `json:"status"`
	Response      struct {
		Result []PeriodResult `json:"result"`
	} `json:"response"`
}

type LotteryClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewLotteryClient(baseURL string) *LotteryClient {
	if baseURL == "" {
		baseURL = "https://www.glo.or.th"
	}
	return &LotteryClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// FetchLatestResult fetches the latest lottery results from GLO.
func (c *LotteryClient) FetchLatestResult(ctx context.Context) (*GLOLotteryResponse, error) {
	url := fmt.Sprintf("%s/api/lottery/getLatestLottery", c.baseURL)

	respBody, err := c.retry(ctx, func() ([]byte, error) {
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer([]byte("{}")))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		return io.ReadAll(resp.Body)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest result after retries: %w", err)
	}

	var lotteryResp GLOLotteryResponse
	if err := json.Unmarshal(respBody, &lotteryResp); err != nil {
		return nil, fmt.Errorf("decode GLO response: %w", err)
	}

	if !lotteryResp.Status {
		return nil, fmt.Errorf("GLO API returned status false: %s", lotteryResp.StatusMessage)
	}

	return &lotteryResp, nil
}

// FetchDrawSchedule fetches the draw dates schedule for the given year.
func (c *LotteryClient) FetchDrawSchedule(ctx context.Context, year int) ([]time.Time, error) {
	url := fmt.Sprintf("%s/api/lottery/getPeriodsByYear", c.baseURL)

	reqPayload := GLOScheduleRequest{
		Year: year,
		Type: "ALL",
	}
	reqBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("marshal schedule request: %w", err)
	}

	respBody, err := c.retry(ctx, func() ([]byte, error) {
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBytes))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		return io.ReadAll(resp.Body)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedule after retries: %w", err)
	}

	var scheduleResp GLOScheduleResponse
	if err := json.Unmarshal(respBody, &scheduleResp); err != nil {
		return nil, fmt.Errorf("decode GLO schedule response: %w", err)
	}

	if !scheduleResp.Status {
		return nil, fmt.Errorf("GLO API returned schedule status false: %s", scheduleResp.StatusMessage)
	}

	// Extract unique draw dates in chronological order
	uniqueDatesMap := make(map[string]bool)
	var uniqueDates []time.Time

	for _, res := range scheduleResp.Response.Result {
		if res.Date == "" {
			continue
		}
		if _, exists := uniqueDatesMap[res.Date]; !exists {
			uniqueDatesMap[res.Date] = true
			parsedDate, err := time.Parse("2006-01-02", res.Date)
			if err == nil {
				uniqueDates = append(uniqueDates, parsedDate)
			} else {
				log.Printf("[scheduler] failed to parse draw date '%s': %v", res.Date, err)
			}
		}
	}

	return uniqueDates, nil
}

// retry helper that executes a function up to 5 times with exponential backoff.
// It is context-aware and exits immediately if the context is cancelled.
func (c *LotteryClient) retry(ctx context.Context, fn func() ([]byte, error)) ([]byte, error) {
	maxAttempts := 5
	backoff := 500 * time.Millisecond

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		data, err := fn()
		if err == nil {
			return data, nil
		}

		lastErr = err
		log.Printf("[client] GLO API attempt %d failed: %v", attempt, err)

		if attempt < maxAttempts {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
				backoff *= 2
			}
		}
	}

	return nil, lastErr
}

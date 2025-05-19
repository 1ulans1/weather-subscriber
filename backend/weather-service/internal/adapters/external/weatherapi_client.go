package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"weather-service/internal/core/domain"
	"weather-service/internal/core/ports"
)

type weatherAPIClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewWeatherAPIClient(apiKey string) ports.WeatherAPIClient {
	baseURL := "http://api.weatherapi.com/v1"
	return &weatherAPIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

type weatherAPIResponse struct {
	Location struct {
		Name      string  `json:"name"`
		Region    string  `json:"region"`
		Country   string  `json:"country"`
		Lat       float64 `json:"lat"`
		Lon       float64 `json:"lon"`
		Localtime string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Humidity  int     `json:"humidity"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

func (w *weatherAPIClient) Fetch(ctx context.Context, location string) (*domain.Weather, error) {
	endpoint := fmt.Sprintf("%s/current.json", w.baseURL)
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("key", w.apiKey)
	q.Set("q", location)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weatherapi error: %s", resp.Status)
	}
	var data weatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &domain.Weather{
		Location:    data.Location.Name,
		Temperature: data.Current.TempC,
		Humidity:    data.Current.Humidity,
		Condition:   data.Current.Condition.Text,
		UpdatedAt:   time.Now(),
	}, nil
}

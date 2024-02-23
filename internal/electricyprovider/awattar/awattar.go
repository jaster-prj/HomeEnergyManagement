package awattar

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tobiasjaster/HomeEnergyManagement/internal/electricyprovider"
)

var (
	apiUrlMap = map[string]string{
		"de": "https://api.awattar.de/v1",
		"at": "https://api.awattar.at/v1",
	}
)

type AwattarElectricyProviderBuilder struct{}

func (b *AwattarElectricyProviderBuilder) Build() electricyprovider.IElectricyProvider {
	return &AwattarElectricyProvider{}
}

type AwattarElectricyProvider struct {
	location string
}

func (a *AwattarElectricyProvider) getUrl() (string, error) {
	apiUrl, ok := apiUrlMap[a.location]
	if !ok {
		return "", fmt.Errorf("no valid location: %v", a.location)
	}
	return apiUrl, nil
}

func (a *AwattarElectricyProvider) GetMarketDataRequest(start *time.Time, end *time.Time) (*http.Request, error) {
	apiUrl, err := a.getUrl()
	apiUrl = apiUrl + "/marketdata"
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "application/json")
	if start == nil && end == nil {
		return request, nil
	}
	values := url.Values{}
	if start != nil {
		timestamp := strconv.FormatInt(start.UTC().UnixMilli(), 10)
		values.Add("start", timestamp)
	}
	if start != nil {
		timestamp := strconv.FormatInt(end.UTC().UnixMilli(), 10)
		values.Add("end", timestamp)
	}
	request.URL.RawQuery = values.Encode()
	return request, nil
}

func (a *AwattarElectricyProvider) GetMarketData(start *time.Time, end *time.Time) ([]AwattarData, error) {
	request, err := a.GetMarketDataRequest(start, end)
	if err != nil {
		return []AwattarData{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return []AwattarData{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []AwattarData{}, err
	}

	var awattarResponse AwattarResponse
	if err := json.Unmarshal(body, &awattarResponse); err != nil {
		return []AwattarData{}, err
	}
	return awattarResponse.Data, nil
}

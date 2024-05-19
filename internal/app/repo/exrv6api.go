package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"genesis_test_task/internal/app/model"
)

type ExchangeRateRepo struct {
	apiUrl string
	log    *log.Logger
}

const apiUrlTemplate = "https://v6.exchangerate-api.com/v6/%s/pair/USD/UAH"

func NewExchangeRateRepo(apiKey string, logger *log.Logger) *ExchangeRateRepo {
	apiUrl := fmt.Sprintf(apiUrlTemplate, apiKey)
	return &ExchangeRateRepo{apiUrl: apiUrl, log: logger}
}

func (exrr ExchangeRateRepo) GetExchangeRate(ctx context.Context) (model.ExchangeRate, error) {
	response, err := http.Get(exrr.apiUrl)
	if err != nil {
		exrr.log.Panicf("Failed to get exchange rate: %v", err)
		return model.ExchangeRate{}, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	currencyRate := model.ExchangeRate{}
	err = json.Unmarshal(body, &currencyRate)
	if err != nil {
		log.Panicf("Failed to unmarshal JSON: %v", err)
		return model.ExchangeRate{}, err
	}
	return currencyRate, nil
}

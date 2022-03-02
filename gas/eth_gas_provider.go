package gas

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type EthQuote struct {
	Error  string `json:"error"`
	Result struct {
		Instant     *GasStruct `json:"instant"`
		Fast        *GasStruct `json:"fast"`
		Eco         *GasStruct `json:"eco"`
		BaseFee     float64    `json:"baseFee"`
		NativePrice float64    `json:"ethPrice"`
	} `json:"result"`
}

type EthGasProvider struct {
	url   string
	rates *EthQuote
}

func NewEthGasProvider() GasProvider {
	return &EthGasProvider{
		url:   "https://api.gasprice.io/v1/estimates",
		rates: nil,
	}
}

func (gp *EthGasProvider) Start() error {
	for {
		resp, err := http.Get(gp.url)
		if err != nil {
			log.Println("Error retreiving gas prices", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &gp.rates)
		if err != nil {
			log.Println("Error unmarshalling gas", err)
			return err
		}
		time.Sleep(5 * time.Second)
	}
}

func (gp *EthGasProvider) Instant() *GasStruct {
	return gp.rates.Result.Instant
}

func (gp *EthGasProvider) Fast() *GasStruct {
	return gp.rates.Result.Fast
}

func (gp *EthGasProvider) Eco() *GasStruct {
	return gp.rates.Result.Eco
}

func (gp *EthGasProvider) BaseFee() float64 {
	return gp.rates.Result.BaseFee
}

func (gp *EthGasProvider) NativePrice() float64 {
	return gp.rates.Result.NativePrice
}

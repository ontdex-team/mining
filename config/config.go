package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	DEFAULT_PRCIE_API     = "https://api.binance.com/api/v3/ticker/price"
	DEFAULT_WALLET        = "wallet.dat"
	DEFAULT_NET_ADDR      = "http://dappnode3.ont.io:20336"
	DEFAULT_CONTRACT_ADDR = "d3b733f12df9a6efb13ca547be5ee4e4dbe6d41e"
	DEAFULT_AMOUNT        = 10
	MIN_INTERVAL          = 10
)

type Config struct {
	PriceApi        string `json:"price_api"`
	Wallet          string `json:"wallet"`
	OntologyNetAddr string `json:"ontology_net_addr"`
	ContractAddr    string `json:"contract_addr"`
	Amount          int    `json:"amount"`
	Interval        uint   `json:"interval"`
}

func ParseConfig() (*Config, error) {
	fileContent, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(fileContent, config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed, err: %s", err)
	}
	config.checkDefaultCfg()
	return config, nil
}
func (this *Config) checkDefaultCfg() {
	if this.PriceApi == "" {
		this.PriceApi = DEFAULT_PRCIE_API
	}
	if this.Wallet == "" {
		this.Wallet = DEFAULT_WALLET
	}
	if this.OntologyNetAddr == "" {
		this.OntologyNetAddr = DEFAULT_NET_ADDR
	}
	if this.ContractAddr == "" {
		this.ContractAddr = DEFAULT_CONTRACT_ADDR
	}
	if this.Amount <= 0 {
		this.Amount = DEAFULT_AMOUNT
	}
	if this.Interval < MIN_INTERVAL {
		this.Interval = MIN_INTERVAL
	}
}

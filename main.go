package main

import (
	"encoding/json"
	"fmt"
	"github.com/ontdex-team/mining/config"
	"github.com/ontdex-team/mining/http"
	"strconv"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		fmt.Printf("parse config failed, err: %s\n", err)
		return
	}
	rest := http.NewRestClient()
	ongParam := map[string]string{"symbol": "ONGUSDT"}
	ontParam := map[string]string{"symbol": "ONTUSDT"}
	ongRespData, err := rest.SendGetRequest(cfg.PriceApi, ongParam)
	if err != nil {
		fmt.Printf("req ong price failed, err: %s\n", err)
		return
	}
	ontRespData, err := rest.SendGetRequest(cfg.PriceApi, ontParam)
	if err != nil {
		fmt.Printf("req ont price failed, err: %s\n", err)
		return
	}
	ong := &http.PriceResp{}
	if err := json.Unmarshal(ongRespData, ong); err != nil {
		fmt.Printf("unmarshal ong resp failed, err: %s\n", err)
		return
	}
	ont := &http.PriceResp{}
	if err := json.Unmarshal(ontRespData, ont); err != nil {
		fmt.Printf("unmarshal ont resp failed, err: %s\n", err)
		return
	}
	fmt.Printf("ong price is %s, ont price is %s\n", ong.Price, ont.Price)
	ongPrice, err := strconv.ParseFloat(ong.Price, 64)
	if err != nil {
		fmt.Printf("parse ong price failed, err: %s\n", err)
	}
	ontPrice, err := strconv.ParseFloat(ont.Price, 64)
	if err != nil {
		fmt.Printf("parse ont price failed, err: %s\n", err)
	}
	fmt.Printf("final price is %f ONT/ONG\n", ontPrice/ongPrice)
}

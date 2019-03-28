package main

import (
	"encoding/json"
	"fmt"
	"github.com/ontio/ontology-go-sdk/client"
	"github.com/ontio/ontology/common/password"
	"strconv"
	"time"

	"github.com/ontdex-team/mining/config"
	"github.com/ontdex-team/mining/http"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/log"
)

const (
	BUY       = 0
	SELL      = 1
	ADD_ORDER = "addOrder"
	PAIR_ID   = 1

	AMOUNT_MULTIPLE = 1000000000
	PRICE_MULTIPLE  = 1000000000
)

func main() {
	log.InitLog(log.InfoLog, log.PATH, log.Stdout)
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Errorf("parse config failed, err: %s", err)
		return
	}
	wallet, err := sdk.OpenWallet(cfg.Wallet)
	if err != nil {
		log.Errorf("open wallet failed, err: %s", err)
		return
	}
	pwd, err := password.GetPassword()
	if err != nil {
		log.Errorf("get password error: %sr", err)
		return
	}
	acc, err := wallet.GetDefaultAccount(pwd)
	if err != nil {
		log.Errorf("get account error: %sr", err)
		return
	}
	ontSdk := sdk.NewOntologySdk()
	ontSdk.SetDefaultClient(client.NewRpcClient().SetAddress(cfg.OntologyNetAddr))
	for {
		price, err := getPrice(cfg)
		if err != nil {
			log.Error(err)
		} else {
			log.Infof("ONT/ONG price is %f", price)
			addOrder(ontSdk, acc, cfg.ContractAddr, BUY, cfg.Amount, price)
			addOrder(ontSdk, acc, cfg.ContractAddr, SELL, cfg.Amount, price)
		}
		<-time.After(time.Duration(cfg.Interval) * time.Second)
	}
}

func getPrice(cfg *config.Config) (float64, error) {
	rest := http.NewRestClient()
	ongParam := map[string]string{"symbol": "ONGUSDT"}
	ontParam := map[string]string{"symbol": "ONTUSDT"}
	ongRespData, err := rest.SendGetRequest(cfg.PriceApi, ongParam)
	if err != nil {
		return 0, fmt.Errorf("req ong price failed, err: %s", err)
	}
	ontRespData, err := rest.SendGetRequest(cfg.PriceApi, ontParam)
	if err != nil {
		return 0, fmt.Errorf("req ont price failed, err: %s", err)
	}
	ong := &http.PriceResp{}
	if err := json.Unmarshal(ongRespData, ong); err != nil {
		return 0, fmt.Errorf("unmarshal ong resp failed, err: %s", err)
	}
	ont := &http.PriceResp{}
	if err := json.Unmarshal(ontRespData, ont); err != nil {
		return 0, fmt.Errorf("unmarshal ont resp failed, err: %s", err)
	}
	log.Infof("ong price is %s, ont price is %s", ong.Price, ont.Price)
	ongPrice, err := strconv.ParseFloat(ong.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("parse ong price failed, err: %s", err)
	}
	ontPrice, err := strconv.ParseFloat(ont.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("parse ont price failed, err: %s", err)
	}
	return ontPrice / ongPrice, nil
}

func addOrder(ontSdk *sdk.OntologySdk, acc *sdk.Account, contract string, orderType int, amount int, price float64) {
	//owner,pairid:1,amount,ordertype,price
	priceParam := int(price * PRICE_MULTIPLE)
	amountParam := amount * AMOUNT_MULTIPLE
	txHash, err := invokeSmartContract(ontSdk, acc, contract, ADD_ORDER,
		[]interface{}{acc.Address[:], PAIR_ID, amountParam, orderType, priceParam, "", ""})
	if err != nil {
		log.Errorf("add %d order failed, owner %s, pair %d, amount %d, price %d, err: %s", orderType,
			acc.Address.ToBase58(), PAIR_ID, amountParam, priceParam, err)
		return
	}
	log.Infof("add %d order success, owner %s, pair %d, amount %d, price %f, txHash:%s", orderType,
		acc.Address.ToBase58(), PAIR_ID, amount, price, txHash.ToHexString())
}

func invokeSmartContract(ontSdk *sdk.OntologySdk, signer *sdk.Account, contract, operation string,
	params []interface{}) (common.Uint256, error) {
	contractAddress, err := common.AddressFromHexString(contract)
	if err != nil {
		return common.UINT256_EMPTY, fmt.Errorf("InvokeSmartContract common.AddressFromHexString error:%s", err)
	}

	args := []interface{}{operation, params}
	txHash, err := ontSdk.NeoVM.InvokeNeoVMContract(500, 60000, signer, contractAddress, args)
	if err != nil {
		return common.UINT256_EMPTY, fmt.Errorf("InvokeSmartContract InvokeNeoVMContract error:%s", err)
	}
	return txHash, nil
}

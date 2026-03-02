package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// TODO: refactoring
// usage:
// -i <wallet_id>, -t <d> or <w>, -a <amount> -- deposit or withdraw amount
// or -i <wallet_id> -- get balance

func main() {
	walletID := flag.String("i", "", "wallet id")
	operationType := flag.String("t", "", "operation type: deposit (short d) or withdraw (short w)")
	amount := flag.Int("a", 0, "amount")

	flag.Parse()

	var url string
	var method string
	var reqBody string

	switch {
	case *walletID != "" && *operationType == "d" && *amount > 0:
		url = "http://localhost:8080/api/v1/wallets"
		method = http.MethodPost
		reqBody = fmt.Sprintf(`{"walletId":"%s","operationType":"DEPOSIT","amount":%d}`, *walletID, *amount)

	case *walletID != "" && *operationType == "w" && *amount > 0:
		url = "http://localhost:8080/api/v1/wallets"
		method = http.MethodPost
		reqBody = fmt.Sprintf(`{"walletId":"%s","operationType":"WITHDRAW","amount":%d}`, *walletID, *amount)

	case *walletID != "" && *operationType == "" && *amount == 0:
		url = "http://localhost:8080/api/v1/wallets/" + *walletID
		method = http.MethodGet

	default:
		log.Fatal(errors.New("flag error"))
	}

	client := http.Client{}

	req, err := http.NewRequest(method, url, strings.NewReader(reqBody))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp.Status, string(respBody))
}

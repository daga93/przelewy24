package przelewy24

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
)

type VerificationData struct {
	SessionId string
	OrderId   int
	Amount    int
	Currency  Currency
}

type signatureTransactionRegistration struct {
	SessionId  string   `json:"sessionId"`
	MerchantId int      `json:"merchantId"`
	Amount     int      `json:"amount"`
	Currency   Currency `json:"currency"`
	Crc        string   `json:"crc"`
}

type signatureTransactionVerification struct {
	SessionId string   `json:"sessionId"`
	OrderId   int      `json:"orderId"`
	Amount    int      `json:"amount"`
	Currency  Currency `json:"currency"`
	Crc       string   `json:"crc"`
}

func signRegistration(sessionId string, merchantId int, amount int, currency Currency, crc string) [48]byte {
	signature := signatureTransactionRegistration{
		SessionId:  sessionId,
		MerchantId: merchantId,
		Amount:     amount,
		Currency:   currency,
		Crc:        crc,
	}

	marshalled, err := json.Marshal(signature)
	if err != nil {
		fmt.Println(err)
		return [48]byte{}
	}
	return sha512.Sum384(marshalled)
}

func signVerification(sessionId string, orderId int, amount int, currency Currency, crc string) [48]byte {
	signature := signatureTransactionVerification{
		SessionId: sessionId,
		OrderId:   orderId,
		Amount:    amount,
		Currency:  currency,
		Crc:       crc,
	}

	marshalled, err := json.Marshal(signature)
	if err != nil {
		fmt.Println(err)
		return [48]byte{}
	}
	return sha512.Sum384(marshalled)
}

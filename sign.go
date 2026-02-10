package przelewy24

import (
	"crypto/sha512"
	"encoding/json"
)

type VerificationData struct {
	SessionId string
	OrderId   int
	Amount    int
	Currency  Currency
}

type RefundVerificationData struct {
	SessionId   string
	OrderId     int
	Amount      int
	Currency    Currency
	RefundsUUID string
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

type signatureRefundVerification struct {
	OrderId     int      `json:"orderId"`
	SessionId   string   `json:"sessionId"`
	RefundsUuid string   `json:"refundsUuid"`
	MerchantId  int      `json:"merchantId"`
	Amount      int      `json:"amount"`
	Currency    Currency `json:"currency"`
	Status      int      `json:"status"`
	Crc         string   `json:"crc"`
}

func getRegisterSign(sessionId string, merchantId int, amount int, currency Currency, crc string) ([48]byte, error) {
	signature := signatureTransactionRegistration{
		SessionId:  sessionId,
		MerchantId: merchantId,
		Amount:     amount,
		Currency:   currency,
		Crc:        crc,
	}

	marshalled, err := json.Marshal(signature)
	if err != nil {
		return [48]byte{}, err
	}
	return sha512.Sum384(marshalled), nil
}

func getVerificationSign(sessionId string, orderId int, amount int, currency Currency, crc string) ([48]byte, error) {
	signature := signatureTransactionVerification{
		SessionId: sessionId,
		OrderId:   orderId,
		Amount:    amount,
		Currency:  currency,
		Crc:       crc,
	}

	marshalled, err := json.Marshal(signature)
	if err != nil {
		return [48]byte{}, err
	}
	return sha512.Sum384(marshalled), nil
}

func getRefundSign(orderId int, sessionId string, refundsUuid string, merchantId int, amount int, currency Currency, status int, crc string) ([48]byte, error) {
	signature := signatureRefundVerification{
		OrderId:     orderId,
		SessionId:   sessionId,
		RefundsUuid: refundsUuid,
		MerchantId:  merchantId,
		Amount:      amount,
		Currency:    currency,
		Status:      status,
		Crc:         crc,
	}

	marshalled, err := json.Marshal(signature)
	if err != nil {
		return [48]byte{}, err
	}
	return sha512.Sum384(marshalled), nil
}

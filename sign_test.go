package przelewy24

import (
	"encoding/hex"
	"testing"
)

func TestSignRegistration(t *testing.T) {
	// Assert
	sessionID := "123"
	merchantID := 123
	amount := 100
	currency := PLN
	crc := "ABC"

	// Calculated on P24 tool
	expectedSignHex := "148fbf17ccdd4d070c24e84e040b5ac2fadd4700a3e7fcd31a37492969506a1e5bf466ac23292495b74b2612e4a8e037"
	expectedBytes, err := hex.DecodeString(expectedSignHex)
	if err != nil {
		t.Fatalf("Failed to decode hex string: %v", err)
	}
	var ex [48]byte
	copy(ex[:], expectedBytes)

	// Act
	got, _ := getRegisterSign(sessionID, merchantID, amount, currency, crc)

	// Assert
	if ex != got {
		t.Errorf("Wrong signature for SignRegistration calculated. Want: %s, Have: %s", ex, got)
	}
}

func TestSignVerify(t *testing.T) {
	// Assert
	sessionID := "123"
	orderID := 234
	amount := 100
	currency := PLN
	crc := "ABC"

	// Calculated on P24 tool
	expectedSignHex := "991389947144a6864ff0ed46ffda24111928c4549d961df88797bdca4d583da9cc7328c0848b5c024c56e22b58ebf4fb"
	expectedBytes, err := hex.DecodeString(expectedSignHex)
	if err != nil {
		t.Fatalf("Failed to decode hex string: %v", err)
	}
	var ex [48]byte
	copy(ex[:], expectedBytes)

	// Act
	got, _ := getVerificationSign(sessionID, orderID, amount, currency, crc)

	// Assert
	if ex != got {
		t.Errorf("Wrong signature for SignRegistration calculated. Want: %s, Have: %s", ex, got)
	}
}

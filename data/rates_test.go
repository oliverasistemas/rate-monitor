package data

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"strings"
	"testing"
)

func ErrorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}

func TestExchangeRates_GetRate(t *testing.T) {
	tr, err := NewExchangeRates(hclog.Default())

	if err != nil {
		t.Fatal(err)
	}

	_, err = tr.GetRate("not_exists", "usd")

	//Test that it returns error if symbol does not exist
	if !ErrorContains(err, "Rate not found") {
		t.Errorf("unexpected error: %v", err)
	}

	rate, err := tr.GetRate("eth", "usd")

	if err != nil {
		t.Errorf("Error getting exchange rates.")
		return
	}

	fmt.Printf("ETH/USD rate is %.4f", rate)
}

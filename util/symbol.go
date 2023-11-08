package util

const (
	ETH = "ETH"
	BTC = "BTC"
	ADA = "ADA"
)

func IsSupportedSymbol(symbol string) bool {
	switch symbol {
	case ETH, BTC, ADA:
		return true
	}
	return false
}

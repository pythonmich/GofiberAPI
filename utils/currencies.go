package utils

type CurrencyCode string

const (
	USD CurrencyCode = "USD"
	EUR CurrencyCode = "EUR"
	KES CurrencyCode = "KES"
	UGX CurrencyCode = "UGX"
	TZS CurrencyCode = "TZS"
)

func IsSupportedCurrency(currency CurrencyCode) bool {
	switch currency {
	case KES, USD, TZS, UGX, EUR:
		return true
	}
	return false
}

package tools

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

// ##########################################
// ## Helper functions specific to the EVM ##
// ##########################################

// returns ETH, MATIC, or first 6 chars
func TokenSymbol(client *ethclient.Client, chainId *big.Int, tokenAddr string) (string, error) {
	// TODO - add more
	if tokenAddr == "" {
		if chainId.Cmp(big.NewInt(1)) == 0 || chainId.Cmp(big.NewInt(4)) == 0 {
			return "ETH", nil
		} else if chainId.Cmp(big.NewInt(137)) == 0 || chainId.Cmp(big.NewInt(80001)) == 0 {
			return "MATIC", nil
		} else {
			return "", fmt.Errorf("native token for chainId %s %s", chainId, " not recognized")
		}
	}
	return tokenAddr[:6], nil
}

func ERC20TransferGas() *big.Int {
	return big.NewInt(52500)
}

func NativeTransferGas() *big.Int {
	return big.NewInt(21000)
}

// EthToWei automatically inputs 18 decimals
func EthToWei(iamount decimal.Decimal) *big.Int {
	return ToWei(iamount, 18)
}

// Gwei to Wei automatically inputs 9 decimals
func EthToGwei(iamount decimal.Decimal) *big.Int {
	return ToWei(iamount, 9)
}

// EthToWei automatically inputs 18 decimals
func GweiToWei(iamount decimal.Decimal) *big.Int {
	return ToWei(iamount, 9)
}

// WeiToEth automatically inputs 18 decimals
func WeiToEth(iamount *big.Int) decimal.Decimal {
	return ToDecimal(iamount, 18)
}

// GWeiToEth automatically inputs 18 decimals
func GweiToEth(iamount decimal.Decimal) decimal.Decimal {
	amt := big.NewInt(iamount.IntPart())
	return ToDecimal(amt, 9)
}

// GWeiToEth automatically inputs 18 decimals
func WeiToGwei(iamount *big.Int) decimal.Decimal {
	return ToDecimal(iamount, 9)
}

// ToDecimal wei to decimals
func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

// ToWei decimals to wei
func ToWei(iamount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

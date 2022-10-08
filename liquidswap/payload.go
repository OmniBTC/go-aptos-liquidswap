package liquidswap

import (
	"errors"
	"math/big"

	"github.com/shopspring/decimal"
)

type Payload struct {
	Function string
	TypeArgs []string
	Args     []string
}

type SwapParams struct {
	Script           string // eg. 0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9::scripts
	FromCoin         string
	ToCoin           string
	FromAmount       *big.Int
	ToAmount         *big.Int
	InteractiveToken string // from|to
	Slippage         decimal.Decimal
	Pool             Pool
}

func CreateSwapPayload(params *SwapParams) (*Payload, error) {
	if nil == params {
		return nil, errors.New("invalid params: nil")
	}

	if params.Slippage.LessThan(decimal.Zero) || params.Slippage.GreaterThan(decimal.New(1, 0)) {
		return nil, errors.New("Invalid slippage value:" + params.Slippage.String())
	}

	f := "swap"
	if params.InteractiveToken != "from" {
		f = "swap_into"
	}
	functionName := params.Script + "::" + f
	typeArgs := []string{
		params.FromCoin,
		params.ToCoin,
		params.Pool.CurveStructType,
	}
	if params.InteractiveToken != "from" {
		params.FromAmount = withSlippage(params.FromAmount, params.Slippage, 1)
	}
	if params.InteractiveToken == "from" {
		params.ToAmount = withSlippage(params.ToAmount, params.Slippage, -1)
	}
	args := []string{
		params.FromAmount.String(),
		params.ToAmount.String(),
	}
	return &Payload{
		Function: functionName,
		TypeArgs: typeArgs,
		Args:     args,
	}, nil
}

package liquidswap

import "math/big"

type CoinAmount struct {
	Coin
	Amount *big.Int
}

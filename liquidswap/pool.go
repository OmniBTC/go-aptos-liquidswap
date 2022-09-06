package liquidswap

import "math/big"

type PoolResource struct {
	CoinXReserve *big.Int
	CoinYReserve *big.Int
}

type Pool struct {
	LpToken       string
	ModuleAddress string
	Address       string
}

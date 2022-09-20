package liquidswap

import "math/big"

const (
	// Stable curve (like Solidly).
	StableCurve = 1
	// Uncorellated curve (like Uniswap).
	Uncorellated = 2
)

type PoolResource struct {
	CoinXReserve *big.Int
	CoinYReserve *big.Int
	CurveType    int // 0: same as Uncorellated
}

type Pool struct {
	LpToken       string
	ModuleAddress string
	Address       string
}

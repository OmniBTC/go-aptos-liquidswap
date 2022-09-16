package liquidswap

import (
	"math/big"

	"github.com/shopspring/decimal"
)

const (
	EQUAL        = 0
	LESS_THAN    = 1
	GREATER_THAN = 2
)

func AmountMinOut(val *big.Int, slippage decimal.Decimal) *big.Int {
	return withSlippage(val, slippage, -1)
}

func AmountMaxIn(val *big.Int, slippage decimal.Decimal) *big.Int {
	return withSlippage(val, slippage, 1)
}

func withSlippage(val *big.Int, slippage decimal.Decimal, mod int) *big.Int {
	if mod > 0 {
		return decimal.NewFromBigInt(val, 0).Add(decimal.NewFromBigInt(val, 0).Mul(slippage)).BigInt()
	} else {
		return decimal.NewFromBigInt(val, 0).Sub(decimal.NewFromBigInt(val, 0).Mul(slippage)).BigInt()
	}
}

func GetAmountIn(fromCoin, toCoin Coin, amountOut *big.Int, pool PoolResource) *big.Int {
	return calculateRates(fromCoin, toCoin, amountOut, "to", pool)
}

func GetAmountOut(fromCoin, toCoin Coin, amountIn *big.Int, pool PoolResource) *big.Int {
	return calculateRates(fromCoin, toCoin, amountIn, "from", pool)
}

func calculateRates(fromCoin, toCoin Coin, amount *big.Int, interactiveToken string, pool PoolResource) *big.Int {
	isSorted := IsSortedSymbols(fromCoin.Symbol, toCoin.Symbol)
	var (
		reserveX *big.Int
		reserveY *big.Int
	)
	reserveX, reserveY = pool.CoinXReserve, pool.CoinYReserve
	if !isSorted {
		reserveX, reserveY = reserveY, reserveX
	}
	if interactiveToken != "from" {
		reserveX, reserveY = reserveY, reserveX
	}

	if interactiveToken == "from" {
		return getCoinOutWithFees(amount, reserveX, reserveY)
	} else {
		return getCoinInWithFees(amount, reserveX, reserveY)
	}
}

func getCoinOutWithFees(coinIn, reserveIn, reserveOut *big.Int) *big.Int {
	feePct, feeScale := getFee()
	feeMultiplier := big.NewInt(0).Sub(feeScale, feePct)
	coinInAfterFees := big.NewInt(0).Mul(coinIn, feeMultiplier)
	newReservesInSize := big.NewInt(0).Add(big.NewInt(0).Mul(reserveIn, feeScale), coinInAfterFees)
	return big.NewInt(0).Div(big.NewInt(0).Mul(coinInAfterFees, reserveOut), newReservesInSize)
}

func getCoinInWithFees(coinOut, reserveOut, reserveIn *big.Int) *big.Int {
	feePct, feeScale := getFee()
	feeMultiplier := big.NewInt(0).Sub(feeScale, feePct)
	newReservesOut := big.NewInt(0).Mul(feeMultiplier, big.NewInt(0).Sub(reserveOut, coinOut))
	return big.NewInt(0).Div(big.NewInt(0).Mul(reserveIn, big.NewInt(0).Mul(coinOut, feeScale)), newReservesOut)
}

func getFee() (*big.Int, *big.Int) {
	return big.NewInt(3), big.NewInt(1000)
}

func IsSortedSymbols(symbolX, symbolY string) bool {
	return compare(symbolX, symbolY) == LESS_THAN
}

func compare(symbolX, symbolY string) int {
	ix := len(symbolX)
	iy := len(symbolY)
	lenCmp := cmp(ix, iy)
	// &bcs::to_bytes(utf8(b"hello")) in Aptos the first bytes contains length of string
	if lenCmp != EQUAL {
		return lenCmp
	}
	i := 0
	for i < ix && i < iy {
		elemCmp := cmp(int(symbolX[i]), int(symbolY[i]))
		if elemCmp != EQUAL {
			return elemCmp
		}
		i++
	}
	return EQUAL
}

func cmp(a, b int) int {
	if a == b {
		return EQUAL
	} else if a < b {
		return LESS_THAN
	} else {
		return GREATER_THAN
	}
}

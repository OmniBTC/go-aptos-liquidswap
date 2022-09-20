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

var one_e8 = big.NewInt(10000_0000) // 1e8

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
	coinX, coinY := fromCoin, toCoin
	if !isSorted {
		reserveX, reserveY = reserveY, reserveX
		coinX, coinY = coinY, coinX
	}
	if interactiveToken != "from" {
		reserveX, reserveY = reserveY, reserveX
		coinX, coinY = coinY, coinX
	}

	scaleX := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(coinX.Decimals)), nil)
	scaleY := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(coinY.Decimals)), nil)

	if interactiveToken == "from" {
		switch pool.CurveType {
		case StableCurve:
			return getStableCoinOutWithFees(amount, reserveX, reserveY, scaleX, scaleY)
		default:
			return getCoinOutWithFees(amount, reserveX, reserveY)
		}
	} else {
		switch pool.CurveType {
		case StableCurve:
			return getStableCoinInWithFees(amount, reserveX, reserveY, scaleX, scaleY)
		default:
			return getCoinInWithFees(amount, reserveX, reserveY)
		}
	}
}

func getStableCoinOutWithFees(coinIn, reserveIn, reserveOut, scaleIn, scaleOut *big.Int) *big.Int {
	feePct, feeScale := getFee()
	feeMultiplier := big.NewInt(0).Sub(feeScale, feePct)
	coinInAfterFees := big.NewInt(0).Mul(coinIn, feeMultiplier)
	if big.NewInt(0).Mod(coinInAfterFees, feeScale).Cmp(big.NewInt(0)) == 0 {
		coinInAfterFees = big.NewInt(0).Div(coinInAfterFees, feeScale)
	} else {
		coinInAfterFees = big.NewInt(0).Add(big.NewInt(1), big.NewInt(0).Div(coinInAfterFees, feeScale))
	}
	xy := lp_value(reserveIn, reserveOut, scaleIn, scaleOut)
	_reserveIn := toStableDecimal(reserveIn, scaleIn)
	_reserveOut := toStableDecimal(reserveOut, scaleOut)
	_coinIn := toStableDecimal(coinInAfterFees, scaleIn)
	totalReserve := big.NewInt(0).Add(_coinIn, _reserveIn)
	y := big.NewInt(0).Sub(_reserveOut, gety(totalReserve, xy, _reserveOut))
	return big.NewInt(0).Div(big.NewInt(0).Mul(y, scaleOut), one_e8)
}

func getStableCoinInWithFees(coinOut, reserveOut, reserveIn, scaleOut, scaleIn *big.Int) *big.Int {
	feePct, feeScale := getFee()
	feeMultiplier := big.NewInt(0).Sub(feeScale, feePct)
	xy := lp_value(reserveIn, reserveOut, scaleIn, scaleOut)
	_reserveIn := toStableDecimal(reserveIn, scaleIn)
	_reserveOut := toStableDecimal(reserveOut, scaleOut)
	_coinOut := toStableDecimal(coinOut, scaleOut)
	totalReserve := big.NewInt(0).Sub(_reserveOut, _coinOut)
	x := big.NewInt(0).Sub(gety(totalReserve, xy, _reserveIn), _reserveIn)
	coinIn := big.NewInt(0).Add(
		big.NewInt(1),
		big.NewInt(0).Div(big.NewInt(0).Mul(x, scaleIn), one_e8),
	)
	return big.NewInt(0).Add(
		big.NewInt(1),
		big.NewInt(0).Div(
			big.NewInt(0).Mul(coinIn, feeScale),
			feeMultiplier,
		),
	)
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

func lp_value(x, y, scaleX, scaleY *big.Int) *big.Int {
	_x := toStableDecimal(x, scaleX)
	_y := toStableDecimal(y, scaleY)
	_a := big.NewInt(0).Mul(_x, _y)
	_b := big.NewInt(0).Add(big.NewInt(0).Mul(_x, _x), big.NewInt(0).Mul(_y, _y))
	// _x*_y * (_x*_x + _y*_y)
	return big.NewInt(0).Mul(_a, _b)
}

func toStableDecimal(amount, scale *big.Int) *big.Int {
	return big.NewInt(0).Div(big.NewInt(0).Mul(amount, one_e8), scale)
}

func gety(x0, xy, y *big.Int) *big.Int {
	one := big.NewInt(1)
	for i := 0; i < 255; i++ {
		k := f(x0, y)
		_dy := big.NewInt(0)
		cmp := k.Cmp(xy)
		if cmp < 0 {
			// (xy-k)/d(x0,y) + 1  Round Up
			_dy := big.NewInt(0).Add(big.NewInt(0).Div(big.NewInt(0).Sub(xy, k), d(x0, y)), one)
			y = big.NewInt(0).Add(y, _dy)
		} else {
			_dy := big.NewInt(0).Div(big.NewInt(0).Sub(xy, k), d(x0, y))
			y = big.NewInt(0).Add(y, _dy)
		}
		cmp = _dy.Cmp(one)
		if cmp <= 0 {
			return y
		}
	}
	return y
}

func f(x0, y *big.Int) *big.Int {
	yy := big.NewInt(0).Mul(y, y)
	yyy := big.NewInt(0).Mul(y, yy)
	a := big.NewInt(0).Mul(x0, yyy)
	xx := big.NewInt(0).Mul(x0, x0)
	xxx := big.NewInt(0).Mul(xx, x0)
	b := big.NewInt(0).Mul(xxx, y)
	return big.NewInt(0).Add(a, b)
}

func d(x, y *big.Int) *big.Int {
	three := big.NewInt(3)
	x3 := big.NewInt(0).Mul(three, x)
	yy := big.NewInt(0).Mul(y, y)
	xyy3 := big.NewInt(0).Mul(x3, yy)
	xx := big.NewInt(0).Mul(x, x)
	xxx := big.NewInt(0).Mul(xx, x)
	return big.NewInt(0).Add(xyy3, xxx)
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

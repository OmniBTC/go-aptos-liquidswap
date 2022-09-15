package liquidswap

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestCalculateRates(t *testing.T) {
	type args struct {
		fromCoin         Coin
		toCoin           Coin
		amount           *big.Int
		interactiveToken string
		pool             PoolResource
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{
			name: "case out",
			args: args{
				fromCoin:         Coin{Symbol: "USDT"},
				toCoin:           Coin{Symbol: "BTC"},
				amount:           big.NewInt(1000000),
				interactiveToken: "from",
				pool: PoolResource{
					CoinXReserve: big.NewInt(10415880990),
					CoinYReserve: big.NewInt(3004784231600),
				},
			},
			want: big.NewInt(3456),
		},
		{
			name: "case in",
			args: args{
				fromCoin:         Coin{Symbol: "USDT"},
				toCoin:           Coin{Symbol: "BTC"},
				amount:           big.NewInt(200000),
				interactiveToken: "to",
				pool: PoolResource{
					CoinXReserve: big.NewInt(10419434957),
					CoinYReserve: big.NewInt(3005809484015),
				},
			},
			want: big.NewInt(57870929),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateRates(tt.args.fromCoin, tt.args.toCoin, tt.args.amount, tt.args.interactiveToken, tt.args.pool); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateRates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithSlippage(t *testing.T) {
	type args struct {
		val      *big.Int
		slippage decimal.Decimal
		mod      int
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{
			name: "case out(from)",
			args: args{
				val:      big.NewInt(266607),
				slippage: decimal.NewFromFloat(0.005),
				mod:      -1,
			},
			want: big.NewInt(265273),
		},
		{
			name: "case in(to)",
			args: args{
				val:      big.NewInt(750174),
				slippage: decimal.NewFromFloat(0.005),
				mod:      1,
			},
			want: big.NewInt(753924),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := withSlippage(tt.args.val, tt.args.slippage, tt.args.mod); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithSlippage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAmountMinOut(t *testing.T) {
	type args struct {
		val      *big.Int
		slippage decimal.Decimal
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{
			name: "case out(from)",
			args: args{
				val:      big.NewInt(266607),
				slippage: decimal.NewFromFloat(0.005),
			},
			want: big.NewInt(265273),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AmountMinOut(tt.args.val, tt.args.slippage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AmountMinOut() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAmountMaxIn(t *testing.T) {
	type args struct {
		val      *big.Int
		slippage decimal.Decimal
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{
			name: "case in(to)",
			args: args{
				val:      big.NewInt(750174),
				slippage: decimal.NewFromFloat(0.005),
			},
			want: big.NewInt(753924),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AmountMaxIn(tt.args.val, tt.args.slippage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AmountMaxIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compare(t *testing.T) {
	type args struct {
		symbolX string
		symbolY string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "APTOS BTC",
			args: args{"APTOS", "BTC"},
			want: LESS_THAN,
		},
		{
			name: "BTC APTOS",
			args: args{"BTC", "APTOS"},
			want: GREATER_THAN,
		},
		{
			name: "BTC USDT",
			args: args{"BTC", "USDT"},
			want: LESS_THAN,
		},
		{
			name: "USDT BTC",
			args: args{"USDT", "BTC"},
			want: GREATER_THAN,
		},
		{
			name: "USDC USDT",
			args: args{"USDC", "USDT"},
			want: LESS_THAN,
		},
		{
			name: "APTOS USDT",
			args: args{"APTOS", "USDT"},
			want: LESS_THAN,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compare(tt.args.symbolX, tt.args.symbolY); got != tt.want {
				t.Errorf("compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

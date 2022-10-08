package liquidswap

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestCreateTxPayload(t *testing.T) {
	type args struct {
		params *SwapParams
	}
	tests := []struct {
		name    string
		args    args
		want    *Payload
		wantErr bool
	}{
		{
			name: "case out",
			args: args{
				&SwapParams{
					Script:           "0x123::scripts",
					FromCoin:         "0x123::BTC",
					ToCoin:           "0x123::APT",
					FromAmount:       big.NewInt(1),
					ToAmount:         big.NewInt(266607),
					InteractiveToken: "from",
					Slippage:         decimal.NewFromFloat(0.005),
					Pool: Pool{
						CurveStructType: "0x4e9fce03284c0ce0b86c88dd5a46f050cad2f4f33c4cdd29d98f501868558c81::curves::Uncorrelated",
					},
				},
			},
			want: &Payload{
				Function: "0x123::scripts::swap",
				TypeArgs: []string{
					"0x123::BTC",
					"0x123::APT",
					"0x4e9fce03284c0ce0b86c88dd5a46f050cad2f4f33c4cdd29d98f501868558c81::curves::Uncorrelated",
				},
				Args: []string{
					"1",
					"265273",
				},
			},
			wantErr: false,
		},
		{
			name: "case in",
			args: args{
				&SwapParams{
					Script:           "0x123::scripts",
					FromCoin:         "0x123::BTC",
					ToCoin:           "0x123::APT",
					FromAmount:       big.NewInt(750174),
					ToAmount:         big.NewInt(1),
					InteractiveToken: "to",
					Slippage:         decimal.NewFromFloat(0.005),
					Pool: Pool{
						CurveStructType: "0x4e9fce03284c0ce0b86c88dd5a46f050cad2f4f33c4cdd29d98f501868558c81::curves::Uncorrelated",
					},
				},
			},
			want: &Payload{
				Function: "0x123::scripts::swap_into",
				TypeArgs: []string{
					"0x123::BTC",
					"0x123::APT",
					"0x4e9fce03284c0ce0b86c88dd5a46f050cad2f4f33c4cdd29d98f501868558c81::curves::Uncorrelated",
				},
				Args: []string{
					"753924",
					"1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateSwapPayload(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTxPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateTxPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

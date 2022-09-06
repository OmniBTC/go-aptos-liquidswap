package liquidswap

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestCreateTxPayload(t *testing.T) {
	type args struct {
		params *CreateTxPayloadParams
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
				&CreateTxPayloadParams{
					Script:           "0x123::scripts",
					FromCoin:         "0x123::BTC",
					ToCoin:           "0x123::APT",
					FromAmount:       big.NewInt(1),
					ToAmount:         big.NewInt(266607),
					InteractiveToken: "from",
					Slippage:         decimal.NewFromFloat(0.005),
					Pool: Pool{
						LpToken:       "0x123::lp<0x123::APT,0x123::BTC>",
						ModuleAddress: "0x1234",
						Address:       "0x12345",
					},
				},
			},
			want: &Payload{
				Function: "0x123::scripts::swap",
				TypeArgs: []string{
					"0x123::BTC",
					"0x123::APT",
					"0x123::lp<0x123::APT,0x123::BTC>",
				},
				Args: []string{
					"0x12345",
					"1",
					"265273",
				},
			},
			wantErr: false,
		},
		{
			name: "case in",
			args: args{
				&CreateTxPayloadParams{
					Script:           "0x123::scripts",
					FromCoin:         "0x123::BTC",
					ToCoin:           "0x123::APT",
					FromAmount:       big.NewInt(750174),
					ToAmount:         big.NewInt(1),
					InteractiveToken: "to",
					Slippage:         decimal.NewFromFloat(0.005),
					Pool: Pool{
						LpToken:       "0x123::lp<0x123::APT,0x123::BTC>",
						ModuleAddress: "0x1234",
						Address:       "0x12345",
					},
				},
			},
			want: &Payload{
				Function: "0x123::scripts::swap_into",
				TypeArgs: []string{
					"0x123::BTC",
					"0x123::APT",
					"0x123::lp<0x123::APT,0x123::BTC>",
				},
				Args: []string{
					"0x12345",
					"753924",
					"1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateTxPayload(tt.args.params)
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

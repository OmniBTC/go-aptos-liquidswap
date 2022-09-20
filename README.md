# go-aptos-liquidswap
A go sdk for liquidswap on aptos.

## Install

```sh
go get github.com/coming-chat/go-aptos-liquidswap
```

## Usage

Get amount out:
```go
amountOut := GetAmountOut(Coin{Symbol: "USDT"}, Coin{Symbol: "BTC"}, big.NewInt(1000000), PoolResource{
    CoinXReserve: big.NewInt(10415880990),
    CoinYReserve: big.NewInt(3004784231600),
    CurveType: Uncorellated,
})
```

Calc amount min out:
```go
amountMinOut := AmountMinOut(amountOut, decimal.NewFromFloat(0.005))
```


Get amount in:
```go
amountIn := GetAmountIn(Coin{Symbol: "USDT"}, Coin{Symbol: "BTC"}, big.NewInt(1000000), PoolResource{
    CoinXReserve: big.NewInt(10415880990),
    CoinYReserve: big.NewInt(3004784231600),
})
```

Get amount max in:
```go
amountMaxOut := AmountMaxIn(amountIn, decimal.NewFromFloat(0.005))
```

Use **StableCurve**

> liquidswap pool type can be StableCurve for stable coins swap.

```go
amountIn := GetAmountIn(Coin{Symbol: "USDT"}, Coin{Symbol: "USDC"}, big.NewInt(1000000), PoolResource{
    CoinXReserve: big.NewInt(81442051331),
    CoinYReserve: big.NewInt(136352475461),
    CurveType:    StableCurve,
})
```

Create payload info:
```go
params := &SwapParams{
        Script:           "0x123::scripts",
        FromCoin:         "0x123::BTC",
        ToCoin:           "0x123::APT",
        FromAmount:       big.NewInt(1),
        ToAmount:         big.NewInt(266607),
        InteractiveToken: "from",  // from|to， from - exactIn  to - exactOut
        Slippage:         decimal.NewFromFloat(0.005),
        Pool: Pool{
            LpToken:       "0x123::lp<0x123::APT,0x123::BTC>",
            ModuleAddress: "0x1234",
            Address:       "0x12345",
        },
    }
payload, err := CreateSwapPayload(params)
```

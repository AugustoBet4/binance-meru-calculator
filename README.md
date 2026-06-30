# Binance Meru Calculator

Interactive Go calculator for comparing the BOB cost of buying USDC directly from Meru versus buying USDC on Binance and transferring it to Meru.

## Requirements

- Go installed
- A terminal that supports interactive prompts

The project uses Go modules, so dependencies are installed from `go.mod` and `go.sum`.

## Run the Calculator

From the project directory:

```powershell
go run .
```

Or run the file directly:

```powershell
go run calculator.go
```

The calculator will prompt for:

- Target USDC amount
- Binance exchange rate in BOB per USDC
- Meru exchange rate in BOB per USDC

Default values are already filled in:

- Target USDC amount: `400`
- Binance rate: `9.6`
- Meru rate: `9.8`

Select `Yes!` when asked `Calculate?` to print the comparison.

## Install or Refresh Dependencies

If dependencies are missing, run:

```powershell
go mod download
```

You can also tidy the module file after dependency changes:

```powershell
go mod tidy
```

## Build an Executable

To build a Windows executable:

```powershell
go build -o binance-meru-calculator.exe .
```

Then run:

```powershell
.\binance-meru-calculator.exe
```

## Calculation Notes

The Go version currently uses these fixed fees inside `calculator.go`:

- Binance transfer fee: `1.0` USDC
- Meru flat fee: `1.0` USDC
- Meru percentage fee: `1.0%`

The app compares:

- Buying the target USDC amount directly from Meru
- Buying enough USDC on Binance to cover Binance and Meru fees, then transferring to Meru

It prints the total BOB cost, effective rate, best option, and estimated savings.

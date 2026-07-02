//go:build sell
// +build sell

package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
)

// SellResult holds the sell-side comparison results.
type SellResult struct {
	MeruNetBOB    float64
	BinanceNetBOB float64
	DifferenceBOB float64
	BestOption    string
}

// calculateUSDCSellComparison compares selling USDC in Meru versus sending it to Binance and selling there.
func calculateUSDCSellComparison(usdcAmount, binanceFeeUSDC, binanceFeePercent, binanceRate, meruRate float64) SellResult {
	// Option 1: Sell directly in Meru.
	// User clarified there is no direct Meru fee in this path.
	bobNetMeru := usdcAmount * meruRate

	// Option 2: Transfer to Binance, then sell there.
	// Stellar transfer fee: fixed 1 USDC plus 0.5% of the amount.
	usdcAfterTransferFee := usdcAmount - binanceFeeUSDC - (usdcAmount * binanceFeePercent / 100)
	if usdcAfterTransferFee < 0 {
		usdcAfterTransferFee = 0
	}
	bobNetBinance := usdcAfterTransferFee * binanceRate

	// Determine best option.
	differenceBOB := bobNetBinance - bobNetMeru
	bestOption := ""
	profitPercent := 0.0

	if bobNetMeru > bobNetBinance {
		profitPercent = ((bobNetMeru - bobNetBinance) / bobNetBinance) * 100
		bestOption = "meru"
	} else {
		if bobNetMeru > 0 {
			profitPercent = ((bobNetBinance - bobNetMeru) / bobNetMeru) * 100
		}
		bestOption = "binance"
	}

	// Display results.
	fmt.Println("\n" + strings.Repeat("━", 70))
	fmt.Println("                  💰 USDC SELL COMPARISON")
	fmt.Println(strings.Repeat("━", 70))

	fmt.Printf("\n  Target: %.2f USDC\n", usdcAmount)
	fmt.Printf("  Rates:  Binance %.2f BOB/USDC  |  Meru %.2f BOB/USDC\n", binanceRate, meruRate)

	fmt.Println("\n" + strings.Repeat("─", 70))
	fmt.Println("  📊 OPTION 1: Sell in Meru")
	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("     Start with:              %.4f USDC\n", usdcAmount)
	fmt.Printf("     After Meru fee:          %.4f USDC\n", usdcAmount)
	fmt.Printf("     BOB received:             %.2f BOB\n", bobNetMeru)

	fmt.Println("\n" + strings.Repeat("─", 70))
	fmt.Println("  📊 OPTION 2: Transfer to Binance + Sell")
	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("     Start with:              %.4f USDC\n", usdcAmount)
	fmt.Printf("     After transfer fee:      %.4f USDC\n", usdcAfterTransferFee)
	fmt.Printf("     BOB received:            %.2f BOB\n", bobNetBinance)
	fmt.Printf("     Binance flat fee:        %.4f USDC\n", binanceFeeUSDC)
	fmt.Printf("     Binance transfer fee:     %.2f%% + %.4f USDC\n", binanceFeePercent, binanceFeeUSDC)

	fmt.Println("\n" + strings.Repeat("━", 70))
	if bestOption == "meru" {
		fmt.Println("  ✓ BEST OPTION: Sell in Meru")
	} else {
		fmt.Println("  ✓ BEST OPTION: Send to Binance and sell there")
	}
	fmt.Printf("  💵 You make more: %.2f BOB (%.2f%%)\n", abs(differenceBOB), profitPercent)
	fmt.Println(strings.Repeat("━", 70))

	return SellResult{
		MeruNetBOB:    bobNetMeru,
		BinanceNetBOB: bobNetBinance,
		DifferenceBOB: differenceBOB,
		BestOption:    bestOption,
	}
}

// abs returns the absolute value of a float64.
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	var (
		usdcAmount     string
		binanceRate    string
		meruRate       string
		runCalculation bool
	)

	// Fixed fees.
	const (
		binanceFeeUSDC    = 1.0
		binanceFeePercent = 0.5
	)

	// Default values.
	usdcAmount = "400"
	binanceRate = "9.98"
	meruRate = "9.6"

	for {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Target USDC Amount").
					Value(&usdcAmount).
					Placeholder("400"),

				huh.NewInput().
					Title("Binance Exchange Rate (BOB per USDC)").
					Value(&binanceRate).
					Placeholder("9.6"),

				huh.NewInput().
					Title("Meru Exchange Rate (BOB per USDC)").
					Value(&meruRate).
					Placeholder("9.8"),

				huh.NewConfirm().
					Title("Calculate?").
					Value(&runCalculation).
					Affirmative("Yes!").
					Negative("Exit"),
			),
		)

		err := form.Run()
		if err != nil {
			fmt.Println("Exiting...")
			break
		}

		if !runCalculation {
			fmt.Println("Goodbye!")
			break
		}

		var (
			usdcAmountFloat  float64
			binanceRateFloat float64
			meruRateFloat    float64
		)

		_, err = fmt.Sscanf(usdcAmount, "%f", &usdcAmountFloat)
		if err != nil {
			fmt.Printf("Invalid USDC amount: %v\n", err)
			continue
		}

		_, err = fmt.Sscanf(binanceRate, "%f", &binanceRateFloat)
		if err != nil {
			fmt.Printf("Invalid Binance rate: %v\n", err)
			continue
		}

		_, err = fmt.Sscanf(meruRate, "%f", &meruRateFloat)
		if err != nil {
			fmt.Printf("Invalid Meru rate: %v\n", err)
			continue
		}

		fmt.Println()
		calculateUSDCSellComparison(
			usdcAmountFloat,
			binanceFeeUSDC,
			binanceFeePercent,
			binanceRateFloat,
			meruRateFloat,
		)

		fmt.Println("\n  Press Enter to calculate again...")
		fmt.Scanln()
		runCalculation = false
	}
}

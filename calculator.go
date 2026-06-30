package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
)

// Result holds the calculation results
type Result struct {
	MeruCostBOB    float64
	BinanceCostBOB float64
	DifferenceBOB  float64
	BestOption     string
}

// calculateUSDCPurchase calculates the best way to buy USDC - compare costs in BOB
func calculateUSDCPurchase(usdcAmount, binanceFeeUSDC, meruFeeUSDC, meruFeePercent, binanceRate, meruRate float64) Result {
	// Calculate Option 1: Buy directly from Meru
	bobCostMeru := usdcAmount * meruRate

	// Calculate Option 2: Buy from Binance and transfer to Meru
	usdcBeforeMeruPercentFee := usdcAmount / (1 - meruFeePercent/100)
	usdcBeforeMeruFees := usdcBeforeMeruPercentFee + meruFeeUSDC
	usdcToBuyBinance := usdcBeforeMeruFees + binanceFeeUSDC
	bobCostBinance := usdcToBuyBinance * binanceRate

	// Determine best option
	differenceBOB := bobCostMeru - bobCostBinance
	bestOption := ""
	savingsPercent := 0.0

	if bobCostMeru < bobCostBinance {
		savingsPercent = ((bobCostBinance - bobCostMeru) / bobCostBinance) * 100
		bestOption = "meru"
	} else {
		savingsPercent = ((bobCostMeru - bobCostBinance) / bobCostMeru) * 100
		bestOption = "binance"
	}

	// Display results
	fmt.Println("\n" + strings.Repeat("━", 70))
	fmt.Println("                     💰 USDC PURCHASE COMPARISON")
	fmt.Println(strings.Repeat("━", 70))

	fmt.Printf("\n  Target: %.2f USDC\n", usdcAmount)
	fmt.Printf("  Rates:  Binance %.2f BOB/USDC  |  Meru %.2f BOB/USDC\n", binanceRate, meruRate)

	fmt.Println("\n" + strings.Repeat("─", 70))
	fmt.Println("  📊 OPTION 1: Buy directly from Meru")
	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("     Cost: %.2f BOB\n", bobCostMeru)

	fmt.Println("\n" + strings.Repeat("─", 70))
	fmt.Println("  📊 OPTION 2: Buy from Binance + Transfer")
	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("     Buy on Binance:           %.4f USDC → %.2f BOB\n", usdcToBuyBinance, bobCostBinance)
	fmt.Printf("     After Binance fee:        %.4f USDC\n", usdcBeforeMeruFees)
	fmt.Printf("     After Meru flat fee:      %.4f USDC\n", usdcBeforeMeruPercentFee)
	fmt.Printf("     Final amount:             %.2f USDC\n", usdcAmount)
	fmt.Printf("     Effective rate:           %.2f BOB/USDC\n", bobCostBinance/usdcAmount)

	fmt.Println("\n" + strings.Repeat("━", 70))
	if bestOption == "meru" {
		fmt.Println("  ✓ BEST OPTION: Buy directly from Meru")
	} else {
		fmt.Println("  ✓ BEST OPTION: Buy from Binance and transfer")
	}
	fmt.Printf("  💵 You save: %.2f BOB (%.2f%%)\n", abs(differenceBOB), savingsPercent)
	fmt.Println(strings.Repeat("━", 70))

	return Result{
		MeruCostBOB:    bobCostMeru,
		BinanceCostBOB: bobCostBinance,
		DifferenceBOB:  differenceBOB,
		BestOption:     bestOption,
	}
}

// abs returns the absolute value of a float64
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

	// Fixed fees
	const (
		binanceFeeUSDC = 1.0
		meruFeeUSDC    = 1.0
		meruFeePercent = 1.0
	)

	// Set default values
	usdcAmount = "400"
	binanceRate = "9.6"
	meruRate = "9.8"

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

		// Parse inputs
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

		// Run calculation
		fmt.Println()
		calculateUSDCPurchase(
			usdcAmountFloat,
			binanceFeeUSDC,
			meruFeeUSDC,
			meruFeePercent,
			binanceRateFloat,
			meruRateFloat,
		)

		fmt.Println("\n  Press Enter to calculate again...")
		fmt.Scanln()
		runCalculation = false
	}
}

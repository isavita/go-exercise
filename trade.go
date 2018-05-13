package main

import (
	"fmt"
	"math"
)

const epsilonMinDifference = 0.01

var priceLookup = map[int]float64{
	1: 2.44,
	2: 0.6,
	3: 2.0,
	4: 1.1,
	5: 1.5,
}

type Split struct {
	fund       int
	percentage float64
}

type SplitSlice []Split

func (p *SplitSlice) removeAtIndex(index int) {
	splits := *p
	splits = append(splits[:index], splits[index+1:]...)
	*p = splits
}

type Investment struct {
	fund       int
	percentage float64
	units      float64
}

type InvestmentSlice []Investment

func (p *InvestmentSlice) removeAtIndex(index int) {
	investments := *p
	investments = append(investments[:index], investments[index+1:]...)
	*p = investments
}

type SwitchInvestment struct {
	from_fund int
	to_fund   int
	units     float64
	state     string
}

func initOldInvestments() InvestmentSlice {
	oldInvestments := InvestmentSlice{
		Investment{fund: 1, percentage: 30, units: 1.45},
		Investment{fund: 2, percentage: 20, units: 2.50},
		Investment{fund: 3, percentage: 50, units: 6.15},
	}

	return oldInvestments
}

func initNewSplits() SplitSlice {
	newSplits := SplitSlice{
		Split{fund: 2, percentage: 50},
		Split{fund: 4, percentage: 50},
	}

	return newSplits
}

func printInvestmentsInfo(investments InvestmentSlice) {
	println("Old investments:")
	var sum float64
	for _, inv := range investments {
		amount := inv.units * priceLookup[inv.fund]
		fmt.Println("\tfund:", inv.fund, "->", inv, "(amount:", amount, ")")
		sum += amount
	}
	fmt.Println("\tvaluation amount:", sum)
}

func printSplitsInfo(splits SplitSlice) {
	println("New splits:")
	for _, split := range splits {
		fmt.Println("fund:", split.fund, "percentage:", split.percentage)
	}
}

func printSwitchInvestmentsInfo(switchInvestments []SwitchInvestment) {
	println("New switch investments:")
	for _, switch_inv := range switchInvestments {
		fmt.Println(
			"from fund:", switch_inv.from_fund,
			"to fund:", switch_inv.to_fund,
			"units:", switch_inv.units,
			"state:", switch_inv.state,
		)
	}
}

func resolveIntersect(investments *InvestmentSlice, splits *SplitSlice) []SwitchInvestment {
	switchInvestments := make([]SwitchInvestment, 0, len(*splits))
	var index int
	for splitIndex, split := range *splits {
		index = findIndexByFund(*investments, split.fund)
		if index != -1 {
			switchInvestment := createSwitchInvestmentSameFund(investments, index, splits, splitIndex)
			switchInvestments = append(switchInvestments, switchInvestment)
		}
	}

	return switchInvestments
}

func createSwitchInvestmentSameFund(investments *InvestmentSlice, invIndex int, splits *SplitSlice, spIndex int) SwitchInvestment {
	investment := (*investments)[invIndex]
	split := &(*splits)[spIndex]
	var units float64

	if investment.percentage > split.percentage {
		ratio := split.percentage / investment.percentage
		units = ratio * investment.units
		investment.units -= units
		investment.percentage -= split.percentage
		splits.removeAtIndex(spIndex)
	} else if investment.percentage < split.percentage {
		units = investment.units
		split.percentage -= investment.percentage
		investments.removeAtIndex(invIndex)
	} else {
		units = investment.units
		splits.removeAtIndex(spIndex)
		investments.removeAtIndex(invIndex)
	}

	switchInvestment := SwitchInvestment{
		from_fund: investment.fund,
		to_fund:   investment.fund,
		units:     units,
		state:     "completed",
	}

	return switchInvestment
}

func resolveSimmetricDifference(investments *InvestmentSlice, splits *SplitSlice) []SwitchInvestment {
	totalAmount := calcTotalAmount(investments)
	totalPercentage := calcTotalPercentage(investments)
	switchInvestments := make([]SwitchInvestment, 0, len(*splits)+len(*investments))

	for _, currentSplit := range *splits {
		currentSplitAmount := totalAmount * (currentSplit.percentage / totalPercentage)

		for invIndex, inv := range *investments {
			investmentAmount := inv.units * priceLookup[inv.fund]

			if math.Abs(investmentAmount-currentSplitAmount) <= epsilonMinDifference {
				totalAmount -= investmentAmount
				totalPercentage -= currentSplit.percentage
				investments.removeAtIndex(invIndex)
				switchInvestments = append(switchInvestments, SwitchInvestment{from_fund: inv.fund, to_fund: currentSplit.fund, units: inv.units, state: "created"})
				break

			} else if investmentAmount > currentSplitAmount {
				totalAmount -= currentSplitAmount
				totalPercentage -= currentSplit.percentage
				units := (currentSplitAmount / investmentAmount) * inv.units
				inv.units -= units
				switchInvestments = append(switchInvestments, SwitchInvestment{from_fund: inv.fund, to_fund: currentSplit.fund, units: units, state: "created"})
				break
			} else {
				totalAmount -= investmentAmount
				totalPercentage -= (investmentAmount / currentSplitAmount) * float64(currentSplit.percentage)
				investments.removeAtIndex(invIndex)
				switchInvestments = append(switchInvestments, SwitchInvestment{from_fund: inv.fund, to_fund: currentSplit.fund, units: inv.units, state: "created"})
			}
		}
	}

	return switchInvestments
}

func calcTotalAmount(investments *InvestmentSlice) float64 {
	var amount float64
	for _, inv := range *investments {
		amount += inv.units * priceLookup[inv.fund]
	}

	return amount
}

func calcTotalPercentage(investments *InvestmentSlice) float64 {
	var percentage float64
	for _, inv := range *investments {
		percentage += inv.percentage
	}

	return percentage
}

func findIndexByFund(investments InvestmentSlice, fund int) int {
	for i, inv := range investments {
		if inv.fund == fund {
			return i
		}
	}

	return -1
}

func main() {
	// Initialize old investments and the new fund splits
	oldInvestments := initOldInvestments()
	printInvestmentsInfo(oldInvestments)
	newSplits := initNewSplits()
	printSplitsInfo(newSplits)

	// Resolve overlapping funds
	completedSwitchInvestments := resolveIntersect(&oldInvestments, &newSplits)
	printSwitchInvestmentsInfo(completedSwitchInvestments)

	// Old investments and new fund splits after resolving the basic
	println("\nAfter resolved intersection:\n")
	printInvestmentsInfo(oldInvestments)
	printSplitsInfo(newSplits)

	// Resolve the rest
	createdSwitchInvestments := resolveSimmetricDifference(&oldInvestments, &newSplits)
	printSwitchInvestmentsInfo(createdSwitchInvestments)

	// Old investments and new fund splits after resolving the basic
	println("\nAfter evrything resolved:\n")
	printInvestmentsInfo(oldInvestments)
	printSplitsInfo(newSplits)
}

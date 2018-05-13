package main

import (
	"reflect"
	"testing"
)

func TestResolveIntersectWithNoOverlapping(t *testing.T) {
	t.Log("Investments do not have overlapping funds with the new splits. (Expected no switch investments)")

	investments := InvestmentSlice{
		Investment{fund: 1, percentage: 50, units: 1.45},
		Investment{fund: 2, percentage: 50, units: 2.50},
	}
	splits := SplitSlice{Split{fund: 3, percentage: 100}}

	if switchInvs := resolveIntersect(&investments, &splits); len(switchInvs) != 0 {
		t.Error("Expected no switch investments, but get", switchInvs)
	}
}

func TestResolveIntersectWithPartOverlapping(t *testing.T) {
	t.Log("Investments have some overlapping funds with the new splits. (Expected one switch investments)")

	investments := InvestmentSlice{
		Investment{fund: 1, percentage: 50, units: 1.45},
		Investment{fund: 2, percentage: 50, units: 2.50},
	}
	splits := SplitSlice{
		Split{fund: 3, percentage: 70}, Split{fund: 2, percentage: 30},
	}

	expected := []SwitchInvestment{
		SwitchInvestment{from_fund: 2, to_fund: 2, units: 1.5, state: "completed"},
	}

	if switchInvs := resolveIntersect(&investments, &splits); !reflect.DeepEqual(switchInvs, expected) {
		t.Error("Expected", expected, ", but get", switchInvs)
	}
}

func TestResolveIntersectWithAllOverlapping(t *testing.T) {
	t.Log("Investments have exatly the same funds with the new splits. (Expected two switch investments)")

	investments := InvestmentSlice{
		Investment{fund: 1, percentage: 50, units: 1.45},
		Investment{fund: 2, percentage: 50, units: 2.50},
	}
	splits := SplitSlice{
		Split{fund: 1, percentage: 90}, Split{fund: 2, percentage: 10},
	}

	expected := []SwitchInvestment{
		SwitchInvestment{from_fund: 1, to_fund: 1, units: 1.45, state: "completed"},
		SwitchInvestment{from_fund: 2, to_fund: 2, units: 0.5, state: "completed"},
	}

	if switchInvs := resolveIntersect(&investments, &splits); !reflect.DeepEqual(switchInvs, expected) {
		t.Error("Expected", expected, ", but get", switchInvs)
	}
}

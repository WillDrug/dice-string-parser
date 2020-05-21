package diceparser

import (
	"testing"
)

func TestParse(t *testing.T) {
	var parsed RollParse
	var err error
	parsed, err = Parse("1d100-10+5d1l1")
	if err != nil {
		t.Error(err)
	}
	if (parsed.Rolls()[0].total != parsed.rolls[0].total) {
		t.Error("Methods failed")
	}
	if (parsed.arithmetic[0] != -1 || parsed.arithmetic[1] != 1) {
		t.Error("Arithmetic failed")
	}
	if (parsed.rolls[0].dice != 100 || parsed.rolls[0].num != 1) {
		t.Error("Roll 1 failed")
	}
	if (parsed.rolls[1].dice != 1 || parsed.rolls[1].total != 10) {
		t.Error("Roll 2 failed")
	}
	if (parsed.rolls[2].dice != 1 || parsed.rolls[2].num != 5 || parsed.rolls[2].total != 1) {
		t.Error("Roll 3 failed")
	}
}

package core__test

import (
	"GoFyS/core"
	"testing"
)

var bytetobintests = []struct {
	name    string
	byteArr []byte
	binary  []bool
}{
	{
		name:    "zeros",
		byteArr: []byte{'\x00'},
		binary:  []bool{false, false, false, false, false, false, false, false},
	},
	{
		name:    "mixed",
		byteArr: []byte{'\x75'},
		binary:  []bool{false, true, true, true, false, true, false, true},
	},
	{
		name:    "ones",
		byteArr: []byte{'\xFF'},
		binary:  []bool{true, true, true, true, true, true, true, true},
	},
	{
		name:    "multiple",
		byteArr: []byte{'\x2E', '\x3D', '\x7E'},
		binary: []bool{
			false, false, true, false, true, true, true, false,
			false, false, true, true, true, true, false, true,
			false, true, true, true, true, true, true, false,
		},
	},
}

func TestByteToBin(t *testing.T) {
	for _, btbt := range bytetobintests {
		t.Run(btbt.name, func(t *testing.T) {
			result := core.ByteToBin(btbt.byteArr)
			haveError := false
			for i, val := range result {
				if val != btbt.binary[i] {
					haveError = true
				}
			}
			if haveError {
				t.Errorf("got %t, want %t", result, btbt.binary)
			}
		})
	}
}

func TestBinToByteValue(t *testing.T) {
	for _, btbt := range bytetobintests[:3] {
		t.Run(btbt.name, func(t *testing.T) {
			result := core.BinToByteValue(btbt.binary)
			if result != btbt.byteArr[0] {
				t.Errorf("got %q, want %q", result, btbt.byteArr)
			}
		})
	}
}

package coreunit__test

import (
	"testing"

	"github.com/bojodimitrov/gofys/core"
)

var storageLen = 5

var writetests = []struct {
	name     string
	content  string
	expected []byte
	offset   int
}{
	{
		name:     "ordinary",
		content:  "aaa",
		expected: []byte{'\x61', '\x61', '\x61', '\x00', '\x00'},
		offset:   0,
	},
	{
		name:     "fill storage",
		content:  "_full",
		expected: []byte{'\x5F', '\x66', '\x75', '\x6C', '\x6C'},
		offset:   0,
	},
	{
		name:     "middle out",
		content:  "out",
		expected: []byte{'\x00', '\x00', '\x6F', '\x75', '\x74'},
		offset:   2,
	},
	{
		name:     "no content",
		content:  "",
		expected: []byte{'\x00', '\x00', '\x00', '\x00', '\x00'},
		offset:   0,
	},
}

func TestWrite(t *testing.T) {
	for _, wt := range writetests {
		t.Run(wt.name, func(t *testing.T) {
			storage := make([]byte, storageLen)
			core.Write(storage, wt.content, wt.offset)
			haveError := false
			for i, val := range storage {
				if val != wt.expected[i] {
					haveError = true
				}
			}
			if haveError {
				t.Errorf("got %q, want %q", storage, wt.expected)
			}
		})
	}
}

var writetestspanic = []struct {
	name     string
	content  string
	expected []byte
	offset   int
}{
	{
		name:     "content too long",
		content:  "too long",
		expected: []byte{},
		offset:   0,
	},
	{
		name:     "offset out of bound",
		content:  "1",
		expected: []byte{},
		offset:   6,
	},
}

func TestWritePanic(t *testing.T) {
	for _, wtp := range writetestspanic {
		t.Run(wtp.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("The code did not panic")
				}
			}()
			storage := make([]byte, storageLen)
			core.Write(storage, wtp.content, wtp.offset)
		})
	}
}

var writebytetests = []struct {
	name     string
	content  []byte
	expected []byte
	offset   int
}{
	{
		name:     "ordinary",
		content:  []byte{'\x61', '\x61', '\x61'},
		expected: []byte{'\x61', '\x61', '\x61', '\x00', '\x00'},
		offset:   0,
	},
	{
		name:     "fill storage",
		content:  []byte{'\x5F', '\x66', '\x75', '\x6C', '\x6C'},
		expected: []byte{'\x5F', '\x66', '\x75', '\x6C', '\x6C'},
		offset:   0,
	},
	{
		name:     "middle out",
		content:  []byte{'\x6F', '\x75', '\x74'},
		expected: []byte{'\x00', '\x00', '\x6F', '\x75', '\x74'},
		offset:   2,
	},
	{
		name:     "no content",
		content:  []byte{},
		expected: []byte{'\x00', '\x00', '\x00', '\x00', '\x00'},
		offset:   0,
	},
}

func TestWriteByte(t *testing.T) {
	for _, wbt := range writebytetests {
		t.Run(wbt.name, func(t *testing.T) {
			storage := make([]byte, storageLen)
			core.WriteByte(storage, wbt.content, wbt.offset)
			haveError := false
			for i, val := range storage {
				if val != wbt.expected[i] {
					haveError = true
				}
			}
			if haveError {
				t.Errorf("got %q, want %q", storage, wbt.expected)
			}
		})
	}
}

var writebytetestspanic = []struct {
	name     string
	content  []byte
	expected []byte
	offset   int
}{
	{
		name:     "content too long",
		content:  []byte{'\x74', '\x64', '\x64', '\x20', '\x6C', '\x64', '\x6E', '\x67'},
		expected: []byte{},
		offset:   0,
	},
	{
		name:     "offset out of bound",
		content:  []byte{'\x31'},
		expected: []byte{},
		offset:   6,
	},
}

func TestWriteBytePanic(t *testing.T) {
	for _, wtbp := range writebytetestspanic {
		t.Run(wtbp.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("The code did not panic")
				}
			}()
			storage := make([]byte, storageLen)
			core.WriteByte(storage, wtbp.content, wtbp.offset)
		})
	}
}

var readtests = []struct {
	name     string
	expected string
	offset   int
	length   int
}{
	{
		name:     "ordinary",
		expected: "_re",
		offset:   0,
		length:   3,
	},
	{
		name:     "full",
		expected: "_read",
		offset:   0,
		length:   5,
	},
	{
		name:     "middle out",
		expected: "ad",
		offset:   3,
		length:   2,
	},
}

var storageТоRead = []byte{'\x5F', '\x72', '\x65', '\x61', '\x64'}

func TestRead(t *testing.T) {
	for _, rt := range readtests {
		t.Run(rt.name, func(t *testing.T) {
			storageRead := core.Read(storageТоRead, rt.offset, rt.length)
			haveError := false
			for i, _ := range storageRead {
				if storageRead[i] != rt.expected[i] {
					haveError = true
				}
			}
			if haveError {
				t.Errorf("got %q, want %q", storageRead, rt.expected)
			}
		})
	}
}

var readtestspanic = []struct {
	name     string
	expected string
	offset   int
	length   int
}{
	{
		name:     "negative length",
		expected: "",
		offset:   0,
		length:   -1,
	},
	{
		name:     "offset out of bound",
		expected: "",
		offset:   6,
		length:   1,
	},
	{
		name:     "length out of bound",
		expected: "",
		offset:   0,
		length:   6,
	},
}

func TestReadPanic(t *testing.T) {
	for _, rtp := range readtestspanic {
		t.Run(rtp.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("The code did not panic")
				}
			}()
			core.Read(storageТоRead, rtp.offset, rtp.length)
		})
	}
}

func TestReadRaw(t *testing.T) {
	rawStorage := []byte{'\x00', '\x72', '\x65', '\x61', '\x64'}
	rawExpected := []byte{'\x00', '\x72', '\x65'}
	rawStorageRead := core.ReadRaw(rawStorage, 0, 3)
	haveError := false
	for i, _ := range rawStorageRead {
		if rawStorageRead[i] != rawExpected[i] {
			haveError = true
		}
	}
	if haveError {
		t.Errorf("got %q, want %q", rawStorageRead, rawExpected)
	}
}

func TestReadByte(t *testing.T) {
	rawStorage := []byte{'\x00', '\x72', '\x65', '\x61', '\x64'}
	rawExpected := []byte{'\x00', '\x72', '\x65'}
	rawStorageRead := core.ReadByte(rawStorage, 0, 3)
	haveError := false
	for i, _ := range rawStorageRead {
		if rawStorageRead[i] != rawExpected[i] {
			haveError = true
		}
	}
	if haveError {
		t.Errorf("got %q, want %q", rawStorageRead, rawExpected)
	}
}

package diractsunit__test

import (
	"testing"

	"github.com/bojodimitrov/byfiri/diracts"
	s "github.com/bojodimitrov/byfiri/structures"
)

var encodetests = []struct {
	name    string
	decoded []s.DirectoryEntry
	encoded string
}{
	{
		name:    "ordinary",
		decoded: []s.DirectoryEntry{s.DirectoryEntry{Inode: 1, FileName: "abv.txt"}},
		encoded: "1\x00\x00\x00\x00\x00\x00\x00\x00\x007\x00\x00abv.txt"},
	{
		name:    "long inode",
		decoded: []s.DirectoryEntry{s.DirectoryEntry{Inode: 1234567890, FileName: "abv.txt"}},
		encoded: "12345678907\x00\x00abv.txt"},
	{
		name: "more than file",
		decoded: []s.DirectoryEntry{
			s.DirectoryEntry{Inode: 1234567890, FileName: "."},
			s.DirectoryEntry{Inode: 1234567891, FileName: "abc.txt"},
			s.DirectoryEntry{Inode: 1234567892, FileName: "somefilename.txt"},
		},
		encoded: "12345678901\x00\x00.12345678917\x00\x00abc.txt123456789216\x00somefilename.txt",
	},
	{
		name:    "long file name",
		decoded: []s.DirectoryEntry{s.DirectoryEntry{Inode: 123, FileName: "thisisanamewithoveronehundredcharactersreallyreallyreallylongnamethisisandIdonotknowwhyyoustillreadthename.txt"}},
		encoded: "123\x00\x00\x00\x00\x00\x00\x00110thisisanamewithoveronehundredcharactersreallyreallyreallylongnamethisisandIdonotknowwhyyoustillreadthename.txt",
	},
}

//TestEncodeDirectoryContent tests EncodeDirectoryContent function
func TestEncodeDirectoryContent(t *testing.T) {
	for _, et := range encodetests {
		t.Run(et.name, func(t *testing.T) {
			result, err := diracts.EncodeDirectoryContent(et.decoded)
			if result != et.encoded || err != nil {
				t.Errorf("got %q, want %q", result, et.encoded)
			}
		})
	}
}

var veryLongFileName string = "ThreeRingsfortheElvenkingsundertheskySevenfortheDwarflordsintheirhallsofstoneNineforMortalMendoomedtodieOnefortheDarkLordonhisdarkthroneIntheLandofMordorwheretheShadowslieOneRingtorulethemallOneRingtofindthemOneRingtobringthemallandinthedarknessbindthemIntheLandofMordorwheretheShadowslie"

var encodetestserrors = []struct {
	name    string
	decoded []s.DirectoryEntry
	encoded string
}{
	{
		name:    "empty file name",
		decoded: []s.DirectoryEntry{s.DirectoryEntry{Inode: 1, FileName: ""}},
		encoded: "",
	},
	{
		name:    "too long file name",
		decoded: []s.DirectoryEntry{s.DirectoryEntry{Inode: 1, FileName: veryLongFileName}},
		encoded: "",
	},
}

//TestEncodeDirectoryContentErrors tests EncodeDirectoryContent errors returned
func TestEncodeDirectoryContentErrors(t *testing.T) {
	for _, et := range encodetestserrors {
		t.Run(et.name, func(t *testing.T) {
			result, err := diracts.EncodeDirectoryContent(et.decoded)
			if result != et.encoded || err == nil {
				t.Errorf("got %q, want %q", result, et.encoded)
			}
		})
	}
}

//TestDecodeDirectoryContent tests DecodeDirectoryContent
func TestDecodeDirectoryContent(t *testing.T) {
	for _, et := range encodetests {
		t.Run(et.name, func(t *testing.T) {
			result := diracts.DecodeDirectoryContent(et.encoded)
			for i, dir := range result {
				if dir != et.decoded[i] {
					t.Errorf("got %q, want %q", result, et.encoded)
				}
			}
		})
	}
}

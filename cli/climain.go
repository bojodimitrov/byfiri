package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bojodimitrov/byfiri/core"
	"github.com/bojodimitrov/byfiri/structures"
)

//Start contains mainloop of cli
func Start(storage []byte, currentDirectory *structures.DirectoryIterator) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println()
		fmt.Println(core.GetPath(storage, currentDirectory))
		fmt.Print("->")
		scanner.Scan()
		commands := parseInput(scanner.Text())

		action := parseCommand(commands[0])

		if action == nil {
			fmt.Println("unknown command")
		} else {
			currentDirectory = action(storage, currentDirectory, commands)
		}
	}
}

func parseCommand(text string) func([]byte, *structures.DirectoryIterator, []string) *structures.DirectoryIterator {
	switch text {
	case Open.String():
		return open
	case Make.String():
		return make
	case Delete.String():
		return delete
	case Edit.String():
		return edit
	case Rename.String():
		return rename
	case List.String():
		return list
	case Tree.String():
		return tree
	case Exit.String():
		return exit
	default:
		return nil
	}
}

func getFileContent() string {
	scanner := bufio.NewScanner(os.Stdin)
	var content strings.Builder
	for {
		fmt.Print(">")
		scanner.Scan()
		if scanner.Text() == ":q" || scanner.Text() == ":Q" {
			break
		} else {
			content.WriteString(scanner.Text())
			content.WriteString("\n")
		}
	}
	return content.String()
}

func parseInput(command string) []string {
	if strings.Count(command, "\"")%2 == 1 {
		return []string{""}
	}
	var content strings.Builder
	fields := []string{}
	ignoreSpaces := false
	for _, char := range command {
		if char == ' ' && !ignoreSpaces {
			fields = append(fields, content.String())
			content.Reset()
		}
		if char == '"' {
			ignoreSpaces = !ignoreSpaces
		}
		content.WriteRune(char)
	}
	fields = append(fields, content.String())
	for i, field := range fields {
		fields[i] = strings.Trim(field, " \"")
	}
	return fields
}

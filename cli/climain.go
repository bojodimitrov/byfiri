package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bojodimitrov/byfiri/core"
	"github.com/bojodimitrov/byfiri/graphic"
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
		// text = text[:len(text)-2]
		commands := strings.Fields(scanner.Text())

		action := parseCommand(commands[0])

		if action == nil {
			fmt.Println("unrecognized command")
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

func list(storage []byte, currentDirectory *structures.DirectoryIterator, _ []string) *structures.DirectoryIterator {
	for _, entry := range currentDirectory.DirectoryContent {
		fmt.Println(entry.FileName)
	}
	return currentDirectory
}

func tree(storage []byte, currentDirectory *structures.DirectoryIterator, _ []string) *structures.DirectoryIterator {
	graphic.DisplayDirectoryTree(storage, currentDirectory)
	return currentDirectory
}

func exit(_ []byte, currentDirectory *structures.DirectoryIterator, _ []string) *structures.DirectoryIterator {
	os.Exit(0)
	return currentDirectory
}

func open(storage []byte, currentDirectory *structures.DirectoryIterator, commands []string) *structures.DirectoryIterator {
	var err error
	if strings.Contains(commands[1], "\\") {
		currentDirectory, err = core.EnterDirectory(storage, currentDirectory, commands[1])
	} else {
		if core.IsDirectory(storage, currentDirectory, commands[1]) {
			currentDirectory, err = core.EnterDirectory(storage, currentDirectory, commands[1])
		} else {
			fmt.Print(core.ReadFile(storage, core.GetInode(storage, currentDirectory, commands[1])))
		}
	}
	if err != nil {
		fmt.Println(err)
	}
	return currentDirectory
}

func edit(storage []byte, currentDirectory *structures.DirectoryIterator, commands []string) *structures.DirectoryIterator {
	if core.IsDirectory(storage, currentDirectory, commands[1]) {
		fmt.Println("edit: cannot edit directory")
	} else {
		core.UpdateFile(storage, core.GetInode(storage, currentDirectory, commands[1]), getFileContent())
	}

	return currentDirectory
}

func make(storage []byte, currentDirectory *structures.DirectoryIterator, commands []string) *structures.DirectoryIterator {
	switch commands[1] {
	case File.String():
		core.AllocateFile(storage, currentDirectory, commands[2], getFileContent())
	case Directory.String():
		core.AllocateDirectory(storage, currentDirectory, commands[2])
	default:
		fmt.Println("unrecognized command")
		return currentDirectory
	}

	return currentDirectory
}

func delete(storage []byte, currentDirectory *structures.DirectoryIterator, commands []string) *structures.DirectoryIterator {
	if core.IsDirectory(storage, currentDirectory, commands[1]) {
		core.DeleteDirectory(storage, currentDirectory, core.GetInode(storage, currentDirectory, commands[1]))
	} else {
		core.DeleteFile(storage, currentDirectory, core.GetInode(storage, currentDirectory, commands[1]))
	}
	return currentDirectory
}

func rename(storage []byte, currentDirectory *structures.DirectoryIterator, commands []string) *structures.DirectoryIterator {
	core.RenameFile(storage, currentDirectory, core.GetInode(storage, currentDirectory, commands[1]), commands[2])
	return currentDirectory
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
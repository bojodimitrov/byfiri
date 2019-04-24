package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/bojodimitrov/byfiri/core"
	"github.com/bojodimitrov/byfiri/graphic"
	"github.com/bojodimitrov/byfiri/structures"
)

const helpOption = "-h"

func list(storage []byte, currentDirectory *structures.DirectoryIterator, commands []string) *structures.DirectoryIterator {
	if len(commands) > 2 {
		fmt.Println("too many parameters:", commands[1:])
		return currentDirectory
	}
	if len(commands) == 2 {
		if strings.HasPrefix(commands[1], "-") {
			printHelp(commands, listHelp)
		} else {
			fmt.Println("too many parameters:", commands[1:])
		}
		return currentDirectory
	}
	for _, entry := range currentDirectory.DirectoryContent {
		fmt.Println(entry.FileName)
	}
	return currentDirectory
}

func tree(storage []byte, currentDirectory *structures.DirectoryIterator, commands []string) *structures.DirectoryIterator {
	if len(commands) > 2 {
		fmt.Println("too many parameters:", commands[1:])
		return currentDirectory
	}
	if len(commands) == 2 {
		if strings.HasPrefix(commands[1], "-") {
			printHelp(commands, treeHelp)
		} else {
			fmt.Println("too many parameters:", commands[1:])
		}
		return currentDirectory
	}
	graphic.DisplayDirectoryTree(storage, currentDirectory)
	return currentDirectory
}

func exit(_ []byte, currentDirectory *structures.DirectoryIterator, commands []string) *structures.DirectoryIterator {
	if len(commands) > 2 {
		fmt.Println("too many parameters:", commands[1:])
		return currentDirectory
	}
	if len(commands) == 2 {
		if strings.HasPrefix(commands[1], "-") {
			printHelp(commands, exitHelp)
		} else {
			fmt.Println("too many parameters:", commands[1:])
		}
		return currentDirectory
	}
	os.Exit(0)
	return currentDirectory
}

func open(storage []byte, currentDirectory *structures.DirectoryIterator, commands []string) *structures.DirectoryIterator {
	if len(commands) > 2 {
		fmt.Println("too many parameters:", commands[2:])
		return currentDirectory
	}
	if len(commands) == 2 && strings.HasPrefix(commands[1], "-") {
		printHelp(commands, openHelp)
		return currentDirectory
	}
	var err error
	if strings.Contains(commands[1], "\\") {
		preserveCurrent := currentDirectory
		path := strings.Split(commands[1], "\\")
		for _, value := range path {
			isDir, err := core.IsDirectory(storage, currentDirectory, commands[1])
			if err != nil {
				//Log err
				fmt.Println("open: could not read file")
				return currentDirectory
			}
			if isDir {
				currentDirectory, err = core.EnterDirectory(storage, currentDirectory, value)
			} else {
				fmt.Print(core.ReadFile(storage, core.GetInode(storage, currentDirectory, value)))
				return preserveCurrent
			}
		}
	} else {
		isDir, err := core.IsDirectory(storage, currentDirectory, commands[1])
		if err != nil {
			//Log err
			fmt.Println("open: could not read file")
			return currentDirectory
		}
		if isDir {
			currentDirectory, err = core.EnterDirectory(storage, currentDirectory, commands[1])
		} else {
			content, err := core.ReadFile(storage, core.GetInode(storage, currentDirectory, commands[1]))
			if err != nil {
				fmt.Println(err)
			}
			fmt.Print(content)
		}
	}
	if err != nil {
		fmt.Println(err)
	}
	return currentDirectory
}

func edit(storage []byte, currentDirectory *structures.DirectoryIterator, commands []string) *structures.DirectoryIterator {
	if len(commands) > 2 {
		fmt.Println("too many parameters:", commands[2:])
		return currentDirectory
	}
	if len(commands) == 2 && strings.HasPrefix(commands[1], "-") {
		printHelp(commands, editHelp)
		return currentDirectory
	}
	value, err := core.IsDirectory(storage, currentDirectory, commands[1])
	if err != nil {
		//Log err
		fmt.Println("edit: could not read file")
		return currentDirectory
	}
	if value {
		fmt.Println("edit: cannot edit directory")
		return currentDirectory
	}
	core.UpdateFile(storage, core.GetInode(storage, currentDirectory, commands[1]), getFileContent())

	return currentDirectory
}

func make(storage []byte, currentDirectory *structures.DirectoryIterator, commands []string) *structures.DirectoryIterator {
	if len(commands) == 1 {
		fmt.Println("option must be provided: 'f' for file and 'd' for directory")
		fmt.Println("name must be provided")
		return currentDirectory
	}
	if len(commands) == 2 {
		if commands[1] == helpOption {
			fmt.Println(makeHelp)
		} else {
			fmt.Println("name must be provided")
		}
		return currentDirectory
	}
	if len(commands) > 3 {
		fmt.Println("too many parameters:", commands[3:])
		return currentDirectory
	}
	switch commands[1] {
	case File.String():
		_, err := core.AllocateFile(storage, currentDirectory, commands[2], getFileContent())
		if err != nil {
			fmt.Println(err)
		}
	case Directory.String():
		core.AllocateDirectory(storage, currentDirectory, commands[2])
	default:
		fmt.Println("unknown option: ", commands[1])
		return currentDirectory
	}

	return currentDirectory
}

func delete(storage []byte, currentDirectory *structures.DirectoryIterator, commands []string) *structures.DirectoryIterator {
	if len(commands) == 1 {
		fmt.Println("name not provided")
		return currentDirectory
	}
	if len(commands) > 2 {
		fmt.Println("too many parameters:", commands[2:])
		return currentDirectory
	}
	if len(commands) == 2 && strings.HasPrefix(commands[1], "-") {
		printHelp(commands, deleteHelp)
		return currentDirectory
	}
	isDir, err := core.IsDirectory(storage, currentDirectory, commands[1])
	if err != nil {
		//Log err
		fmt.Println("delete: could not delete file")
		return currentDirectory
	}
	if isDir {
		core.DeleteDirectory(storage, currentDirectory, core.GetInode(storage, currentDirectory, commands[1]))
	} else {
		core.DeleteFile(storage, currentDirectory, core.GetInode(storage, currentDirectory, commands[1]))
	}
	return currentDirectory
}

func rename(storage []byte, currentDirectory *structures.DirectoryIterator, commands []string) *structures.DirectoryIterator {
	if len(commands) == 1 {
		fmt.Println("old and new file names must be provided")
		fmt.Println("new file name must be provided")
		return currentDirectory
	}
	if len(commands) == 2 && strings.HasPrefix(commands[1], "-") {
		printHelp(commands, renameHelp)
		return currentDirectory
	}
	if len(commands) == 2 {
		fmt.Println("new file name must be provided")
		return currentDirectory
	}
	if len(commands) > 3 {
		fmt.Println("too many parameters:", commands[3:])
		return currentDirectory
	}
	core.RenameFile(storage, currentDirectory, core.GetInode(storage, currentDirectory, commands[1]), commands[2])
	return currentDirectory
}

func printHelp(commands []string, help string) {
	if commands[1] == helpOption {
		fmt.Println(help)
	}
}

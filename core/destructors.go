package core

import (
	"fmt"

	"github.com/bojodimitrov/byfiri/structures"
)

func deleteInode(storage []byte, fsdata *structures.Metadata, inode int) {
	markOnBitmap(storage, fsdata, false, inode, structures.Inodes)
}

func deleteBlocks(storage []byte, fsdata *structures.Metadata, inodeInfo *structures.Inode) {
	for _, block := range inodeInfo.BlocksLocations {
		markOnBitmap(storage, fsdata, false, int(block), structures.Blocks)
	}
}

func removeFileFromDirectory(storage []byte, current *structures.DirectoryIterator, fileInode int) {
	content, err := ReadDirectory(storage, current.DirectoryInode)
	if err != nil {
		fmt.Print(err)
	}

	index := 0
	for i, entry := range content {
		if entry.Inode == uint32(fileInode) {
			index = i
		}
	}

	content[index] = content[len(content)-1]
	content = content[:len(content)-1]

	UpdateDirectory(storage, current.DirectoryInode, content)

	current.DirectoryContent[index] = current.DirectoryContent[len(current.DirectoryContent)-1]
	current.DirectoryContent = current.DirectoryContent[:len(current.DirectoryContent)-1]
}

//DeleteFile deletes file
func DeleteFile(storage []byte, currentDirectory *structures.DirectoryIterator, inode int) error {
	if inode == 0 {
		return fmt.Errorf("delete file: inode cannot be 0")
	}

	fsdata := ReadMetadata(storage)
	inodeInfo, err := ReadInode(storage, fsdata, inode)
	if err != nil {
		return err
	}
	clearFile(storage, inodeInfo.BlocksLocations, fsdata)
	clearInode(storage, fsdata, inode)
	removeFileFromDirectory(storage, currentDirectory, inode)
	deleteBlocks(storage, fsdata, inodeInfo)
	deleteInode(storage, fsdata, inode)
	return nil
}

func directoryContainsDirectory(storage []byte, fsdata *structures.Metadata, currentDirectory *structures.DirectoryIterator) (bool, error) {
	for _, entry := range currentDirectory.DirectoryContent {
		if entry.FileName != "." && entry.FileName != ".." {
			inodeInfo, err := ReadInode(storage, fsdata, int(entry.Inode))
			if err != nil {
				return false, err
			}
			if inodeInfo.Mode == 0 {
				return true, nil
			}
		}
	}
	return false, nil
}

func iterateDirectoryRecursively(storage []byte, fsdata *structures.Metadata, currentDirectory *structures.DirectoryIterator) error {
	//bottom of recursion
	val, err := directoryContainsDirectory(storage, fsdata, currentDirectory)
	if err != nil {
		return err
	}
	if !val {
		// we delete all other files
		for _, entry := range currentDirectory.DirectoryContent {
			if entry.FileName != "." && entry.FileName != ".." {
				DeleteFile(storage, currentDirectory, int(entry.Inode))
			}
		}
		return nil
	}
	for _, entry := range currentDirectory.DirectoryContent {
		if entry.FileName != "." && entry.FileName != ".." {
			inodeInfo, err := ReadInode(storage, fsdata, int(entry.Inode))
			if err != nil {
				return err
			}
			if inodeInfo.Mode == 0 {
				child, err := EnterDirectory(storage, currentDirectory, entry.FileName)
				if err != nil {
					return err
				}
				err = iterateDirectoryRecursively(storage, fsdata, child)
				if err != nil {
					return err
				}
				DeleteFile(storage, currentDirectory, int(entry.Inode))
			}
			if inodeInfo.Mode == 1 {
				DeleteFile(storage, currentDirectory, int(entry.Inode))
			}
		}
	}
	return nil
}

//DeleteDirectory deletes directory recursively
func DeleteDirectory(storage []byte, currentDirectory *structures.DirectoryIterator, inode int) error {
	if inode == 1 {
		return fmt.Errorf("delete file: cannot delete root")
	}

	fsdata := ReadMetadata(storage)
	fileName := ""
	for _, entry := range currentDirectory.DirectoryContent {
		if entry.Inode == uint32(inode) {
			fileName = entry.FileName
		}
	}
	deletedDirectory, err := EnterDirectory(storage, currentDirectory, fileName)
	if err != nil {
		// Log
		return fmt.Errorf("delete directory: could not enter directory")
	}
	err = iterateDirectoryRecursively(storage, fsdata, deletedDirectory)
	if err != nil {
		// Log
		return fmt.Errorf("delete directory: could not delete contents")
	}
	err = DeleteFile(storage, currentDirectory, inode)
	if err != nil {
		// Log
		return fmt.Errorf("delete directory: could not delete directory")
	}
	return nil
}

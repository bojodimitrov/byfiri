package core_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bojodimitrov/byfiri/core"
	"github.com/bojodimitrov/byfiri/structures"
	. "github.com/bojodimitrov/byfiri/structures"
)

var _ = Describe("Mutators", func() {
	fileContent := "content"
	fileUpdatedContent := "update content"

	var (
		storage   []byte
		root      *DirectoryIterator
		err       error
		dirInode  int
		fileInode int
	)

	BeforeEach(func() {
		size := 512 * 1048576
		storage = InitFsSpace(size)
		root, err = AllocateAllStructures(storage, size, DefaultBlockSize)
		Expect(err).NotTo(HaveOccurred())
		var errUpdate error
		dirInode, errUpdate = AllocateDirectory(storage, root, "dir_name")
		Expect(errUpdate).NotTo(HaveOccurred())
		fileInode, errUpdate = AllocateFile(storage, root, "file_name", fileContent)
		Expect(errUpdate).NotTo(HaveOccurred())
	})

	Context("update file", func() {
		It("should update file", func() {
			errUpdate := UpdateFile(storage, fileInode, fileUpdatedContent)
			content, errRead := ReadFile(storage, fileInode)
			Expect(content).To(Equal(fileUpdatedContent))
			Expect(errUpdate).NotTo(HaveOccurred())
			Expect(errRead).NotTo(HaveOccurred())
		})

		It("should update large file", func() {
			largeContent := strings.Repeat("a", DefaultBlockSize+10)
			errUpdate := UpdateFile(storage, fileInode, largeContent)
			content, errRead := ReadFile(storage, fileInode)
			Expect(content).To(Equal(largeContent))
			Expect(errUpdate).NotTo(HaveOccurred())
			Expect(errRead).NotTo(HaveOccurred())
		})
	})

	Context("update directory", func() {
		It("should update directory", func() {
			updatedContent := []structures.DirectoryEntry{
				structures.DirectoryEntry{FileName: ".", Inode: uint32(root.DirectoryInode)},
				structures.DirectoryEntry{FileName: "..", Inode: 0},
				structures.DirectoryEntry{FileName: "test file", Inode: 2},
			}
			errUpdate := UpdateDirectory(storage, dirInode, updatedContent)
			content, errRead := ReadDirectory(storage, dirInode)
			Expect(content).To(Equal(updatedContent))
			Expect(errUpdate).NotTo(HaveOccurred())
			Expect(errRead).NotTo(HaveOccurred())
		})
	})

	Context("update directory unsuccessfully", func() {
		It("should not update inode 0", func() {
			updatedContent := []structures.DirectoryEntry{}
			errUpdate := UpdateDirectory(storage, 0, updatedContent)
			Expect(errUpdate).To(HaveOccurred())
		})

		It("should not update file", func() {
			updatedContent := []structures.DirectoryEntry{}
			errUpdate := UpdateDirectory(storage, fileInode, updatedContent)
			Expect(errUpdate).To(HaveOccurred())
		})

		It("should not update non existant", func() {
			updatedContent := []structures.DirectoryEntry{}
			errUpdate := UpdateDirectory(storage, 100, updatedContent)
			Expect(errUpdate).To(HaveOccurred())
		})
	})

	Context("update file unsuccessfully", func() {
		It("should not update inode 0", func() {
			errUpdate := UpdateFile(storage, 0, fileUpdatedContent)
			Expect(errUpdate).To(HaveOccurred())
		})

		It("should not update directory", func() {
			errUpdate := UpdateFile(storage, dirInode, fileUpdatedContent)
			Expect(errUpdate).To(HaveOccurred())
		})

		It("should not update non existant", func() {
			errUpdate := UpdateFile(storage, 100, fileUpdatedContent)
			Expect(errUpdate).To(HaveOccurred())
		})
	})

	Context("rename file", func() {
		It("should rename file", func() {
			newName := "new_name"
			errRen := RenameFile(storage, root, fileInode, newName)
			rootContent, errRead := ReadDirectory(storage, root.DirectoryInode)
			Expect(errRen).ToNot(HaveOccurred())
			Expect(errRead).ToNot(HaveOccurred())
			rootFiles := []string{}
			for _, val := range rootContent {
				rootFiles = append(rootFiles, val.FileName)
			}
			Expect(rootFiles).To(ContainElement(newName))
		})
	})

	Context("rename file unsuccessfully", func() {
		It("should not rename non existant file", func() {
			newName := "new_name"
			errRen := RenameFile(storage, root, 100, newName)
			Expect(errRen).To(HaveOccurred())
		})
	})
})

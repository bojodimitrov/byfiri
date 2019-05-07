package core_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bojodimitrov/byfiri/core"
	. "github.com/bojodimitrov/byfiri/structures"
)

var _ = Describe("Readers", func() {
	fileContent := "content"

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
		var errRead error
		dirInode, errRead = AllocateDirectory(storage, root, "dir_name")
		Expect(errRead).NotTo(HaveOccurred())
		fileInode, errRead = AllocateFile(storage, root, "file_name", fileContent)
		Expect(errRead).NotTo(HaveOccurred())
	})

	Context("read directory", func() {
		It("should read root with files", func() {
			rootContent, errRead := ReadDirectory(storage, root.DirectoryInode)
			Expect(root.DirectoryContent).To(Equal(rootContent))
			Expect(errRead).NotTo(HaveOccurred())
		})

		It("should read dir content", func() {
			emptyDirContent := []DirectoryEntry{
				DirectoryEntry{
					FileName: ".",
					Inode:    uint32(dirInode),
				},
				DirectoryEntry{
					FileName: "..",
					Inode:    uint32(root.DirectoryInode),
				}}
			dirContent, errRead := ReadDirectory(storage, dirInode)
			Expect(dirContent).To(Equal(emptyDirContent))
			Expect(errRead).NotTo(HaveOccurred())
		})
	})

	Context("read file unsusccessful", func() {
		It("should not read non existant file", func() {
			_, errRead := ReadFile(storage, 100)
			Expect(errRead).To(HaveOccurred())
		})

		It("should not read inode 0", func() {
			_, errRead := ReadFile(storage, 0)
			Expect(errRead).To(HaveOccurred())
		})

		It("should not read directory", func() {
			_, errRead := ReadFile(storage, dirInode)
			Expect(errRead).To(HaveOccurred())
		})
	})

	Context("read directory unsusccessful", func() {
		It("should not read non existant directory", func() {
			_, errRead := ReadDirectory(storage, 100)
			Expect(errRead).To(HaveOccurred())
		})

		It("should not read inode 0", func() {
			_, errRead := ReadDirectory(storage, 0)
			Expect(errRead).To(HaveOccurred())
		})

		It("should not read file", func() {
			_, errRead := ReadDirectory(storage, fileInode)
			Expect(errRead).To(HaveOccurred())
		})
	})

	Context("read file", func() {
		It("should read file", func() {
			fileContent, errRead := ReadFile(storage, fileInode)
			Expect(fileContent).To(Equal(fileContent))
			Expect(errRead).NotTo(HaveOccurred())
		})

	})
})

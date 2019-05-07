package core_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bojodimitrov/byfiri/core"
	. "github.com/bojodimitrov/byfiri/structures"
)

var _ = Describe("Destructors", func() {
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
		var errDelete error
		dirInode, errDelete = AllocateDirectory(storage, root, "dir_name")
		Expect(errDelete).NotTo(HaveOccurred())
		fileInode, errDelete = AllocateFile(storage, root, "file_name", fileContent)
		Expect(errDelete).NotTo(HaveOccurred())
	})

	Context("delete file", func() {
		It("should delete file", func() {
			errDelete := DeleteFile(storage, root, fileInode)
			_, errRead := ReadFile(storage, fileInode)
			Expect(errDelete).NotTo(HaveOccurred())
			Expect(errRead).To(HaveOccurred())
		})
	})

	Context("delete directory", func() {
		It("should delete directory", func() {
			errDelete := DeleteDirectory(storage, root, dirInode)
			_, errRead := ReadDirectory(storage, dirInode)
			Expect(errDelete).NotTo(HaveOccurred())
			Expect(errRead).To(HaveOccurred())
		})
	})

	Context("delete directory unsuccessful", func() {
		It("should not delete root", func() {
			errDelete := DeleteDirectory(storage, root, root.DirectoryInode)
			Expect(errDelete).To(HaveOccurred())
		})

		It("should not delete non existant", func() {
			errDelete := DeleteDirectory(storage, root, 100)
			Expect(errDelete).To(HaveOccurred())
		})
	})

	Context("delete file unsuccessful", func() {
		It("should not delete non existant", func() {
			errDelete := DeleteFile(storage, root, 100)
			Expect(errDelete).To(HaveOccurred())
		})
	})
})

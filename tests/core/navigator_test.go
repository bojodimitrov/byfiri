package coreint_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bojodimitrov/byfiri/core"
	. "github.com/bojodimitrov/byfiri/structures"
)

var _ = Describe("Allocators", func() {
	var (
		storage []byte
		root    *DirectoryIterator
		err     error
	)
	BeforeEach(func() {
		size := 512 * 1048576
		storage = InitFsSpace(size)
		root, err = AllocateAllStructures(storage, size, DefaultBlockSize)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Entering", func() {
		Context("directory", func() {
			It("should be successful", func() {
				dirInode, err := AllocateDirectory(storage, root, "dir")
				Expect(err).NotTo(HaveOccurred())

				dir, err := EnterDirectory(storage, root, "dir")
				Expect(err).NotTo(HaveOccurred())
				Expect(dir.DirectoryInode).To(Equal(dirInode))
			})
		})

		Context("non existant directory", func() {
			It("should not be successful", func() {
				_, err := EnterDirectory(storage, root, "dir")
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

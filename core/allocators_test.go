package core_test

import (
	"strings"

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

	Describe("Allocating", func() {
		Context("basic structures", func() {
			It("should return root", func() {
				expectedInode := 1
				Expect(root.DirectoryInode).To(Equal(expectedInode))
			})
		})

		Context("new file", func() {
			It("should be successful", func() {
				expectedInode := 2
				inode, errAlloc := AllocateFile(storage, root, "file_name", "content")
				Expect(inode).To(Equal(expectedInode))
				Expect(errAlloc).NotTo(HaveOccurred())
			})

			It("should be successful with long content", func() {
				content := strings.Repeat("a", DefaultBlockSize+10)
				expectedInode := 2
				inode, errAlloc := AllocateFile(storage, root, "file_name", content)
				Expect(inode).To(Equal(expectedInode))
				Expect(errAlloc).NotTo(HaveOccurred())
			})

			It("should be successful with multiple files", func() {
				expectedInodes := []int{2, 3, 4}
				inode1, errAlloc := AllocateFile(storage, root, "file_name1", "content")
				Expect(errAlloc).NotTo(HaveOccurred())
				inode2, errAlloc := AllocateFile(storage, root, "file_name2", "content")
				Expect(errAlloc).NotTo(HaveOccurred())
				inode3, errAlloc := AllocateFile(storage, root, "file_name3", "content")
				Expect(errAlloc).NotTo(HaveOccurred())
				Expect(inode1).To(Equal(expectedInodes[0]))
				Expect(inode2).To(Equal(expectedInodes[1]))
				Expect(inode3).To(Equal(expectedInodes[2]))
			})
		})

		Context("duplicate file", func() {
			It("should return duplicate error", func() {
				AllocateFile(storage, root, "file_name", "content")
				_, errAlloc := AllocateFile(storage, root, "file_name", "content")
				Expect(errAlloc).To(HaveOccurred())
			})
		})

		Context("file with forbidden symbols", func() {
			It("should return symbol error", func() {
				_, errAlloc := AllocateFile(storage, root, "file-name", "content")
				Expect(errAlloc).To(HaveOccurred())
			})
		})

		Context("new directory", func() {
			It("should be successful", func() {
				expectedInode := 2
				inode, errAlloc := AllocateDirectory(storage, root, "dir_name")
				Expect(inode).To(Equal(expectedInode))
				Expect(errAlloc).NotTo(HaveOccurred())
			})
		})

		Context("duplicate directory", func() {
			It("should return duplicate error", func() {
				AllocateDirectory(storage, root, "dir_name")
				_, errAlloc := AllocateDirectory(storage, root, "dir_name")
				Expect(errAlloc).To(HaveOccurred())
			})
		})
	})
})

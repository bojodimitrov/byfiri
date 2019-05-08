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

	Describe("Allocating", func() {
		Context("basic structures", func() {
			It("root should contain reference to itself only", func() {
				rootContent, errRead := ReadDirectory(storage, 1)
				Expect(errRead).NotTo(HaveOccurred())
				rootEntry := rootContent[0]
				expected := DirectoryEntry{FileName: ".", Inode: 1}
				Expect(rootEntry).To(Equal(expected))
				nullEntry := rootContent[1]
				expected = DirectoryEntry{FileName: "..", Inode: 0}
				Expect(nullEntry).To(Equal(expected))
			})
		})

		Context("new directory", func() {
			It("should contain references to parent and self", func() {
				inode, errAlloc := AllocateDirectory(storage, root, "dir_name")
				Expect(errAlloc).NotTo(HaveOccurred())
				dirContent, errRead := ReadDirectory(storage, inode)
				Expect(errRead).NotTo(HaveOccurred())
				expected := []DirectoryEntry{
					DirectoryEntry{FileName: ".", Inode: uint32(inode)},
					DirectoryEntry{FileName: "..", Inode: uint32(root.DirectoryInode)},
				}
				Expect(dirContent).To(Equal(expected))
			})
		})

		Context("new files", func() {
			It("directory should cointain inodes of file", func() {
				dirInode, errAlloc := AllocateDirectory(storage, root, "dir_name")
				Expect(errAlloc).NotTo(HaveOccurred())

				fileInode, errAlloc := AllocateFile(storage, root, "test_file", "file content")
				Expect(errAlloc).NotTo(HaveOccurred())

				result, errRead := ReadDirectory(storage, root.DirectoryInode)
				Expect(errRead).NotTo(HaveOccurred())

				expected := []DirectoryEntry{
					DirectoryEntry{FileName: ".", Inode: uint32(root.DirectoryInode)},
					DirectoryEntry{FileName: "..", Inode: 0},
					DirectoryEntry{FileName: "dir_name", Inode: uint32(dirInode)},
					DirectoryEntry{FileName: "test_file", Inode: uint32(fileInode)},
				}
				Expect(result).To(Equal(expected))
			})
		})
	})
})

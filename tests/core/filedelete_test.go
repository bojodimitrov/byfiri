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

	Describe("Deleting", func() {
		Context("file", func() {
			It("directory should not contain reference to it", func() {
				content := "file content"

				fileInode, errAlloc := AllocateFile(storage, root, "test_file", content)
				Expect(errAlloc).NotTo(HaveOccurred())

				resultWithFile, errRead := ReadDirectory(storage, root.DirectoryInode)
				Expect(errRead).NotTo(HaveOccurred())

				expected := []DirectoryEntry{
					DirectoryEntry{FileName: ".", Inode: uint32(root.DirectoryInode)},
					DirectoryEntry{FileName: "..", Inode: 0},
					DirectoryEntry{FileName: "test_file", Inode: uint32(fileInode)},
				}
				Expect(resultWithFile).To(Equal(expected))

				DeleteFile(storage, root, fileInode)
				resultWithNoFile, errRead := ReadDirectory(storage, root.DirectoryInode)
				Expect(errRead).NotTo(HaveOccurred())

				noFile := []DirectoryEntry{
					DirectoryEntry{FileName: ".", Inode: uint32(root.DirectoryInode)},
					DirectoryEntry{FileName: "..", Inode: 0},
				}
				Expect(resultWithNoFile).To(Equal(noFile))
			})
		})

		Context("directory", func() {
			It("files inside should be deleted", func() {
				dirInode, err := AllocateDirectory(storage, root, "root lv1 dir1")
				Expect(err).NotTo(HaveOccurred())

				dir, err := EnterDirectory(storage, root, "root lv1 dir1")
				Expect(err).NotTo(HaveOccurred())
				_, err = AllocateDirectory(storage, dir, "root lv2 dir1")
				Expect(err).NotTo(HaveOccurred())
				dir, err = EnterDirectory(storage, dir, "root lv2 dir1")
				Expect(err).NotTo(HaveOccurred())
				_, err = AllocateFile(storage, dir, "root lv2 f1", "hello there")
				Expect(err).NotTo(HaveOccurred())
				dir, err = EnterDirectory(storage, dir, "..")
				Expect(err).NotTo(HaveOccurred())
				dir, err = EnterDirectory(storage, dir, "..")
				Expect(err).NotTo(HaveOccurred())

				DeleteDirectory(storage, dir, dirInode)
				_, err = ReadDirectory(storage, dirInode)
				Expect(err).To(HaveOccurred())

				fileInode1, err := AllocateFile(storage, dir, "test file", "test content")
				Expect(err).NotTo(HaveOccurred())

				fileInode2, err := AllocateFile(storage, dir, "test file 2", "test content")
				Expect(err).NotTo(HaveOccurred())
				Expect(fileInode1).To(Equal(2))
				Expect(fileInode2).To(Equal(3))
			})
		})
	})
})

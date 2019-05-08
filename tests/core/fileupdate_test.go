package coreint_test

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

	Describe("Updating", func() {
		Context("file", func() {
			It("should be successful", func() {
				content := "file content"
				updatedContent := "update file content"
				fileInode, err := AllocateFile(storage, root, "test file", content)
				Expect(err).NotTo(HaveOccurred())

				resultPure, err := ReadFile(storage, fileInode)
				Expect(err).NotTo(HaveOccurred())

				UpdateFile(storage, fileInode, updatedContent)
				resultTouched, err := ReadFile(storage, fileInode)
				Expect(err).NotTo(HaveOccurred())
				Expect(resultPure).To(Equal(content))
				Expect(resultTouched).To(Equal(updatedContent))
			})

			It("should enlarge successfully", func() {
				content := "file content"
				enlargeContent := strings.Repeat("x", 20000)

				fileInode1, err := AllocateFile(storage, root, "test_file", content)
				Expect(err).NotTo(HaveOccurred())

				fileInode2, err := AllocateFile(storage, root, "test_file2", content)
				Expect(err).NotTo(HaveOccurred())

				UpdateFile(storage, fileInode1, enlargeContent)

				resultFile1, err := ReadFile(storage, fileInode1)
				Expect(err).NotTo(HaveOccurred())

				resultFile2, err := ReadFile(storage, fileInode2)
				Expect(err).NotTo(HaveOccurred())
				Expect(resultFile1).To(Equal(enlargeContent))
				Expect(resultFile2).To(Equal(content))
			})
		})

		Context("directory", func() {
			It("should be successful", func() {

				updatedContent := []DirectoryEntry{
					DirectoryEntry{FileName: ".", Inode: uint32(root.DirectoryInode)},
					DirectoryEntry{FileName: "..", Inode: 0},
					DirectoryEntry{FileName: "test file", Inode: 2},
				}

				expectedPure := []DirectoryEntry{
					DirectoryEntry{FileName: ".", Inode: uint32(root.DirectoryInode)},
					DirectoryEntry{FileName: "..", Inode: 0},
				}

				resultPure, err := ReadDirectory(storage, root.DirectoryInode)
				Expect(err).NotTo(HaveOccurred())

				UpdateDirectory(storage, root.DirectoryInode, updatedContent)
				resultUpdated, err := ReadDirectory(storage, root.DirectoryInode)
				Expect(err).NotTo(HaveOccurred())

				Expect(resultPure).To(Equal(expectedPure))
				Expect(updatedContent).To(Equal(resultUpdated))
			})
		})
	})
})

package core_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bojodimitrov/byfiri/core"
	. "github.com/bojodimitrov/byfiri/structures"
)

var _ = Describe("Navigators", func() {
	var (
		storage   []byte
		root      *DirectoryIterator
		err       error
		dirInode1 int
		dirInode2 int
	)

	BeforeEach(func() {
		size := 512 * 1048576
		storage = InitFsSpace(size)
		root, err = AllocateAllStructures(storage, size, DefaultBlockSize)
		Expect(err).NotTo(HaveOccurred())
		var errEnter error
		dirInode1, errEnter = AllocateDirectory(storage, root, "dir_name1")
		Expect(errEnter).NotTo(HaveOccurred())
		dirInode2, errEnter = AllocateDirectory(storage, root, "dir_name2")
		Expect(errEnter).NotTo(HaveOccurred())
	})

	Context("navigate directory", func() {
		It("should enter first directory", func() {
			dir1, errEnter := EnterDirectory(storage, root, "dir_name1")
			Expect(dir1.DirectoryInode).To(Equal(dirInode1))
			Expect(errEnter).NotTo(HaveOccurred())
		})

		It("should stay in current directory", func() {
			currentDir, errEnter := EnterDirectory(storage, root, ".")
			Expect(currentDir.DirectoryInode).To(Equal(root.DirectoryInode))
			Expect(errEnter).NotTo(HaveOccurred())
		})

		It("should navigate back to root", func() {
			dir1, errEnter := EnterDirectory(storage, root, "dir_name1")
			Expect(dir1.DirectoryInode).To(Equal(dirInode1))
			Expect(errEnter).NotTo(HaveOccurred())
			rootIt, errEnter := EnterDirectory(storage, dir1, "..")
			Expect(rootIt.DirectoryInode).To(Equal(root.DirectoryInode))
			Expect(errEnter).NotTo(HaveOccurred())
		})

		It("should navigate back and forth", func() {
			dir1, errEnter := EnterDirectory(storage, root, "dir_name1")
			Expect(dir1.DirectoryInode).To(Equal(dirInode1))
			Expect(errEnter).NotTo(HaveOccurred())
			rootIt, errEnter := EnterDirectory(storage, dir1, "..")
			Expect(rootIt.DirectoryInode).To(Equal(root.DirectoryInode))
			Expect(errEnter).NotTo(HaveOccurred())
			dir2, errEnter := EnterDirectory(storage, root, "dir_name2")
			Expect(dir2.DirectoryInode).To(Equal(dirInode2))
			Expect(errEnter).NotTo(HaveOccurred())
		})
	})

	Context("navigate directory unsuccessfully", func() {
		It("should not enter root's parent", func() {
			_, errEnter := EnterDirectory(storage, root, "..")
			Expect(errEnter).To(HaveOccurred())
		})

		It("should not enter non existant directory", func() {
			_, errEnter := EnterDirectory(storage, root, "non existant")
			Expect(errEnter).To(HaveOccurred())
		})
	})
})

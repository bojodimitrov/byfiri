package core_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bojodimitrov/byfiri/core"
)

var _ = Describe("StorageIO", func() {
	storageLen := 5
	storageТоRead := []byte{'\x5F', '\x72', '\x65', '\x61', '\x64'}

	var (
		storage []byte
	)

	BeforeEach(func() {
		storage = make([]byte, storageLen)
	})

	Context("write", func() {
		It("should write on storage normal", func() {
			content := "aaa"
			offset := 0
			expected := []byte{'\x61', '\x61', '\x61', '\x00', '\x00'}
			Write(storage, content, offset)
			Expect(storage).To(Equal(expected))
		})
		It("should fill storage", func() {
			content := "_full"
			offset := 0
			expected := []byte{'\x5F', '\x66', '\x75', '\x6C', '\x6C'}
			Write(storage, content, offset)
			Expect(storage).To(Equal(expected))
		})
		It("should write in middle of storage", func() {
			content := "out"
			offset := 2
			expected := []byte{'\x00', '\x00', '\x6F', '\x75', '\x74'}
			Write(storage, content, offset)
			Expect(storage).To(Equal(expected))
		})
		It("should write nothing", func() {
			content := ""
			offset := 0
			expected := []byte{'\x00', '\x00', '\x00', '\x00', '\x00'}
			Write(storage, content, offset)
			Expect(storage).To(Equal(expected))
		})
	})

	Context("write panics", func() {
		It("should not write too long content", func() {
			content := "too long"
			offset := 0
			Expect(func() {
				Write(storage, content, offset)
			}).To(Panic())
		})

		It("should go out of bound", func() {
			content := "1"
			offset := 6
			Expect(func() {
				Write(storage, content, offset)
			}).To(Panic())
		})

	})

	Context("write byte", func() {
		It("should write on storage normal", func() {
			content := []byte{'\x61', '\x61', '\x61'}
			offset := 0
			expected := []byte{'\x61', '\x61', '\x61', '\x00', '\x00'}
			WriteByte(storage, content, offset)
			Expect(storage).To(Equal(expected))
		})
		It("should fill storage", func() {
			content := []byte{'\x5F', '\x66', '\x75', '\x6C', '\x6C'}
			offset := 0
			expected := []byte{'\x5F', '\x66', '\x75', '\x6C', '\x6C'}
			WriteByte(storage, content, offset)
			Expect(storage).To(Equal(expected))
		})
		It("should write in middle of storage", func() {
			content := []byte{'\x6F', '\x75', '\x74'}
			offset := 2
			expected := []byte{'\x00', '\x00', '\x6F', '\x75', '\x74'}
			WriteByte(storage, content, offset)
			Expect(storage).To(Equal(expected))
		})
		It("should write nothing", func() {
			content := []byte{}
			offset := 0
			expected := []byte{'\x00', '\x00', '\x00', '\x00', '\x00'}
			WriteByte(storage, content, offset)
			Expect(storage).To(Equal(expected))
		})
	})

	Context("write byte panics", func() {
		It("should not write too long content", func() {
			content := []byte{'\x74', '\x64', '\x64', '\x20', '\x6C', '\x64', '\x6E', '\x67'}
			offset := 0
			Expect(func() {
				WriteByte(storage, content, offset)
			}).To(Panic())
		})

		It("should go out of bound", func() {
			content := []byte{'\x31'}
			offset := 6
			Expect(func() {
				WriteByte(storage, content, offset)
			}).To(Panic())
		})
	})

	Context("read", func() {
		It("should read from storage normal", func() {
			expected := "_re"
			offset := 0
			length := 3
			storageRead := Read(storageТоRead, offset, length)
			Expect(storageRead).To(Equal(expected))
		})
		It("should read full storage", func() {
			expected := "_read"
			offset := 0
			length := 5
			storageRead := Read(storageТоRead, offset, length)
			Expect(storageRead).To(Equal(expected))
		})
		It("should read middle out", func() {
			expected := "ad"
			offset := 3
			length := 2
			storageRead := Read(storageТоRead, offset, length)
			Expect(storageRead).To(Equal(expected))
		})
	})

	Context("read panics", func() {
		It("should not read negative length", func() {
			offset := 5
			length := -6
			Expect(func() {
				Read(storageТоRead, offset, length)
			}).To(Panic())
		})

		It("should not go out of bound", func() {
			offset := 5
			length := 6
			Expect(func() {
				Read(storageТоRead, offset, length)
			}).To(Panic())
		})
		It("should not read too long", func() {
			offset := 5
			length := 6
			Expect(func() {
				Read(storageТоRead, offset, length)
			}).To(Panic())
		})
	})

	Context("read byte", func() {
		It("should read from storage ", func() {
			rawStorage := []byte{'\x00', '\x72', '\x65', '\x61', '\x64'}
			rawExpected := []byte{'\x00', '\x72', '\x65'}
			offset := 0
			length := 3
			storageRead := ReadByte(rawStorage, offset, length)
			Expect(storageRead).To(Equal(rawExpected))
		})
	})
})

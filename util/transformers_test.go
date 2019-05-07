package util_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bojodimitrov/byfiri/util"
)

var _ = Describe("Transformers", func() {
	Context("byte to bin", func() {
		It("should transform zeros", func() {
			byteArr := []byte{'\x00'}
			binary := []bool{false, false, false, false, false, false, false, false}
			result := ByteToBin(byteArr)
			Expect(result).To(Equal(binary))
		})
		It("should transform mixed", func() {
			byteArr := []byte{'\x75'}
			binary := []bool{false, true, true, true, false, true, false, true}
			result := ByteToBin(byteArr)
			Expect(result).To(Equal(binary))
		})
		It("should transform ones", func() {
			byteArr := []byte{'\xFF'}
			binary := []bool{true, true, true, true, true, true, true, true}
			result := ByteToBin(byteArr)
			Expect(result).To(Equal(binary))
		})
		It("should transform multiple", func() {
			byteArr := []byte{'\x2E', '\x3D', '\x7E'}
			binary := []bool{
				false, false, true, false, true, true, true, false,
				false, false, true, true, true, true, false, true,
				false, true, true, true, true, true, true, false,
			}
			result := ByteToBin(byteArr)
			Expect(result).To(Equal(binary))
		})
	})

	Context("bin to byte value", func() {
		It("should transform zeros", func() {
			byteArr := []byte{'\x00'}
			binary := []bool{false, false, false, false, false, false, false, false}
			result := BinToByteValue(binary)
			Expect(result).To(Equal(byteArr[0]))
		})
		It("should transform mixed", func() {
			byteArr := []byte{'\x75'}
			binary := []bool{false, true, true, true, false, true, false, true}
			result := BinToByteValue(binary)
			Expect(result).To(Equal(byteArr[0]))
		})
		It("should transform ones", func() {
			byteArr := []byte{'\xFF'}
			binary := []bool{true, true, true, true, true, true, true, true}
			result := BinToByteValue(binary)
			Expect(result).To(Equal(byteArr[0]))
		})
	})
})

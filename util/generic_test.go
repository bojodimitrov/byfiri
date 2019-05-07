package util_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bojodimitrov/byfiri/util"
)

var _ = Describe("Generic", func() {
	container := []int{0, 1, 2}

	Context("contains", func() {
		It("should contain", func() {
			element := 1
			result := Contains(container, element)
			Expect(result).To(BeTrue())
		})

		It("should not contain", func() {
			element := 4
			result := Contains(container, element)
			Expect(result).To(BeFalse())
		})
	})

	Context("min", func() {
		It("should find minimum", func() {
			a := 1
			b := 2
			result := Min(a, b)
			Expect(result).To(Equal(a))
		})
	})
})

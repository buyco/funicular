package utils_test

import (
	. "github.com/buyco/funicular/internal/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Slice", func() {

	var (
		testSlice = []string{"foo", "bar"}
	)

	It("should find value in slice", func() {
		exists, index := InArray("foo", testSlice)
		Expect(exists).To(BeTrue())
		Expect(index).To(BeZero())
	})

	It("should not find value in slice", func() {
		exists, index := InArray("woo", testSlice)
		Expect(exists).To(BeFalse())
		Expect(index).To(Equal(-1))
	})
})

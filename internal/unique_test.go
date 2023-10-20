package internal

import (
	"sort"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Unique", func() {
	It("should return unique strings", func() {
		arr := []string{
			"mice",
			"mice",
			"mice",
			"mice",
			"mice",
			"mice",
			"mice",
			"toad",
			"toad",
			"mice",
		}
		result := Unique(arr)
		sort.Strings(result)
		Expect(result).To(Equal([]string{"mice", "toad"}))
	})

	It("should return unique numbers", func() {
		arr := []int{
			1,
			1,
			2,
			3,
			1,
			2,
			3,
		}
		result := Unique(arr)
		sort.Ints(result)
		Expect(result).To(Equal([]int{1, 2, 3}))
	})
})

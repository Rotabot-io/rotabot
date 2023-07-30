package block

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Prefix", func() {
	It("prefixes a string with a value", func() {
		Expect(prefix("foo", "bar")).To(Equal("foo_bar"))
	})
})

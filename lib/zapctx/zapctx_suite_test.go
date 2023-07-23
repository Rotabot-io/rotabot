package zapctx

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestZapctx(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Zapctx Suite")
}

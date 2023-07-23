package codecs

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCodecs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Codecs Suite")
}

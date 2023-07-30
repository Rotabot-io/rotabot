package slackclient

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSlackclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Slackclient Suite")
}

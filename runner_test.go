package runner_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Run Spec")
}

var _ = Describe("", func() {
	var ()

	BeforeEach(func() {
	})

	XIt("", func() {
	})
})

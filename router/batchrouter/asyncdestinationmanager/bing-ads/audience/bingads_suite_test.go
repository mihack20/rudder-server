package bingads_audience

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBingads(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bingads Suite")
}

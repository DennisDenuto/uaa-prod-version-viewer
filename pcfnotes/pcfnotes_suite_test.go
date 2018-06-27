package pcfnotes_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPcfnotes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pcfnotes Suite")
}

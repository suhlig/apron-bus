package atc_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAtc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ATC Suite")
}

package wikipage_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWikipage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wikipage Suite")
}

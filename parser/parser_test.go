package parser_test

import (
	"github.com/mrap/guestimator/models"
	. "github.com/mrap/guestimator/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	Describe("Fetching top wikipages", func() {
		var wikis []models.Wikipage

		BeforeEach(func() {
			wikis = FetchTopWikis()
		})

		It("should return 25 wiki pages", func() {
			Expect(wikis).To(HaveLen(25))
		})
	})
})

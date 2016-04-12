package request_test

import (
	"github.com/mrap/guestimator/models"
	. "github.com/mrap/guestimator/request"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Requests", func() {
	Describe("Requesting for a Wikipage", func() {
		var (
			err  error
			page *models.Wikipage
		)

		Context("When the wikipage does exist", func() {
			BeforeEach(func() {
				page, err = GetWikipage("Christopher Nolan")
			})

			It("should not have an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should have the correct title", func() {
				Expect(page.Title).To(Equal("Christopher Nolan"))
			})

			It("should have a pageid", func() {
				Expect(page.ID).ToNot(BeZero())
			})

			It("should have correct extract text", func() {
				Expect(page.Extract).ToNot(BeEmpty())
			})
		})

		Context("When the wikipage does not exist", func() {
			BeforeEach(func() {
				page, err = GetWikipage("This does not exist")
			})

			It("should not have an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return nil", func() {
				Expect(page).To(BeNil())
			})
		})
	})
})

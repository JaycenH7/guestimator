package models_test

import (
	. "github.com/mrap/guestimator/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WikiPage", func() {

	Describe("Creating a wiki page", func() {
		var wiki WikiPage

		BeforeEach(func() {
			wiki = WikiPage{
				PageID: 1,
				Title:  "My WikiPage",
			}
			err := CreateWikiPage(DB, &wiki)
			Expect(err).NotTo(HaveOccurred())
			Expect(wiki.PageID).NotTo(BeZero())
		})

		It("should be saved in db", func() {
			res, err := GetWikiPage(DB, wiki.PageID)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(&wiki))
		})
	})
})

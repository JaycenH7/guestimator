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

	Describe("Extracting questions out of a wiki page", func() {
		var (
			wiki      WikiPage
			questions []Question
		)

		BeforeEach(func() {
			wiki = WikiPage{
				PageID:  1,
				Title:   "My WikiPage",
				Extract: "",
			}
		})

		Describe("Simple sentence", func() {
			BeforeEach(func() {
				wiki.Extract = "He is 42 years old."
				questions = wiki.ExtractQuestions()
			})

			It("should extract the question correctly", func() {
				Expect(questions).To(HaveLen(1))
				q := questions[0]
				Expect(q.FullText).To(Equal(wiki.Extract))
				Expect(q.Positions).To(Equal([]int{6, 7}))
				Expect(q.PageID).To(Equal(wiki.PageID))
			})
		})
	})
})

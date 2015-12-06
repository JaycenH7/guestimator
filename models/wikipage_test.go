package models_test

import (
	"fmt"

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

		sentences := []Question{
			{FullText: "He is 25 years old.", Positions: []int{6, 7}},
			{FullText: "This sentence has no number.", Positions: []int{}},
			{FullText: "He was born in 1990.", Positions: []int{15, 18}},
		}

		BeforeEach(func() {
			wiki = WikiPage{
				PageID:  1,
				Title:   "My WikiPage",
				Extract: "",
			}
		})

		AssertExtractedQuestion := func(qPos, sPos int) {
			It(fmt.Sprintf("should extract question @%d to match sentence @%d", qPos, sPos), func() {
				q := questions[qPos]
				s := sentences[sPos]
				Expect(q.FullText).To(Equal(s.FullText))
				Expect(q.Positions).To(Equal(s.Positions))
				Expect(q.PageID).To(Equal(wiki.PageID))
			})
		}

		ExtractUntil := func(i int) string {
			extract := ""
			for j := 0; j <= i; j++ {
				extract += sentences[j].FullText + " "
			}
			return extract
		}

		Context(fmt.Sprintf("with extract: %s", ExtractUntil(0)), func() {
			BeforeEach(func() {
				wiki.Extract += ExtractUntil(0)
				questions = wiki.ExtractQuestions()
			})

			It("should have 1 question", func() {
				Expect(questions).To(HaveLen(1))
			})

			AssertExtractedQuestion(0, 0)
		})

		Context(fmt.Sprintf("with extract: %s", ExtractUntil(1)), func() {
			BeforeEach(func() {
				wiki.Extract += ExtractUntil(1)
				questions = wiki.ExtractQuestions()
			})

			It("should have 1 question", func() {
				Expect(questions).To(HaveLen(1))
			})

			AssertExtractedQuestion(0, 0)
		})

		Context(fmt.Sprintf("with extract: %s", ExtractUntil(2)), func() {
			BeforeEach(func() {
				wiki.Extract += ExtractUntil(2)
				questions = wiki.ExtractQuestions()
			})

			It("should have 2 questions", func() {
				Expect(questions).To(HaveLen(2))
			})

			AssertExtractedQuestion(0, 0)
			AssertExtractedQuestion(1, 2)
		})
	})
})

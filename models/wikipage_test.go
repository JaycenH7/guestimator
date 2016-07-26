package models_test

import (
	"fmt"

	. "github.com/mrap/guestimator/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wikipage", func() {

	Describe("Creating a wiki page", func() {
		var wiki Wikipage

		BeforeEach(func() {
			wiki = Wikipage{
				ID:    1,
				Title: "My Wikipage",
			}
			err := CreateWikipage(DB, &wiki)
			Expect(err).NotTo(HaveOccurred())
			Expect(wiki.ID).NotTo(BeZero())
		})

		It("should be saved in db", func() {
			res, err := GetWikipage(DB, wiki.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(&wiki))
		})
	})

	Describe("Extracting questions out of a wiki page", func() {
		var (
			wiki      Wikipage
			questions []Question
		)

		sentences := []Question{
			{FullText: "He is 25 years old.", Positions: []int{6, 7}},
			{FullText: "This sentence has no number.", Positions: []int{}},
			{FullText: "He was born in 1990.", Positions: []int{15, 18}},
			{FullText: "Healthy temp is 98.6.", Positions: []int{16, 19}},
			{FullText: "It is 99.99% effective.", Positions: []int{6, 11}},
		}

		BeforeEach(func() {
			wiki = Wikipage{
				ID:      1,
				Title:   "My Wikipage",
				Extract: "",
			}
		})

		AssertExtractedQuestion := func(qPos, sPos int) {
			It(fmt.Sprintf("should extract question @%d to match sentence @%d", qPos, sPos), func() {
				q := questions[qPos]
				s := sentences[sPos]
				Expect(q.FullText).To(Equal(s.FullText))
				Expect(q.Positions).To(Equal(s.Positions))
				Expect(q.WikipageID).To(Equal(wiki.ID))
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

		Context(fmt.Sprintf("with extract: %s", ExtractUntil(3)), func() {
			BeforeEach(func() {
				wiki.Extract += ExtractUntil(3)
				questions = wiki.ExtractQuestions()
			})

			It("should have 3 questions", func() {
				Expect(questions).To(HaveLen(3))
			})

			AssertExtractedQuestion(0, 0)
			AssertExtractedQuestion(1, 2)
			AssertExtractedQuestion(2, 3)
		})

		Context(fmt.Sprintf("with extract: %s", ExtractUntil(4)), func() {
			BeforeEach(func() {
				wiki.Extract += ExtractUntil(4)
				questions = wiki.ExtractQuestions()
			})

			It("should have 4 questions", func() {
				Expect(questions).To(HaveLen(4))
			})

			AssertExtractedQuestion(0, 0)
			AssertExtractedQuestion(1, 2)
			AssertExtractedQuestion(2, 3)
			AssertExtractedQuestion(3, 4)
		})
	})
})

package models_test

import (
	. "github.com/mrap/guestimator/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Question", func() {

	Describe("Creating a question", func() {
		var (
			wiki     WikiPage
			question Question
		)

		BeforeEach(func() {
			wiki = WikiPage{
				PageID: 1,
				Title:  "My WikiPage",
			}
			err := CreateWikiPage(DB, &wiki)
			Expect(err).NotTo(HaveOccurred())

			question = Question{
				FullText:  "42 has 2 digits",
				Positions: []int{0, 7},
				PageID:    wiki.PageID,
			}
			err = CreateQuestion(DB, &question)
			Expect(err).NotTo(HaveOccurred())
			Expect(question.Id).NotTo(BeZero())
		})

		It("should be saved in db", func() {
			res, err := GetQuestion(DB, question.Id)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(&question))
		})
	})

	Describe("Parsing a string for a question", func() {
		var (
			question *Question
			str      string
		)

		JustBeforeEach(func() {
			question = ParseQuestion(str)
		})

		Context("with one number", func() {
			BeforeEach(func() {
				str = "He is 25."
			})

			It("should set correct positions", func() {
				Expect(question.Positions).To(Equal([]int{6, 7}))
			})
		})

		Context("with two numbers", func() {
			BeforeEach(func() {
				str = "He is 25 not 26."
			})

			It("should set correct positions", func() {
				Expect(question.Positions).To(Equal([]int{6, 7, 13, 14}))
			})
		})
	})

	Describe("Pretty format that shows answers as missing", func() {
		It("should replace the answers with blanks", func() {
			question := Question{
				FullText:  "He will turn 26 in 2016.",
				Positions: []int{13, 14, 19, 22},
			}
			Expect(question.SansAnswers()).To(Equal("He will turn __ in ____."))
		})
	})
})

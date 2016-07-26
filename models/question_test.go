package models_test

import (
	. "github.com/mrap/guestimator/models"
	"github.com/mrap/guestimator/models/fixtures"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Question", func() {

	Describe("Creating a question", func() {
		var (
			wiki     Wikipage
			question Question
		)

		BeforeEach(func() {
			wiki = Wikipage{
				ID:    1,
				Title: "My Wikipage",
			}
			err := CreateWikipage(DB, &wiki)
			Expect(err).NotTo(HaveOccurred())

			question = Question{
				FullText:   "42 has 2 digits",
				Positions:  []int{0, 7},
				WikipageID: wiki.ID,
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

	Describe("Accessing answers", func() {
		var question Question
		BeforeEach(func() {
			question = fixtures.Question()
		})

		It("should be able to return the first answer", func() {
			answer, err := question.FirstAnswer()
			Expect(err).ToNot(HaveOccurred())
			Expect(answer).To(Equal(float64(26)))
		})

		It("should be able to return the second answer", func() {
			answer, err := question.AnswerAt(19)
			Expect(err).ToNot(HaveOccurred())
			Expect(answer).To(Equal(float64(2016)))
		})
	})

	Describe("Sans answers", func() {
		var question Question

		BeforeEach(func() {
			question = fixtures.Question()
		})

		Describe("full text sans answers", func() {
			It("should replace the answers with blanks", func() {
				Expect(question.FullTextSansAnswers()).To(Equal("He will turn __ in ____."))
			})
		})

		Describe("getting a copy rid of any answer-related data", func() {
			var answerless Question

			BeforeEach(func() {
				answerless = question.SansAnswers()
			})

			It("should have full text without the answers", func() {
				Expect(answerless.FullText).To(Equal(question.FullTextSansAnswers()))
			})

			It("should have answerless wikipage", func() {
				Expect(answerless.Wikipage.Extract).To(BeEmpty())
				Expect(answerless.Wikipage.Questions).To(BeEmpty())
			})
		})
	})
})

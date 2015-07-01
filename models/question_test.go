package models_test

import (
	. "github.com/mrap/guestimator/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Question", func() {

	Describe("Creating a question", func() {
		var question Question

		BeforeEach(func() {
			question = Question{
				FullText:  "42 has 2 digits",
				Positions: []int{0, 7},
			}
			err := CreateQuestion(DB, &question)
			Expect(err).NotTo(HaveOccurred())
			Expect(question.Id).NotTo(BeZero())
		})

		It("should be saved in db", func() {
			res, err := GetQuestion(DB, question.Id)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(&question))
		})
	})
})

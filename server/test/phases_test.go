package server_test

import (
	"strconv"

	"github.com/mrap/guestimator/client"
	"github.com/mrap/guestimator/models"
	"github.com/mrap/guestimator/models/fixtures"
	"github.com/mrap/guestimator/server"
	"github.com/mrap/guestimator/server/match"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Match Phases", func() {
	var clients []*client.Client
	var matchID string
	var nextMatchID int

	questions := []models.Question{
		fixtures.Question(),
	}

	// Each run should have a unique match id
	BeforeEach(func() {
		nextMatchID++
		matchID = strconv.Itoa(nextMatchID)
		server.AddMatch(matchID, questions)

		clients = make([]*client.Client, server.MatchSize)
		for i := 0; i < server.MatchSize; i++ {
			clients[i] = client.NewClient(strconv.Itoa(i))
		}
	})

	connect := func(cs ...*client.Client) {
		for _, c := range cs {
			err := connectClient(c, matchID, c.PlayerID)
			Expect(err).ToNot(HaveOccurred())
		}
	}

	Describe("JoinPhase", func() {
		Describe("PlayerJoin events", func() {
			var connClients []*client.Client

			BeforeEach(func() {
				connClients = make([]*client.Client, 0)
			})

			JustBeforeEach(func() {
				connect(connClients...)
			})

			AssertClientsReceivePlayerJoinEvents := func() {
				Specify("just-connected client should receive MatchState message", func() {
					newest := connClients[len(connClients)-1]
					expected := match.Message{
						Type: match.MatchStateMsgType,
						MatchState: &match.MatchState{
							Phase: "Join",
						},
					}

					Eventually(newest.RecvMsg).Should(Receive(Equal(expected)))
				})

				Specify("just-connected client should not receive PlayerJoin event", func() {
					newest := connClients[len(connClients)-1]
					expected := match.Message{
						Type:     match.PlayerJoinMsgType,
						PlayerID: newest.PlayerID,
					}

					Eventually(newest.RecvMsg).ShouldNot(Receive(Equal(expected)))
				})

				Specify("previously connected clients should receive PlayerJoin event", func() {
					newest := connClients[len(connClients)-1]

					expected := match.Message{
						Type:     match.PlayerJoinMsgType,
						PlayerID: newest.PlayerID,
					}

					for _, connC := range connClients[:len(connClients)-1] {
						Eventually(connC.RecvMsg).Should(Receive(Equal(expected)))
					}
				})
			}

			Context("when first player connects", func() {
				BeforeEach(func() {
					connClients = []*client.Client{clients[0]}
				})

				It("should be initially in the JoinPhase", func() {
					cMatch := server.GetMatch(matchID)
					Expect(cMatch).NotTo(BeNil())
					Eventually(func() match.Phase {
						return cMatch.CurrentPhase
					}).Should(BeAssignableToTypeOf(&match.JoinPhase{}))
				})

				AssertClientsReceivePlayerJoinEvents()
			})

			for i, _ := range clients {
				Context("when a player connects", func() {
					BeforeEach(func() {
						connClients = clients[:i+1]
					})

					AssertClientsReceivePlayerJoinEvents()
				})
			}
		})
	})

	Describe("GuessPhase", func() {
		BeforeEach(func() {
			connect(clients...)
		})

		It("should change to GuessPhase", func() {
			cMatch := server.GetMatch(matchID)
			Expect(cMatch).NotTo(BeNil())
			Eventually(func() match.Phase {
				return cMatch.CurrentPhase
			}).Should(BeAssignableToTypeOf(&match.GuessPhase{}))
		})

		It("all players should receive a MatchState message with the question", func() {
			question := fixtures.Question().SansAnswers()
			msg := match.Message{
				Type: match.MatchStateMsgType,
				MatchState: &match.MatchState{
					Phase:    "Guess",
					Question: &question,
				},
			}

			for _, c := range clients {
				Eventually(c.RecvMsg).Should(Receive(Equal(msg)))
			}
		})

		Context("when all players have guessed", func() {
			BeforeEach(func() {
				// All clients send a guess message
				for _, c := range clients {
					c.SendMessage(match.Message{
						Type:  match.GuessMsgType,
						Guess: &match.Guess{0, 0},
					})
				}
			})

			It("should change to GuessResultPhase phase", func() {
				cMatch := server.GetMatch(matchID)
				Expect(cMatch).NotTo(BeNil())
				Eventually(func() match.Phase {
					return cMatch.CurrentPhase
				}).Should(BeAssignableToTypeOf(&match.GuessResultPhase{}))
			})

			Context("when no questions remain", func() {
				It("should change to MatchResultPhase phase", func() {
					cMatch := server.GetMatch(matchID)
					Expect(cMatch).NotTo(BeNil())
					Eventually(func() match.Phase {
						return cMatch.CurrentPhase
					}).Should(BeAssignableToTypeOf(&match.MatchResultPhase{}))
				})

				It("should disconnect all clients", func() {
					for _, c := range clients {
						Eventually(c.RecvMsg).Should(BeClosed())
					}
				})
			})
		})
	})
})

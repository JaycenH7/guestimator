package server_test

import (
	"strconv"

	"github.com/mrap/guestimator/client"
	"github.com/mrap/guestimator/server"
	"github.com/mrap/guestimator/server/match"
	"github.com/mrap/guestimator/server/match/event"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Match Phases", func() {
	var clients [server.MatchSize]client.Client
	var matchID string
	var nextMatchID int

	// Each run should have a unique match id
	BeforeEach(func() {
		matchID = strconv.Itoa(nextMatchID)
		nextMatchID++

		for i := 0; i < server.MatchSize; i++ {
			clients[i] = *client.NewClient(strconv.Itoa(i))
		}
	})

	connect := func(cs ...client.Client) {
		for _, c := range cs {
			connectClient(c, matchID, c.PlayerID)
		}
	}

	Describe("JoinPhase", func() {
		Describe("PlayerJoin events", func() {
			var connClients []client.Client

			BeforeEach(func() {
				connClients = make([]client.Client, 0)
			})

			JustBeforeEach(func() {
				connect(connClients...)
			})

			AssertClientsReceivePlayerJoinEvents := func() {
				Specify("just-connected client should not receive PlayerJoin event", func() {
					newest := connClients[len(connClients)-1]
					Expect(newest.ReceivedEvents).To(BeEmpty())
				})

				Specify("previously connected clients should receive PlayerJoin event", func() {
					newest := connClients[len(connClients)-1]

					expected := event.Event{
						Type:     event.TypePlayerJoin,
						PlayerID: newest.PlayerID,
					}

					for _, connC := range connClients[:len(connClients)-1] {
						Eventually(connC.ReceivedEvents).Should(Receive(Equal(expected)))
					}
				})
			}

			Context("when first player connects", func() {
				BeforeEach(func() {
					connClients = []client.Client{clients[0]}
				})

				It("should be initially in the JoinPhase", func() {
					cMatch := server.GetMatch(matchID)
					Expect(cMatch).NotTo(BeNil())
					Expect(cMatch.CurrentPhase()).To(BeAssignableToTypeOf(match.JoinPhase))
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
})

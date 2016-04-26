package server_test

import (
	"strconv"

	"github.com/mrap/guestimator/client"
	"github.com/mrap/guestimator/server"
	"github.com/mrap/guestimator/server/match"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Match Phases", func() {
	var clients [server.MatchSize]client.Client
	var matchID string
	var nextMatchID int

	// Each run should have a unique match id
	BeforeEach(func() {
		nextMatchID++
		matchID = strconv.Itoa(nextMatchID)
		server.AddMatch(matchID)

		for i := 0; i < server.MatchSize; i++ {
			clients[i] = *client.NewClient(strconv.Itoa(i))
		}
	})

	connect := func(cs ...client.Client) {
		for _, c := range cs {
			if err := connectClient(c, matchID, c.PlayerID); err != nil {
				Expect(err).ToNot(HaveOccurred())
			}
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
				Specify("just-connected client should receive MatchState message", func() {
					newest := connClients[len(connClients)-1]
					expected := match.Message{
						Type: match.MatchStateMsgType,
						MatchState: &match.MatchState{
							Phase: match.JoinPhaseType,
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
					connClients = []client.Client{clients[0]}
				})

				It("should be initially in the JoinPhase", func() {
					cMatch := server.GetMatch(matchID)
					Expect(cMatch).NotTo(BeNil())
					Eventually(func() match.PhaseType {
						return cMatch.PhaseType()
					}).Should(Equal(match.JoinPhaseType))
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

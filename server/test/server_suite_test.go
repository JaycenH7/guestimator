package server_test

import (
	"fmt"
	"log"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/mrap/guestimator/client"
	"github.com/mrap/guestimator/server"
	"github.com/mrap/guestimator/server/match"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

var matchServer *httptest.Server
var matchServerHost string

var _ = BeforeSuite(func() {
	match.PhaseDuration = 50 * time.Millisecond
	SetDefaultEventuallyTimeout(2 * time.Second)
})

var _ = BeforeEach(func() {
	server.ClearMatches()
	matchServer = httptest.NewServer(server.NewMatchHandler())
	matchServerHost = serverHost(*matchServer)
})

var _ = AfterEach(func() {
	matchServer.CloseClientConnections()
	matchServer.Close()
})

func serverHost(s httptest.Server) string {
	u, err := url.Parse(s.URL)
	if err != nil {
		log.Fatalln("Could not get server host from server url:", s.URL, err)
	}
	return u.Host
}

func matchURL(host, matchID, playerID string) string {
	return fmt.Sprintf("ws://%s/match/%s/ws?player=%s", host, matchID, playerID)
}

func connectClient(c *client.Client, matchID, playerID string) error {
	return c.Connect(matchURL(matchServerHost, matchID, playerID))
}

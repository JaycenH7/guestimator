package server_test

import (
	"fmt"
	"log"
	"net/http/httptest"
	"net/url"

	"github.com/mrap/guestimator/client"
	"github.com/mrap/guestimator/server"
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
	server.ClearMatches()
	matchServer = httptest.NewServer(server.NewMatchHandler())
})

var _ = BeforeEach(func() {
	matchServerHost = serverHost(*matchServer)
})

var _ = AfterSuite(func() {
	matchServer.Close()
})

var _ = AfterEach(func() {
	matchServer.CloseClientConnections()
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

func connectClient(c client.Client, matchID, playerID string) error {
	return c.Connect(matchURL(matchServerHost, matchID, playerID))
}

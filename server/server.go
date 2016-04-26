package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrap/guestimator/server/match"
)

const MatchSize = 3

var matches = make(map[string]*match.Match)

func NewMatchHandler() *gin.Engine {
	r := gin.Default()

	r.GET("/match/:id/ws", func(c *gin.Context) {
		id := c.Param("id")

		m, ok := matches[id]
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}

		// TODO: allow player to connect if authorized
		m.Hub.HandleRequest(c.Writer, c.Request)
	})
	return r
}

func AddMatch(matchID string) bool {
	if _, exists := matches[matchID]; exists {
		return false
	}

	matches[matchID] = match.NewMatch(matchID, MatchSize)
	return true
}

func GetMatch(matchID string) *match.Match {
	return matches[matchID]
}

func ClearMatches() {
	matches = make(map[string]*match.Match)
}

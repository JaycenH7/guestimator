package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mrap/guestimator/server/match"
)

const MatchSize = 3

var matches = make(map[string]*match.Match)

func NewMatchHandler() *gin.Engine {
	r := gin.Default()

	r.GET("/match/:id/ws", func(c *gin.Context) {
		id := c.Param("id")

		// get match hub
		m, ok := matches[id]
		if !ok {
			m = match.NewMatch(id, MatchSize)
			matches[id] = m
		}

		// TODO: allow player to connect if authorized
		m.Hub.HandleRequest(c.Writer, c.Request)
	})
	return r
}

func GetMatch(matchID string) *match.Match {
	return matches[matchID]
}

func ClearMatches() {
	matches = make(map[string]*match.Match)
}

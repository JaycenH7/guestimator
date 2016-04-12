package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mrap/guestimator/server/match"
)

var matches = make(map[string]*match.Match)

func main() {
	r := gin.Default()

	r.GET("/match/:id/ws", func(c *gin.Context) {
		id := c.Param("id")

		// get match hub
		m, ok := matches[id]
		if !ok {
			m = match.NewMatch(id, 0)
			matches[id] = m
		}

		// TODO: allow player to connect if authorized
		m.Hub.HandleRequest(c.Writer, c.Request)
	})

	r.Run(":5000")
}

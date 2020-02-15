package controllers

import (
	"github.com/gorilla/websocket"
	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	act(w, r, "Home", func(ctx util.RequestContext) (int, error) {
		return template.Index(ctx, w)
	})
}

var upgrader = websocket.Upgrader{}

func Socket(w http.ResponseWriter, r *http.Request) {
	ctx := util.ExtractContext(r, "WebSocket Connection")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ctx.Logger.Info("Unable to upgrade connection to websocket")
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			ctx.Logger.Info("Unable to read from websocket")
			break
		}
		ctx.Logger.Info("Received message on websocket: " + string(message))
		err = c.WriteMessage(mt, message)
		if err != nil {
			ctx.Logger.Info("Unable to write to websocket")
			break
		}
	}
}

func About(w http.ResponseWriter, r *http.Request) {
	act(w, r, "About", func(ctx util.RequestContext) (int, error) {
		return template.About(ctx, w)
	})
}

package controllers

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
)

func Home(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx util.RequestContext) (int, error) {
		return template.Index(ctx, w)
	})
}

var upgrader = websocket.Upgrader{}

func Socket(w http.ResponseWriter, r *http.Request) {
	ctx := util.ExtractContext(r)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ctx.Logger.Info("Unable to upgrade connection to websocket")
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
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
	act(w, r, func(ctx util.RequestContext) (int, error) {
		ctx.Breadcrumbs = util.BreadcrumbsSimple(ctx.Route("about"), "about")
		return template.About(ctx, w)
	})
}

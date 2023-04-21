package rainier

import (
	"net/http"
)

func (c *RainierController) SubscriptionQuickAdd(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userId := ctx.Value(ContextUserIDKey{}).(int64)
	quickadd := req.URL.Query().Get("quickadd")

	if quickadd == "" {
		w.WriteHeader(400) // TODO: error message
		return
	}

	err := c.Container.Command.SubscribeToFeedByUrl(ctx, quickadd, userId)

	if err != nil {
		c.Container.Render.Text(w, http.StatusInternalServerError, err.Error())
	} else {
		c.Container.Render.Text(w, http.StatusOK, "ok")
	}
}

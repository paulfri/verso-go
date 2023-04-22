package rainier

import (
	"net/http"
)

type SubscriptionQuickAddRequestParams struct {
	Quickadd string `query:"quickadd" validate:"required"`
}

func (c *RainierController) SubscriptionQuickAdd(w http.ResponseWriter, req *http.Request) {
	params := SubscriptionQuickAddRequestParams{}
	err := c.Container.Params(&params, req)

	if err != nil {
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)

	err = c.Container.Command.SubscribeToFeedByURL(ctx, params.Quickadd, userID)

	if err != nil {
		c.Container.Render.Text(w, http.StatusInternalServerError, err.Error())
	} else {
		c.Container.Render.Text(w, http.StatusOK, "ok")
	}
}

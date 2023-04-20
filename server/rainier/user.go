package rainier

import (
	"net/http"
)

type UserResponse struct {
	UserId              string `json:"userId"`
	UserName            string `json:"userName"`
	UserProfileId       string `json:"userProfileId"`
	UserEmail           string `json:"userEmail"`
	IsBloggerUser       bool   `json:"isBloggerUser"`
	SignupTimeSec       int64  `json:"signupTimeSec"`
	IsMultiLoginEnabled bool   `json:"isMultiLoginEnabled"`
	IsPremium           bool   `json:"isPremium"`
}

func (c *RainierController) UserGet(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{})
	user, _ := c.Container.Queries.GetUserById(ctx, userID.(int64))

	c.Container.Render.JSON(w, http.StatusOK, UserResponse{
		UserId:              user.Uuid.String(),
		UserName:            user.Name,
		UserProfileId:       user.Uuid.String(),
		UserEmail:           user.Email,
		IsBloggerUser:       false,
		SignupTimeSec:       user.CreatedAt.Unix(),
		IsMultiLoginEnabled: false,
		IsPremium:           true,
	})
}

func (c *RainierController) UserTokenGet(w http.ResponseWriter, req *http.Request) {
	token := req.Context().Value(ContextAuthTokenKey{}).(string)

	c.Container.Render.Text(w, http.StatusOK, token)
}

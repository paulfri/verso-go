package sierra

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

func (r *SierraRouter) user(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey)
	user, _ := r.Controller.Queries.GetUserById(ctx, userID.(int64))

	r.Controller.Render.JSON(w, http.StatusOK, UserResponse{
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

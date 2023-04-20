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
	SignupTimeSec       int32  `json:"signupTimeSec"`
	IsMultiLoginEnabled bool   `json:"isMultiLoginEnabled"`
	IsPremium           bool   `json:"isPremium"`
}

func (r *SierraRouter) user(w http.ResponseWriter, req *http.Request) {
	r.Controller.Render.JSON(w, http.StatusOK, UserResponse{
		UserId:              "00157a17b192950b65be3791",
		UserName:            "Paul Friedman",
		UserProfileId:       "00157a17b192950b65be3791",
		UserEmail:           "paul@verso.so",
		IsBloggerUser:       false,
		SignupTimeSec:       1370709105,
		IsMultiLoginEnabled: false,
		IsPremium:           true,
	})
}

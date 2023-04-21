package rainier

import (
	"fmt"
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

func (c *RainierController) User(w http.ResponseWriter, req *http.Request) {
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

func (c *RainierController) UserToken(w http.ResponseWriter, req *http.Request) {
	token := req.Context().Value(ContextAuthTokenKey{}).(string)

	c.Container.Render.Text(w, http.StatusOK, token)
}

const UserPreferencesResponse = `{
	"prefs": [{
			"id": "lhn-prefs",
			"value": "{\"subscriptions\":{\"ssa\":\"true\"}}"
	}]
}`

func (c *RainierController) UserPreferences(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(UserPreferencesResponse))
}

const UserStreamPreferencesResponse = `{
    "streamprefs": {}
}`

func (c *RainierController) UserStreamPreferences(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(UserStreamPreferencesResponse))
}

type UserFriend struct {
	P           string    `json:"p"`
	ContactID   string    `json:"contactId"`
	Flags       int32     `json:"flags"`
	Stream      string    `json:"stream"`
	HasShared   bool      `json:"hasSharedItemsOnProfile"`
	ProfileIDs  [1]string `json:"profileIds"`
	UserIDs     [1]string `json:"userIds"`
	Name        string    `json:"givenName"`
	DisplayName string    `json:"displayName"`
	N           string    `json:"n"`
}

func (c *RainierController) UserFriendList(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userId := ctx.Value(ContextUserIDKey{}).(int64)
	strUserId := fmt.Sprintf("%d", userId)
	idSlice := [1]string{strUserId}
	user, _ := c.Container.Queries.GetUserById(ctx, userId)

	c.Container.Render.JSON(w, http.StatusOK, UserFriend{
		ContactID:   strUserId,
		Flags:       1,
		Stream:      "user/" + strUserId + "/state/com.google/broadcast",
		ProfileIDs:  idSlice,
		UserIDs:     idSlice,
		Name:        user.Name,
		DisplayName: user.Name,
	})
}

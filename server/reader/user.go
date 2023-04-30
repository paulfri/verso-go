package reader

import (
	"fmt"
	"net/http"
)

type UserResponse struct {
	UserID              string `json:"userId"`
	UserName            string `json:"userName"`
	UserProfileID       string `json:"userProfileId"`
	UserEmail           string `json:"userEmail"`
	IsBloggerUser       bool   `json:"isBloggerUser"`
	SignupTimeSec       int64  `json:"signupTimeSec"`
	IsMultiLoginEnabled bool   `json:"isMultiLoginEnabled"`
	IsPremium           bool   `json:"isPremium"`
}

func (c *ReaderController) User(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{})
	queries := c.Container.GetQueries(req)
	user, _ := queries.GetUser(ctx, userID.(int64))

	c.Container.Render.JSON(w, http.StatusOK, UserResponse{
		UserID:              user.UUID.String(),
		UserName:            user.Name,
		UserProfileID:       user.UUID.String(),
		UserEmail:           user.Email,
		IsBloggerUser:       false,
		SignupTimeSec:       user.CreatedAt.Unix(),
		IsMultiLoginEnabled: false,
		IsPremium:           true,
	})
}

func (c *ReaderController) UserToken(w http.ResponseWriter, req *http.Request) {
	token := req.Context().Value(ContextAuthTokenKey{}).(string)

	c.Container.Render.Text(w, http.StatusOK, token)
}

const UserPreferencesResponse = `{
	"prefs": [{
			"id": "lhn-prefs",
			"value": "{\"subscriptions\":{\"ssa\":\"true\"}}"
	}]
}`

func (c *ReaderController) UserPreferences(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(UserPreferencesResponse))
}

const UserStreamPreferencesResponse = `{
    "streamprefs": {}
}`

func (c *ReaderController) UserStreamPreferences(w http.ResponseWriter, req *http.Request) {
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

func (c *ReaderController) UserFriendList(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)
	strUserID := fmt.Sprintf("%d", userID)
	idSlice := [1]string{strUserID}
	queries := c.Container.GetQueries(req)
	user, _ := queries.GetUser(ctx, userID)

	c.Container.Render.JSON(w, http.StatusOK, UserFriend{
		ContactID:   strUserID,
		Flags:       1,
		Stream:      "user/" + strUserID + "/state/com.google/broadcast",
		ProfileIDs:  idSlice,
		UserIDs:     idSlice,
		Name:        user.Name,
		DisplayName: user.Name,
	})
}

package reader

import (
	"net/http"

	"github.com/dchest/uniuri"
	"github.com/versolabs/verso/db/query"
	"golang.org/x/crypto/bcrypt"
)

type AuthTokenResponse struct {
	SID  string `json:"SID"`
	LSID string `json:"LSID"`
	Auth string `json:"Auth"`
}

type AuthErrorResponse struct {
	Error string `json:"Error"`
}

type ClientLoginRequest struct {
	Email    string `query:"Email"`
	Password string `query:"Passwd"`
}

func (c *ReaderController) ClientLogin(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	body := ClientLoginRequest{}
	c.Container.BodyParams(&body, req)

	if body.Email == "" || body.Password == "" {
		c.Container.Render.JSON(w, http.StatusBadRequest, AuthErrorResponse{
			Error: "BadAuthentication",
		})

		return
	}

	user, err := c.Container.Queries.GetUserByEmail(ctx, body.Email)
	if err != nil {
		c.Container.Render.JSON(w, http.StatusBadRequest, AuthErrorResponse{
			Error: "BadAuthentication",
		})
	}

	if match(user.Password.String, body.Password) {
		rando := uniuri.NewLen(20)

		// TODO handle error
		c.Container.Queries.CreateReaderToken(
			ctx,
			query.CreateReaderTokenParams{UserID: user.ID, Identifier: rando},
		)

		c.Container.Render.Text(w, http.StatusOK, "SID="+rando+"\nLSID="+rando+"\nAuth="+rando+"\n")
	} else {
		c.Container.Render.JSON(w, http.StatusBadRequest, AuthErrorResponse{
			Error: "BadAuthentication",
		})
	}
}

func match(input string, control string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(input), []byte(control))
	return err == nil
}

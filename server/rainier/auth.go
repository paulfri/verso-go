package rainier

import (
	"net/http"

	"github.com/dchest/uniuri"
	"github.com/versolabs/citra/db/query"
	"golang.org/x/crypto/bcrypt"
)

type AuthTokenResponse struct {
	SID  *string `json:"SID"`
	LSID *string `json:"LSID"`
	Auth *string `json:"Auth"`
}

type AuthErrorResponse struct {
	Error string `json:"Error"`
}

func (c *RainierController) ClientLogin(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	email := req.URL.Query().Get("Email")
	password := req.URL.Query().Get("Passwd")

	if email == "" || password == "" {
		c.Container.Render.JSON(w, http.StatusBadRequest, AuthErrorResponse{
			Error: "BadAuthentication",
		})

		return
	}

	user, err := c.Container.Queries.GetUserByEmail(ctx, email)

	if err != nil {
		c.Container.Render.JSON(w, http.StatusBadRequest, AuthErrorResponse{
			Error: "BadAuthentication",
		})
	}

	if match(user.Password.String, password) {
		rando := uniuri.NewLen(20)

		// TODO handle error
		c.Container.Queries.CreateReaderToken(
			ctx,
			query.CreateReaderTokenParams{UserID: user.ID, Identifier: rando},
		)

		c.Container.Render.JSON(w, http.StatusOK, AuthTokenResponse{
			Auth: &rando,
		})
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

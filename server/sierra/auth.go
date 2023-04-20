package sierra

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

func (r *SierraRouter) token(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	email := req.URL.Query().Get("Email")
	password := req.URL.Query().Get("Passwd")

	if email == "" || password == "" {
		r.Controller.Render.JSON(w, http.StatusBadRequest, AuthErrorResponse{
			Error: "BadAuthentication",
		})

		return
	}

	user, err := r.Controller.Queries.GetUserByEmail(ctx, email)

	if err != nil {
		r.Controller.Render.JSON(w, http.StatusBadRequest, AuthErrorResponse{
			Error: "BadAuthentication",
		})
	}

	if match(user.Password.String, password) {
		rando := uniuri.NewLen(20)

		// TODO handle error
		r.Controller.Queries.CreateToken(ctx, query.CreateTokenParams{UserID: user.ID, Identifier: rando})

		r.Controller.Render.JSON(w, http.StatusOK, AuthTokenResponse{
			Auth: &rando,
		})
	} else {
		r.Controller.Render.JSON(w, http.StatusBadRequest, AuthErrorResponse{
			Error: "BadAuthentication",
		})
	}
}

func match(input string, control string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(input), []byte(control))
	return err == nil
}

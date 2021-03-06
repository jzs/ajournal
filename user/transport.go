package user

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jzs/ajournal/utils"
	"github.com/jzs/ajournal/utils/logger"

	"github.com/gorilla/mux"
)

type userCtx string

const cookieName = "a"

var userContext userCtx

func init() {
	userContext = userCtx("usercontext")
}

// TestContextWithUser sets up a test context with a given user in for testing.
func TestContextWithUser(u *User) context.Context {
	return context.WithValue(context.Background(), userContext, u)
}

// ApplyUserToRequest Applys current user to Request based on current user cookie
func ApplyUserToRequest(r *http.Request, us Service) *http.Request {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return r
	}
	user, err := us.UserWithToken(r.Context(), cookie.Value)
	if err != nil {
		return r
	}

	ctx := context.WithValue(r.Context(), userContext, user)
	nr := r.WithContext(ctx)
	return nr
}

// FromContext gets the currently logged in user from  a given request. Returns nil if no user
func FromContext(ctx context.Context) *User {
	val := ctx.Value(userContext)
	if u, ok := val.(*User); ok {
		return u
	}
	return nil
}

// SetupHandler sets up the handler routes for the user service
func SetupHandler(r *mux.Router, us Service, l logger.Logger) {
	r.Path("/users/me").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usr, err := us.Me(r.Context())
		utils.JSONResp(r.Context(), l, r, w, usr, err)
	})

	r.Path("/users/{username}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]
		user, err := us.User(r.Context(), username)
		utils.JSONResp(r.Context(), l, r, w, user, err)
	})

	// Create user
	r.Path("/users").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := &User{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(u)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, utils.NewAPIError(err, http.StatusBadRequest, "Not valid json"))
			return
		}
		err = us.Register(r.Context(), u)
		utils.JSONResp(r.Context(), l, r, w, nil, err)
	})

	// Log in
	r.Path("/users/login").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := &User{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(u)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, utils.NewAPIError(err, http.StatusBadRequest, "Not valid json"))
		}

		token, err := us.Login(r.Context(), u.Username, u.Password)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, token, err)
			return
		}

		http.SetCookie(w, CreateCookie(token))

		utils.JSONResp(r.Context(), l, r, w, token, err)
	})

	// Log out
	r.Path("/users/logout").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			//TODO: Try to get the cookie from elsewhere
			utils.JSONResp(r.Context(), l, r, w, nil, nil) // nil, nil since we ignore the error and just "log out" the user.
			return
		}
		if err := us.Logout(r.Context(), cookie.Value); err != nil {
			l.Errorf(r.Context(), "Failed to log out, %+v", err.Error())
		}

		cookie.Expires = time.Now().Add(-24 * 7 * time.Hour) // Invalidate cookie by set time to zero time
		cookie.MaxAge = -1
		cookie.Path = "/"
		cookie.Name = cookieName
		cookie.HttpOnly = true
		http.SetCookie(w, cookie)
		utils.JSONResp(r.Context(), l, r, w, nil, nil)
	})
}

// CreateCookie creates an authentication cookie
func CreateCookie(token *Token) *http.Cookie {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    token.Token,
		Path:     "/",
		Expires:  token.Expires,
		HttpOnly: true,
		//Secure: true,
	}
	return cookie
}

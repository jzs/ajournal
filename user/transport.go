package user

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const cookieName = "a"
const userContext = "usercontext"

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
func SetupHandler(r *mux.Router, us Service) {

	// Create user
	r.Path("/users").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Parse User from request!
		u := &User{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(u)
		if err != nil {
			// bad args...
			panic("bad args")
		}
		err = us.Register(r.Context(), u)
		err = JSONResp(w, nil, err)
		if err != nil {
			// TODO Log this error or panic
		}
	})

	// Log in
	r.Path("/users/login").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := &User{}
		dec := json.NewDecoder(r.Body)
		dec.Decode(u)

		token, err := us.Login(r.Context(), u.Username, u.Password)
		if err != nil {
			JSONResp(w, token, err)
			return
		}

		cookie := &http.Cookie{
			Name:     cookieName,
			Value:    token.Token,
			Path:     "/",
			Expires:  token.Expires,
			HttpOnly: true,
			//Secure: true,
		}
		http.SetCookie(w, cookie)

		err = JSONResp(w, token, err)
		if err != nil {
			// TODO: Handle Error.
		}
	})

	// Log out
	r.Path("/users/logout").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			//TODO: Try to get the cookie from elsewhere
			err = JSONResp(w, nil, nil) // nil, nil since we ignore the error and just "log out" the user.
			return
		}
		us.Logout(r.Context(), cookie.Value)
		cookie.Expires = time.Now().Add(-24 * 7 * time.Hour) // Invalidate cookie by set time to zero time
		cookie.MaxAge = -1
		cookie.Path = "/"
		cookie.Name = cookieName
		cookie.HttpOnly = true
		http.SetCookie(w, cookie)
		JSONResp(w, nil, nil)
	})

}

// Json response handling and error handling!

type jsonresp struct {
	Data   interface{}
	Status int64
	Error  string
}

// JSONResp writes a json response to the responsewriter
func JSONResp(w http.ResponseWriter, data interface{}, err error) error {
	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := jsonresp{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}
		err = enc.Encode(resp)
		if err != nil {
			// Log this error or panic!
			return err
		}
		return nil
	}

	resp := jsonresp{
		Data:   data,
		Status: http.StatusOK,
	}
	enc.Encode(resp)
	return nil
}

package routes

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	"github.com/Kzynyi/ass/internal/middleware"
	"github.com/Kzynyi/ass/internal/service/user"
)

func GetNewMux() http.Handler {

	router := httprouter.New()
	router.POST("/api/v1/user/signup", signUpUser)
	router.POST("/api/v1/user/signin", signInUser)
	router.GET("/api/v1/lobbies/all", middleware.Adapt(viewLobbies, middleware.ValidateJWT()))
	return cors.Default().Handler(router)

}

func signUpUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u := new(user.User)
	json.NewDecoder(r.Body).Decode(u)
	_, err := user.CreateUser(r.Context(), u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		m := make(map[string]string)
		m["status"] = err.Error()
		json.NewEncoder(w).Encode(m)
	} else {
		w.WriteHeader(http.StatusOK)
	}

}

func signInUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	m := make(map[string]string)
	json.NewDecoder(r.Body).Decode(&m)
	token, err := user.SignInUser(r.Context(), m)
	response := make(map[string]string)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		response["message"] = err.Error()
	} else {

		cookie := new(http.Cookie)
		cookie.Name = "token"
		cookie.HttpOnly = true
		cookie.Value = token
		http.SetCookie(w, cookie)
	}
	json.NewEncoder(w).Encode(response)
}

func viewLobbies(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	m := make(map[string]string)
	m["lobbies"] = "all lobbies"
	json.NewEncoder(w).Encode(m)
}

package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

type Adapter func(httprouter.Handle) httprouter.Handle

func Adapt(h httprouter.Handle, adapters ...Adapter) httprouter.Handle {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func ValidateJWT() Adapter {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			tk, err := r.Cookie("token")
			if err != nil {
				panic(err)
			}
			valid := ValidateToken(tk.Value)
			if !valid {
				m := make(map[string]string)
				m["status"] = "Invalid token"
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(m)
			} else {
				next(w, r, p)
			}

		}
	}
}
func ValidateToken(tokenString string) bool {

	tk, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s", "Wrong format")
		}
		return []byte("secretkey"), nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if claims, ok := tk.Claims.(jwt.RegisteredClaims); ok {
		fmt.Println(claims.ExpiresAt)
	}
	return true
}

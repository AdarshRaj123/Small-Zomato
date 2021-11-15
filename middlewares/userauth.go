package middlewares

import (
	"SmallZomato/database/dbhelper"
	"SmallZomato/models"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

type ContextKeys string

const (
	userContext ContextKeys = "__userContext"
)

/*func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")
		user, err := dbhelper.GetUserBySession(apiKey)
		if err != nil || user == nil {

			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userContext, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}*/

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apikey:=r.Header.Get("x-api-key")
		if apikey =="" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var mySigningKey = []byte("123")

		_, err := jwt.Parse(apikey, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySigningKey, nil
		})
		if err!=nil{
			fmt.Println(err)
			w.WriteHeader(http.StatusForbidden)
		}


		next.ServeHTTP(w, r)
	})
}


func UserContext(r *http.Request) *models.User {
	if user, ok := r.Context().Value(userContext).(*models.User); ok && user != nil {
		return user
	}
	return nil
}


func AdminCheck() func(http.Handler)http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apikey := r.Header.Get("x-api-key")
			if apikey == "" {

				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			var mySigningKey = []byte("123")

			cs, err := jwt.Parse(apikey, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error in parsing")
				}
				return mySigningKey, nil
			})
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			claims := cs.Claims.(jwt.MapClaims)
			if claims["role"].(string) == "admin" {
				user,_ := dbhelper.GetUserDetail(claims["email"].(string))
				ctx := context.WithValue(r.Context(), userContext, user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		},
		)
	}
}
func UserCheck() func(http.Handler)http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apikey := r.Header.Get("x-api-key")
			if apikey == "" {

				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			var mySigningKey = []byte("123")

			cs, err := jwt.Parse(apikey, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error in parsing")
				}
				return mySigningKey, nil
			})
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			claims := cs.Claims.(jwt.MapClaims)
			if claims["role"].(string) == "user" {
				user,_ := dbhelper.GetUserDetail(claims["email"].(string))
				ctx := context.WithValue(r.Context(), userContext, user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		},
		)
	}
}
func SubAdminCheck()func(http.Handler)http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apikey := r.Header.Get("x-api-key")
			if apikey == "" {

				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			var mySigningKey = []byte("123")

			cs, err := jwt.Parse(apikey, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error in parsing")
				}
				return mySigningKey, nil
			})
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			claims := cs.Claims.(jwt.MapClaims)
			if claims["role"].(string) == "subadmin" {
				user,_ := dbhelper.GetUserDetail(claims["email"].(string))
				ctx := context.WithValue(r.Context(), userContext, user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		},
		)
	}
}



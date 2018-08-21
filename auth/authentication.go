package auth

import (
	"strings"
	"encoding/base64"
	"net/http"
	"github.com/jeevanrd/sms-service/database"
	"context"
	"encoding/json"
)

type  Response struct {
	Message		string `json:"message"`
	Error		string `json:"error"`
}

type  AuthResponse struct {
	Account		database.Account
	Status		bool
	Message		string
}

type Auth struct {
	Repo database.Repository
}

func (a *Auth) ContentTypeHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if (contentType == "") {
			buildResponse(w, "unprocessable entity", http.StatusUnprocessableEntity)
			return
		}
		if(strings.ToLower(contentType) != "application/json") {
			buildResponse(w, "unprocessable entity", http.StatusUnprocessableEntity)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (a *Auth) AuthHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			buildResponse(w, "Please pass valid credentials", http.StatusForbidden)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			buildResponse(w, "Please pass valid credentials", http.StatusForbidden)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			buildResponse(w, "Please pass valid credentials", http.StatusForbidden)
			return
		}


		account, err := a.Repo.GetAccount(pair[0], pair[1])
		if(err != nil) {
			buildResponse(w, "Please pass valid credentials", http.StatusForbidden)
			return
		}

		var ctx context.Context
		ctx = context.WithValue(ctx, "url-parameters", r.URL.Query())
		ctx = context.WithValue(ctx, "request-method", r.Method)
		ctx = context.WithValue(ctx, "account", account)

		reqWithContext := r.WithContext(ctx)
		h.ServeHTTP(w, reqWithContext)
	})
}

func buildResponse(w http.ResponseWriter, err string, code int) {
	response, _ := json.Marshal(&Response{"", err})
	w.WriteHeader(code);
	w.Write(response)
}

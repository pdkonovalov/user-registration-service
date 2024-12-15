package v1

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/pdkonovalov/user-registration-service/pkg/config"
	"github.com/pdkonovalov/user-registration-service/pkg/email"
	"github.com/pdkonovalov/user-registration-service/pkg/email/templates"
	"github.com/pdkonovalov/user-registration-service/pkg/jwt"
	"github.com/pdkonovalov/user-registration-service/pkg/storage"
)

func HandleEmailVerify(config *config.Config, storage storage.Storage, email *email.Email) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rEmail := r.URL.Query().Get("email")
		_, err := mail.ParseAddress(rEmail)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rEmail = strings.ToLower(rEmail)
		rand, _ := rand.Int(rand.Reader, big.NewInt(90000))
		code := 10000 + rand.Int64()
		msg := templates.VerificationCodeMsg(rEmail, int(code))
		err = email.Send(rEmail, msg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = storage.DeleteEmailCode(rEmail)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = storage.WriteEmailCode(rEmail, int(code), time.Now().Add(config.EmailCodeTtl))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func HandleNewUser(storage storage.Storage) http.HandlerFunc {
	type request struct {
		Name      string `json:"Name"`
		Username  string `json:"Username"`
		Password  string `json:"Password"`
		Email     string `json:"Email"`
		EmailCode int    `json:"EmailCode"`
	}
	isValidRequest := func(req *request) bool {
		if req.Name == "" || req.Username == "" || len(req.Password) < 5 || len(req.Password) > 20 {
			return false
		}
		_, err := mail.ParseAddress(req.Email)
		return err != nil
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !isValidRequest(req) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		req.Email = strings.ToLower(req.Email)
		code, exp_time, exist, err := storage.FindEmailCode(req.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !exist ||
			code != req.EmailCode ||
			time.Now().Compare(exp_time) != -1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = storage.DeleteEmailCode(req.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _, _, exist, err = storage.FindUser(req.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if exist {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		password_hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = storage.WriteNewUser(req.Email, req.Name, req.Username, string(password_hash))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func HandleNewPassword(storage storage.Storage) http.HandlerFunc {
	type request struct {
		Password  string `json:"Password"`
		Email     string `json:"Email"`
		EmailCode int    `json:"EmailCode"`
	}
	isValidRequest := func(req *request) bool {
		if len(req.Password) < 5 || len(req.Password) > 20 {
			return false
		}
		_, err := mail.ParseAddress(req.Email)
		return err != nil
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !isValidRequest(req) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		req.Email = strings.ToLower(req.Email)
		code, exp_time, exist, err := storage.FindEmailCode(req.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !exist ||
			code != req.EmailCode ||
			time.Now().Compare(exp_time) != -1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = storage.DeleteEmailCode(req.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _, _, exist, err = storage.FindUser(req.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !exist {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		password_hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = storage.UpdatePassword(req.Email, string(password_hash))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func HandleNewJwt(storage storage.Storage, jwt *jwt.JwtGenerator) http.HandlerFunc {
	type response struct {
		AccessToken  string `json:"AccessToken"`
		RefreshToken string `json:"RefreshToken"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		rEmail := r.URL.Query().Get("email")
		_, err := mail.ParseAddress(rEmail)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rEmail = strings.ToLower(rEmail)
		rPassword := r.URL.Query().Get("password")
		_, _, password_hash, exist, err := storage.FindUser(rEmail)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !exist {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(rPassword))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		accessToken, err := jwt.GenerateAccessToken(rEmail, r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		refreshToken, err := jwt.GenerateRefreshToken(rEmail)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp, err := json.Marshal(&response{accessToken, refreshToken})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func HandleRefreshJwt(storage storage.Storage, email *email.Email, jwt *jwt.JwtGenerator) http.HandlerFunc {
	type request struct {
		AccessToken  string `json:"AccessToken"`
		RefreshToken string `json:"RefreshToken"`
	}
	type response struct {
		AccessToken  string `json:"AccessToken"`
		RefreshToken string `json:"RefreshToken"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		aEmail, ip, _ := jwt.ValidateAccessToken(req.AccessToken)
		rEmail, valid := jwt.ValidateRefreshToken(req.RefreshToken)
		if !valid {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if aEmail != rEmail {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if ip != r.RemoteAddr {
			msg := templates.ChangeIpAllertMsg(rEmail)
			email.Send(rEmail, msg)
		}
		accessToken, err := jwt.GenerateAccessToken(rEmail, r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		refreshToken, err := jwt.GenerateRefreshToken(rEmail)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp, err := json.Marshal(&response{accessToken, refreshToken})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

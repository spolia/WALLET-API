package internal

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/securecookie"
	"github.com/spolia/wallet-api/internal/wallet/movement"
	"github.com/spolia/wallet-api/internal/wallet/user"
)

var validate = validator.New()

var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func login(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginRequest struct {
			Alias    string `json:"alias" validate:"required"`
			Password string `json:"password" validate:"required"`
		}

		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validate.Struct(loginRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ok, err := service.ValidateCredential(r.Context(), strings.ToLower(loginRequest.Alias), loginRequest.Password)
		if err != nil {
			if err == user.ErrorInvalidCredential {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if ok {
			setSession(loginRequest.Alias, w)
		}

		w.Write([]byte("Logged in"))
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	w.Write([]byte("Logged out"))
}

//create user
func createUser(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userRequest user.User
		if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validate.Struct(userRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		userRequest.Alias = strings.ToLower(userRequest.Alias)
		err := service.CreateUser(r.Context(), userRequest)
		if err != nil {
			if err == user.ErrorAlreadyExist {
				http.Error(w, "alias or email already exist", http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("ok")
		return
	}
}

func getBalance(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := getUserAlias(r)
		if alias == "" {
			http.Error(w, "log in is required", http.StatusUnauthorized)
			return
		}

		userBalance, err := service.GetBalance(r.Context(), strings.ToLower(alias))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(userBalance)
		return
	}
}

func getHistory(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := getUserAlias(r)
		if alias == "" {
			http.Error(w, "log in is required", http.StatusUnauthorized)
			return
		}

		history, err := service.GetHistory(r.Context(), strings.ToLower(alias))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(history)
		return
	}
}

func send(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := getUserAlias(r)
		if alias == "" {
			http.Error(w, "log in is required", http.StatusUnauthorized)
			return
		}

		var sendRequest struct {
			Amount           float64 `json:"amount" validate:"required,gt=0"`
			CurrencyName     string  `json:"currencyname" validate:"required,oneof=usdt btc ars"`
			InteractionAlias string  `json:"interactionalias" validate:"required"`
		}

		if err := json.NewDecoder(r.Body).Decode(&sendRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validate.Struct(sendRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var m movement.Movement
		m.Alias = strings.ToLower(alias)
		m.InteractionAlias = strings.ToLower(sendRequest.InteractionAlias)
		m.CurrencyName = strings.ToUpper(sendRequest.CurrencyName)
		m.Type = movement.SendMov
		m.Amount = sendRequest.Amount

		if m.Alias == m.InteractionAlias {
			http.Error(w, "the destiny and origin alias have to be different", http.StatusBadRequest)
			return
		}

		err := service.Send(r.Context(), m)
		if err != nil {
			if err == user.ErrorDestinyUserNotFound {
				http.Error(w, "wrong destiny alias", http.StatusBadRequest)
				return
			}

			if err == movement.ErrorWrongCurrency || err == movement.ErrorInsufficientFunds {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("ok")
		return
	}
}

func deposit(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := getUserAlias(r)
		if alias == "" {
			http.Error(w, "log in is required", http.StatusUnauthorized)
			return
		}

		var depositRequest struct {
			Amount       float64 `json:"amount" validate:"required,gt=0"`
			CurrencyName string  `json:"currencyname" validate:"required,oneof=usdt btc ars"`
		}

		if err := json.NewDecoder(r.Body).Decode(&depositRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validate.Struct(depositRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var m movement.Movement
		m.Alias = strings.ToLower(alias)
		m.CurrencyName = strings.ToUpper(depositRequest.CurrencyName)
		m.InteractionAlias = strings.ToLower(alias)
		m.Type = movement.DepositMov
		m.Amount = depositRequest.Amount

		err := service.AutoDeposit(r.Context(), m)
		if err != nil {
			if err == movement.ErrorWrongCurrency {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("ok")
		return
	}
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func setSession(alias string, w http.ResponseWriter) {
	value := map[string]string{
		"alias": alias,
	}

	encoded, err := cookieHandler.Encode("session", value)
	if err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func getUserAlias(request *http.Request) (alias string) {
	cookie, err := request.Cookie("session")
	if err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			alias = cookieValue["alias"]
		}
	}
	return alias
}

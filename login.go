package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	queryUsers        = "select id, login, pass from hbpro_user.account ha where trash = false and login = $1"
	queryGlobalParams = "select currency_rate, loyalty_coefficient from public.global_parameter"
)

func login(w http.ResponseWriter, r *http.Request) {
	// Set content type in header
	w.Header().Set("Content-Type", `application/json`)

	// Get login data from request
	errorData, loginData := getLoginData(r)
	if nil != errorData {
		writeError(w, errorToJson(errorData.Error), errorData.Code)
		return
	}

	// Get account from DB by login data
	errorData, account := getAccount(loginData)
	if nil != errorData {
		writeError(w, errorToJson(errorData.Error), errorData.Code)
		return
	}

	// Compare bcrypt password
	checkPasswordHash := bcrypt.CompareHashAndPassword(account.Pass, []byte(loginData.Password))

	// If password is valid then generate token, put token to the storage, return token
	if checkPasswordHash == nil {

		// Get global params and add it to token
		globalParams := getGlobalParameters()

		// Generate token
		errorData, tokenString := generateToken(getCustomClaims(*account, globalParams))
		if nil != errorData {
			writeError(w, errorToJson(errorData.Error), errorData.Code)
			return
		}

		// Save token to storage
		errorData = saveTokenToStorage(tokenString)
		if nil != errorData {
			writeError(w, errorToJson(errorData.Error), errorData.Code)
			return
		}

		// Generate response
		response, err := json.Marshal(Response{tokenString})
		if err != nil {
			writeError(w, errorToJson(err.Error()), http.StatusUnauthorized)
			return
		}
		w.Write(response)
	} else {
		writeError(w, errorToJson("Invalid login or password"), http.StatusUnauthorized)
	}
}

func getLoginData(r *http.Request) (*ErrorData, *LoginData) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &ErrorData{err.Error(), http.StatusBadRequest}, nil
	}
	loginData := new(LoginData)
	err = json.Unmarshal(body, &loginData)
	if err != nil {
		return &ErrorData{err.Error(), http.StatusBadRequest}, nil
	}
	return nil, loginData
}

func getAccount(loginData *LoginData) (*ErrorData, *AccountData) {
	row := db.QueryRow(
		queryUsers,
		loginData.Login,
	)
	account := new(AccountData)

	err := row.Scan(&account.Id, &account.Login, &account.Pass)
	if err == sql.ErrNoRows {
		return &ErrorData{"User not found", http.StatusNotFound}, nil
	} else if err != nil {
		return &ErrorData{
			"Something wrong on connecting to database: " + err.Error(),
			http.StatusInternalServerError,
		}, nil
	}
	return nil, account
}

func getGlobalParameters() GlobalParameter {
	// Select global parameters
	globalParams := GlobalParameter{}
	row := db.QueryRow(queryGlobalParams)
	row.Scan(&globalParams.CurrencyRate, &globalParams.LoyaltyCoefficient)
	return globalParams
}

func generateToken(customClaims CustomClaims) (*ErrorData, string) {
	claims := JwtData{
		StandardClaims: jwt.StandardClaims{
			// The token will expire after 30 days.
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		},

		// Store userId, rights etc.
		CustomClaims: customClaims,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(signKey)

	if err != nil {
		return &ErrorData{err.Error(), http.StatusUnauthorized}, ""
	}

	return nil, tokenString
}

func saveTokenToStorage(tokenString string) *ErrorData {
	respStorage, err := requestToStorage(http.MethodPut, settings.TokenStorage+"/token", tokenString)

	if err != nil {
		return &ErrorData{err.Error(), http.StatusInternalServerError}
	}

	if respStorage.Status != http.StatusOK {
		return &ErrorData{respStorage.Body, respStorage.Status}
	}

	if err != nil {
		return &ErrorData{err.Error(), http.StatusInternalServerError}
	}
	return nil
}

func getMd5HashString(password string) string {
	md5Object := md5.New()
	io.WriteString(md5Object, password)
	md5Hash := md5Object.Sum(nil)
	return hex.EncodeToString(md5Hash)
}

func getCustomClaims(account AccountData, globalParams GlobalParameter) CustomClaims {
	customClaims := CustomClaims{}
	customClaims.UserId = account.Id
	if 0 != globalParams.CurrencyRate {
		customClaims.CurrencyRate = globalParams.CurrencyRate
	}
	if 0 != globalParams.LoyaltyCoefficient {
		customClaims.LoyaltyCoefficient = globalParams.LoyaltyCoefficient
	}
	return customClaims
}

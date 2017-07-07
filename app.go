package main

import (
	model "Flexer/Model"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//VARS
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

const (
	secretKey = "tes"
)

var ScreenshotStorage string

//CONNECTION
func (a *App) Initialize(user, password, dbname, screenshotStorage string) {
	connectionString :=
		fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", user, password, dbname)
	ScreenshotStorage = screenshotStorage
	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

//RUN
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

//ROUTES
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/login", a.Login).Methods("POST")
	a.Router.Handle("/addActivity", jwtMiddleware.Handler(http.HandlerFunc(a.AddActivity))).Methods("POST")
	a.Router.Handle("/addActivity/screenshot", jwtMiddleware.Handler(http.HandlerFunc(a.AddActivityScreenshot))).Methods("POST")
}

//HANDLERS
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	var loginX model.Login
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&loginX); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	login := model.Login{Userlogin: loginX.Userlogin, Password: loginX.Password}
	if err := login.DoLogin(a.DB); err != nil {
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "Wrong username or password")
		}
		return
	}

	//Getting Client Info
	client := model.Client{ClientID: login.ClientID}
	if err := client.GetClient(a.DB); err != nil {
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "No Client Found")
		}
		return
	}
	token := GetToken(login.Username)
	result := map[string]string{"token": token, "clientname": client.ClientName, "username": login.Username}
	respondWithJSON(w, http.StatusOK, result)
}

/* AddActivity :
Params:
- TransactionID (GUID)
- Userlogin (string)
- ApplicationName (string)
- URL (string)
- Mouseclick (int)
- Keystroke (int)
*/
func (a *App) AddActivity(w http.ResponseWriter, r *http.Request) {
	var transactions []model.Transaction

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&transactions); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	var errorCount int
	var clientTrID []string
	for _, transaction := range transactions {
		application := model.Application{ApplicationName: transaction.ApplicationName}
		status, err := application.CheckApplicationExist(a.DB)
		if status == 0 {
			if err == sql.ErrNoRows {
				application.CreateApplication(a.DB)
			} else {
				clientTrID = append(clientTrID, strconv.Itoa(transaction.TransactionID))
				errorCount++
				continue
				//respondWithError(w, http.StatusInternalServerError, "Error Checking Application")
			}
		}
		//Getting ApplicationID
		application.GetApplicationID(a.DB)
		transaction.ApplicationID = application.ApplicationID
		//Getting UserID
		user := model.User{UserLogin: transaction.Userlogin}
		err = user.GetUserID(a.DB)
		if err != nil {
			clientTrID = append(clientTrID, strconv.Itoa(transaction.TransactionID))
			errorCount++
			continue
		}
		transaction.UserID = user.UserID
		//Create Transaction
		rows, err := transaction.CreateTransaction(a.DB)
		if err != nil {
			clientTrID = append(clientTrID, strconv.Itoa(transaction.TransactionID))
			errorCount++
			continue
			//respondWithError(w, http.StatusInternalServerError, "Failed Create Transaction")
		}
		row, err := rows.RowsAffected()
		if err != nil {
			clientTrID = append(clientTrID, strconv.Itoa(transaction.TransactionID))
			errorCount++
			continue
			//respondWithError(w, http.StatusInternalServerError, "Failed Get Rows Affected")
		} else if row == 0 {
			clientTrID = append(clientTrID, strconv.Itoa(transaction.TransactionID))
			errorCount++
			continue
			//respondWithError(w, http.StatusInternalServerError, "No Transaction Inserted, ApplicationName : "+transaction.ApplicationName+", URL: "+transaction.URL)
		}
	}
	encjson, _ := json.Marshal(clientTrID)
	var returnMsg string
	if errorCount == 0 {
		returnMsg = "All Transaction Successfully Inserted"
	} else {
		returnMsg = "Error when Inserting"
	}
	result := map[string]string{"errorTrID": string(encjson), "message": returnMsg}
	respondWithJSON(w, http.StatusOK, result)
}

/* AddActivityScreenshot :
header : multipart/form-data
Params :
- screenshot : {file}
- date : "{date}"
*/
func (a *App) AddActivityScreenshot(w http.ResponseWriter, r *http.Request) {
	//var transactions []model.Transaction
	var size int64 = 2 * 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, size)
	err := r.ParseMultipartForm(size)
	if err != nil {
		return
	}

	formdata := r.MultipartForm
	for _, files := range formdata.File {
		for i := range files {
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
			}

			fp := filepath.Join(ScreenshotStorage, files[i].Filename)
			out, err := os.Create(fp)
			defer out.Close()
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
			}

			_, err = io.Copy(out, file)

			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
			}
		}
	}
	result := map[string]interface{}{
		"status":  "OK",
		"message": "Successfully insert screenshot",
	}
	respondWithJSON(w, http.StatusOK, result)
}

//TOKEN
func GetToken(Username string) string {

	/* Create the token */
	token := jwt.New(jwt.SigningMethodHS256)

	/* Create a map to store our claims*/
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims["name"] = Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	/* Sign the token with our secret */
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		panic(err)
	}
	/* Finally, write the token to the browser window */
	return tokenString
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

//HELPER
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

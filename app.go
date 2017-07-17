package main

import (
	model "FlexerAPI/Model"
	query "FlexerAPI/Query"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

//VARS
type App struct {
	Router *mux.Router
	DB     *sql.DB
	Config model.Config
}

const (
	secretKey = "tes"
)

//CONNECTION
func (a *App) Initialize() { //user, password, host, port, dbname, screenshotStorage string) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		a.Config.Database.User,
		a.Config.Database.Password,
		a.Config.Database.Host,
		a.Config.Database.Port,
		a.Config.Database.DBName)
	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err.Error())
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
	a.Router.Handle("/logout", jwtMiddleware.Handler(http.HandlerFunc(a.Logout))).Methods("POST")
	a.Router.Handle("/addActivity", jwtMiddleware.Handler(http.HandlerFunc(a.AddActivity))).Methods("POST")
	a.Router.Handle("/addActivity/screenshot", jwtMiddleware.Handler(http.HandlerFunc(a.AddActivityScreenshot))).Methods("POST")

	a.Router.HandleFunc("/cms/login", a.CMSLogin).Methods("POST")
	a.Router.Handle("/cms/addEmployee", jwtMiddleware.Handler(http.HandlerFunc(a.AddEmployee))).Methods("POST")
	a.Router.Handle("/cms/editEmployee", jwtMiddleware.Handler(http.HandlerFunc(a.EditEmployee))).Methods("POST")
}

//HANDLERS

/* Login :
- Email: string
- Password : string
- LocationType : string
- IPAddress : string
- Lat : string
- Long : string
*/
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	fmt.Println("Login Called")
	var loginX model.Login
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&loginX); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}
	login := model.Login{
		Email:        loginX.Email,
		Password:     loginX.Password,
		LocationType: loginX.LocationType,
		IPAddress:    loginX.IPAddress,
		City:         loginX.City,
		Lat:          loginX.Lat,
		Long:         loginX.Long,
	}
	if err := login.DoLogin(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
	}

	if login.ResultCode == 1 {
		token := GetToken(strconv.Itoa(int(login.Session.SessionID)))
		result := map[string]interface{}{"status": login.ResultCode, "description": login.ResultDescription, "token": token, "session_id": login.Session.SessionID}
		respondWithJSON(w, http.StatusOK, result)
	} else {
		respondWithError(w, http.StatusInternalServerError, login.ResultDescription, login.ResultCode)
	}
}

/* Logout :
- SessionID: int
*/
func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	var Logout model.Logout
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&Logout); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}
	logout := model.Logout{
		SessionID: Logout.SessionID,
	}
	if err := logout.DoLogout(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
	}

	if logout.ResultCode == 1 {
		result := map[string]interface{}{"status": logout.ResultCode, "description": logout.ResultDescription}
		respondWithJSON(w, http.StatusOK, result)
	} else {
		respondWithError(w, http.StatusInternalServerError, logout.ResultDescription, logout.ResultCode)
	}
}

/* AddActivity :
Params:
- TransactionID (GUID)
- SessionID (int)
- ActivityName (string)
- ActivityType (string)
- Mouseclick (int)
- Keystroke (int)
- StartDate (string)
- EndDate (string)
*/
func (a *App) AddActivity(w http.ResponseWriter, r *http.Request) {
	var transactions []model.Transaction
	var Session model.Session
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&Session); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload : Session", -2)
		return
	}
	transactions = Session.Transactions
	if err := Session.FrontCheckSession(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
	}
	//Preparing
	tx, err := a.DB.Begin()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
	}
	defer tx.Rollback()
	stmt, err := a.DB.Prepare(query.SearchQuery("createTransactionQuery"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
	}
	defer stmt.Close()
	for _, transaction := range transactions {
		//Convert Date
		transaction.StartDate = SyncDate(transaction.StartDate, Session.ClientDate, Session.ServerDate)
		transaction.EndDate = SyncDate(transaction.EndDate, Session.ClientDate, Session.ServerDate)
		if !checkValidityPeriod(transaction.StartDate, Session.StartTime, Session.EndTime) {
			continue
		}
		fmt.Println(transaction.StartDate)
		err := stmt.QueryRow(
			Session.SessionID,
			transaction.ActivityName,
			transaction.ActivityType,
			transaction.Keystroke,
			transaction.Mouseclick,
			transaction.StartDate,
			transaction.EndDate,
		).Scan(&transaction.ResultCode, &transaction.ResultDescription)
		if err != nil {
			tx.Rollback()
			respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		}
		if transaction.ResultCode != 1 {
			tx.Rollback()
			respondWithError(w, http.StatusInternalServerError, transaction.ResultDescription, transaction.ResultCode)
		}
	}
	tx.Commit()
	result := map[string]interface{}{"status": 1, "description": "All Transaction Successfully Inserted"}
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
	file, _, err := r.FormFile("img")
	if err != nil {
		return
	}
	var Session model.Session
	var Screenshot model.Screenshot
	Session.SessionID, _ = strconv.ParseInt(r.FormValue("sessionID"), 10, 64)
	Screenshot.SessionID = Session.SessionID
	Screenshot.ScreenshotDate = SyncDate(r.FormValue("screenshotDate"), Session.ClientDate, Session.ServerDate)
	Screenshot.ActivityName = r.FormValue("activityName")
	Screenshot.ActivityType = r.FormValue("activityType")

	if err := Session.FrontCheckSession(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
	}
	if !checkValidityPeriod(Screenshot.ScreenshotDate, Session.StartTime, Session.EndTime) {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload : ScreenshotDate outside Session Date", -2)
		return
	}

	//Getting Parameter
	if err := Screenshot.GetScreenshotParam(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
	}
	if Screenshot.ResultCode != 1 {
		respondWithError(w, http.StatusInternalServerError, Screenshot.ResultDescription, Screenshot.ResultCode)
	}

	//Start Uploading File
	Screenshot.ResultDescription, Screenshot.ResultCode = a.UploadToS3(file, Screenshot.Filename)

	//Reporting Status
	if err := Screenshot.ReportScreenshotStatus(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
	}
	if Screenshot.ResultCode != 1 {
		respondWithError(w, http.StatusInternalServerError, Screenshot.ResultDescription, -1)
	}

	result := map[string]interface{}{
		"status":  1,
		"message": "Successfully insert screenshot",
	}
	respondWithJSON(w, http.StatusOK, result)
}

//CMS API
func (a *App) CMSLogin(w http.ResponseWriter, r *http.Request) {
	var loginX model.Login
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&loginX); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}
	login := model.Login{Email: loginX.Email, Password: loginX.Password}
	if err := login.DoLoginCMS(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
	}

	token := GetToken(strconv.Itoa(int(login.Session.SessionID)))
	result := map[string]interface{}{"token": token, "client_id": login.ClientID}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *App) AddEmployee(w http.ResponseWriter, r *http.Request) {
	var User model.User
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&User); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}
	user := model.User{
		ClientID:     User.ClientID,
		UserName:     User.UserName,
		Role:         User.Role,
		SuperiorID:   User.SuperiorID,
		Email:        User.Email,
		UserPassword: User.UserPassword,
		ActiveStart:  User.ActiveStart,
		ActiveEnd:    User.ActiveEnd,
		EntryUser:    User.EntryUser,
	}

	if err := user.AddEmployee(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
	}
	if user.ResultCode != 1 {
		respondWithError(w, http.StatusInternalServerError, user.ResultDescription, user.ResultCode)
	}
	result := map[string]interface{}{"status": 1}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *App) EditEmployee(w http.ResponseWriter, r *http.Request) {
	var User model.User
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&User); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}
	user := model.User{
		UserID:       User.UserID,
		UserName:     User.UserName,
		Role:         User.Role,
		SuperiorID:   User.SuperiorID,
		Email:        User.Email,
		UserPassword: User.UserPassword,
		ActiveStart:  User.ActiveStart,
		ActiveEnd:    User.ActiveEnd,
		ModifiedBy:   User.ModifiedBy,
	}

	if err := user.EditEmployee(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
	}
	if user.ResultCode != 1 {
		respondWithError(w, http.StatusInternalServerError, user.ResultDescription, user.ResultCode)
	}
	result := map[string]interface{}{"status": 1}
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
func respondWithError(w http.ResponseWriter, code int, message string, errorCode int) {
	respondWithJSON(w, code, map[string]interface{}{"status": errorCode, "error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func SyncDate(destDate string, baseDate string, sourceDate string) string {
	start, _ := time.Parse("2006-01-02 15:04:05", destDate)
	serverDate, _ := time.Parse("2006-01-02 15:04:05", sourceDate)
	clientDate, _ := time.Parse("2006-01-02 15:04:05", baseDate)
	syncedDate := start.Add(serverDate.Sub(clientDate))
	return syncedDate.Format("2006-01-02 15:04:05")
}

func checkValidityPeriod(destDate string, startDate string, endDate sql.NullString) bool {
	start, _ := time.Parse("2006-01-02 15:04:05", startDate)
	dest, _ := time.Parse("2006-01-02 15:04:05", destDate)
	if endDate.Valid {
		end, _ := time.Parse("2006-01-02 15:04:05", endDate.String)
		return dest.After(start) && dest.Before(end)
	} else {
		return dest.After(start)
	}
}

func (a *App) UploadToS3(file multipart.File, filepath string) (string, int) {
	token := ""
	creds := credentials.NewStaticCredentials(a.Config.AWS.S3.AccessKeyID, a.Config.AWS.S3.SecretAccessKey, token)
	_, err := creds.Get()
	if err != nil {
		return fmt.Sprintf("bad credentials: %s", err), 0
	}
	cfg := aws.NewConfig().WithRegion("us-east-1").WithCredentials(creds)
	svc := s3.New(session.New(), cfg)

	//file, err := os.Open("test.jpg")
	// if err != nil {
	// 	return fmt.Sprintf("err opening file: %s", err)
	// }
	// defer file.Close()
	//fileInfo, _ := file.Stat()
	//size := fileInfo.Size()
	//buffer := make([]byte, size) // read file content to buffer

	//file.Read(buffer)
	//fileBytes := bytes.NewReader(buffer)
	//fileType := http.DetectContentType(buffer)
	path := filepath + "." + a.Config.Etc.ScreenshotExt
	fileSize, err := file.Seek(0, 0)
	params := &s3.PutObjectInput{
		Bucket:        aws.String(a.Config.AWS.S3.BucketName),
		Key:           aws.String(path),
		Body:          file,
		ContentLength: aws.Int64(fileSize),
		ContentType:   aws.String("image/jpeg"),
	}
	resp, err := svc.PutObject(params)
	if err != nil {
		return fmt.Sprintf("bad response: %s", err.Error()), 0
	}
	return fmt.Sprintf("response %s", awsutil.StringValue(resp)), 1
}

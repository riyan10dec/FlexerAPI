package main

import (
	model "FlexerAPI/Model"
	query "FlexerAPI/Query"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	handlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	OriginObj := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	HeadersObj := handlers.AllowedHeaders([]string{"X-Requested-With", "Authorization", "Content-Type", "X-Auth-Token", "Origin"})
	MethodsObj := handlers.AllowedMethods([]string{"GET", "HEAD", "PUT", "OPTIONS"})
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(OriginObj, HeadersObj, MethodsObj)(a.Router)))
}

//ROUTES
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/ping", a.Ping).Methods("GET")
	a.Router.HandleFunc("/login", a.Login).Methods("POST")
	a.Router.Handle("/logout", jwtMiddleware.Handler(http.HandlerFunc(a.Logout))).Methods("POST")
	a.Router.Handle("/addActivity", jwtMiddleware.Handler(http.HandlerFunc(a.AddActivity))).Methods("POST")
	a.Router.Handle("/addActivity/screenshot", jwtMiddleware.Handler(http.HandlerFunc(a.AddActivityScreenshot))).Methods("POST")
	a.Router.Handle("/getTask/{sessionID}", jwtMiddleware.Handler(http.HandlerFunc(a.GetTask))).Methods("GET")
	a.Router.Handle("/addTask/{sessionID}", jwtMiddleware.Handler(http.HandlerFunc(a.GetTask))).Methods("POST")

	a.Router.HandleFunc("/cms/login", a.CMSLogin).Methods("POST")
	a.Router.HandleFunc("/cms/addEmployee1", a.AddEmployee).Methods("POST")
	a.Router.Handle("/cms/addEmployee", jwtMiddleware.Handler(http.HandlerFunc(a.AddEmployee))).Methods("POST")
	a.Router.Handle("/cms/editEmployee", jwtMiddleware.Handler(http.HandlerFunc(a.EditEmployee))).Methods("POST")
	a.Router.Handle("/cms/GetActiveSubs/{userID}", jwtMiddleware.Handler(http.HandlerFunc(a.GetActiveSubs))).Methods("GET")
	a.Router.Handle("/cms/CheckSubscription/{clientID}", jwtMiddleware.Handler(http.HandlerFunc(a.CheckSubscription))).Methods("GET")
	a.Router.Handle("/cms/EmployeeTree/first/{clientID}/{activeOnly}", jwtMiddleware.Handler(http.HandlerFunc(a.EmployeeTreeGetFirstLevel))).Methods("GET")
	a.Router.Handle("/cms/EmployeeTree/child/{userID}/{activeOnly}", jwtMiddleware.Handler(http.HandlerFunc(a.EmployeeTreeGetChild))).Methods("GET")
	a.Router.Handle("/cms/EmployeeTree/ChangeSuperior", jwtMiddleware.Handler(http.HandlerFunc(a.EmployeeTreeChangeSuperior))).Methods("POST")
	a.Router.Handle("/cms/EmailValidation", jwtMiddleware.Handler(http.HandlerFunc(a.EmailValidation))).Methods("POST")
	a.Router.Handle("/cms/GetAllEmployees/{userID}", jwtMiddleware.Handler(http.HandlerFunc(a.GetAllEmployees))).Methods("GET")
	a.Router.Handle("/cms/GetAllDepartments/{clientID}", jwtMiddleware.Handler(http.HandlerFunc(a.GetAllDepartment))).Methods("GET")
	a.Router.Handle("/cms/GetActiveDepartments/{clientID}", jwtMiddleware.Handler(http.HandlerFunc(a.GetActiveDepartment))).Methods("GET")
	a.Router.Handle("/cms/ChangePassword", jwtMiddleware.Handler(http.HandlerFunc(a.ChangePassword))).Methods("POST")
	a.Router.Handle("/cms/GetFeatures/{userID}/{positionName}/{subscriptionID}", jwtMiddleware.Handler(http.HandlerFunc(a.GetAllFeatures))).Methods("GET")
	a.Router.Handle("/cms/GetSubs/{userID}", jwtMiddleware.Handler(http.HandlerFunc(a.GetSubs))).Methods("GET")
	a.Router.Handle("/cms/GetAllActivities/{userID}", jwtMiddleware.Handler(http.HandlerFunc(a.GetAllActivities))).Methods("GET")
}

//HANDLERS

func (a *App) Ping(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "pong")
}

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
		return
	}

	if login.ResultCode == 1 {
		token := GetToken(strconv.Itoa(int(login.Session.SessionID)))
		result := map[string]interface{}{"status": login.ResultCode, "description": login.ResultDescription, "token": token, "session_id": login.Session.SessionID}
		respondWithJSON(w, http.StatusOK, result)
	} else {
		respondWithError(w, http.StatusInternalServerError, login.ResultDescription, login.ResultCode)
		return
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
		return
	}

	if logout.ResultCode == 1 {
		result := map[string]interface{}{"status": logout.ResultCode, "description": logout.ResultDescription}
		respondWithJSON(w, http.StatusOK, result)
	} else {
		respondWithError(w, http.StatusInternalServerError, logout.ResultDescription, logout.ResultCode)
		return
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
		return
	}
	//Preparing
	tx, err := a.DB.Begin()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}
	defer tx.Rollback()
	stmt, err := a.DB.Prepare(query.SearchQuery("createTransactionQuery"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
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
			return
		}
		if transaction.ResultCode != 1 {
			tx.Rollback()
			respondWithError(w, http.StatusInternalServerError, transaction.ResultDescription, transaction.ResultCode)
			return
		}
	}
	tx.Commit()
	result := map[string]interface{}{"status": 1, "description": "All Transaction Successfully Inserted"}
	respondWithJSON(w, http.StatusOK, result)
}
func (a *App) GetTask(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	var session model.Session
	vars := mux.Vars(r)
	var err error
	session.SessionID, err = strconv.ParseInt(vars["sessionID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}
	if err := session.GetTasks(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}

	result := map[string]interface{}{"task": session.Tasks}
	respondWithJSON(w, http.StatusOK, result)
}
func (a *App) AddTask(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	var task model.Task
	vars := mux.Vars(r)
	var err error
	task.Session.SessionID, err = strconv.ParseInt(vars["sessionID"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&task); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload : Task", -2)
		return
	}

	if err := task.AddTask(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
	}

	if task.ResultCode == 1 {
		result := map[string]interface{}{"status": task.ResultCode, "description": task.ResultDescription}
		respondWithJSON(w, http.StatusOK, result)
	} else {
		respondWithError(w, http.StatusInternalServerError, task.ResultDescription, task.ResultCode)
		return
	}
}

/* AddActivityScreenshot :
header : multipart/form-data
Params :
- screenshot : {file}
- date : "{date}"
*/
func (a *App) AddActivityScreenshot(w http.ResponseWriter, r *http.Request) {
	//var transactions []model.Transaction
	filedata, _, err := r.FormFile("img")
	if err != nil {
		return
	}
	var Session model.Session
	var Screenshot model.Screenshot
	Session.SessionID, _ = strconv.ParseInt(r.FormValue("sessionID"), 10, 64)
	Screenshot.SessionID = Session.SessionID
	Screenshot.ScreenshotDate = SyncDate(r.FormValue("screenshotDate"), Session.ClientDate, Session.ServerDate)
	Screenshot.ActivityName.String = r.FormValue("activityName")
	Screenshot.ActivityType.String = r.FormValue("activityType")

	if err := Session.FrontCheckSession(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}
	if !checkValidityPeriod(Screenshot.ScreenshotDate, Session.StartTime, Session.EndTime) {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload : ScreenshotDate outside Session Date", -2)
		return
	}

	//Getting Parameter
	if err := Screenshot.GetScreenshotParam(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}
	if Screenshot.ResultCode != 1 {
		respondWithError(w, http.StatusInternalServerError, Screenshot.ResultDescription.String, Screenshot.ResultCode)
		return
	}

	//Start Uploading File
	//Screenshot.ResultDescription, Screenshot.ResultCode = a.UploadToS3(file, Screenshot.Filename)
	Screenshot.Filename.String, Screenshot.ResultDescription.String, Screenshot.ResultCode = a.UploadToGoogleCloud(filedata, Screenshot.Filename.String)
	if Screenshot.ResultCode != 1 {
		respondWithError(w, http.StatusInternalServerError, Screenshot.ResultDescription.String, Screenshot.ResultCode)
		return
	}
	//Reporting Status
	if err := Screenshot.ReportScreenshotStatus(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}
	if Screenshot.ResultCode != 1 {
		respondWithError(w, http.StatusInternalServerError, Screenshot.ResultDescription.String, -1)
		return
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
		return
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
		ClientID:       User.ClientID,
		EmployeeID:     User.EmployeeID,
		UserName:       User.UserName,
		PositionName:   User.PositionName,
		DepartmentName: User.DepartmentName,
		SuperiorID:     User.SuperiorID,
		Email:          User.Email,
		UserPassword:   User.UserPassword,
		ActiveStart:    User.ActiveStart,
		ActiveEnd:      User.ActiveEnd,
		EntryUser:      User.EntryUser,
	}

	if err := user.AddEmployee(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}
	if user.ResultCode != 1 {
		respondWithError(w, http.StatusInternalServerError, user.ResultDescription, user.ResultCode)
		return
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
		UserID:         User.UserID,
		EmployeeID:     User.EmployeeID,
		UserName:       User.UserName,
		PositionName:   User.PositionName,
		DepartmentName: User.DepartmentName,
		SuperiorID:     User.SuperiorID,
		Email:          User.Email,
		UserPassword:   User.UserPassword,
		ActiveStart:    User.ActiveStart,
		ActiveEnd:      User.ActiveEnd,
		ModifiedBy:     User.ModifiedBy,
	}

	if err := user.EditEmployee(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
	}
	if user.ResultCode != 1 {
		respondWithError(w, http.StatusInternalServerError, user.ResultDescription, user.ResultCode)
		return
	}
	result := map[string]interface{}{"status": 1}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *App) GetActiveSubs(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	var User model.User
	vars := mux.Vars(r)
	var err error
	User.UserID, err = strconv.Atoi(vars["userID"])
	user := model.User{
		UserID: User.UserID,
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}

	if err := user.GetActiveSubs(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}
	var res []map[string]interface{}
	for _, u := range user.ReferenceUser {
		res = append(res, map[string]interface{}{
			"userID":         u.UserID,
			"employeeID":     u.EmployeeID,
			"userName":       u.UserName,
			"positionName":   u.PositionName,
			"departmentName": u.DepartmentName,
			"IPAddress":      u.IPAddress,
			"loginDate":      u.LoginDate,
		})
	}
	result := map[string]interface{}{"activeUsers": res}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *App) CheckSubscription(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	var Client model.Client
	vars := mux.Vars(r)
	var err error
	Client.ClientID, err = strconv.Atoi(vars["clientID"])
	client := model.Client{
		ClientID: Client.ClientID,
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}

	if err := client.CheckSubscription(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}
	result := map[string]interface{}{
		"subscriptionType":   client.SubscriptionType,
		"subscriptionStatus": client.SubscriptionStatus,
		"subscriptionStart":  client.SubscriptionStart,
		"subscriptionEnd":    client.SubscriptionEnd,
		"graceUntil":         client.GraceUntil,
		"maxUser":            client.MaxUser,
		"registeredMember":   client.RegisteredMember,
	}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *App) EmployeeTreeGetFirstLevel(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	var User model.User
	vars := mux.Vars(r)
	var err error
	User.ClientID, err = strconv.Atoi(vars["clientID"])
	User.ActiveOnly, err = strconv.ParseBool(vars["activeOnly"])
	user := model.User{
		ClientID: User.ClientID,
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}

	if err := user.EmployeeTreeFirstLevel(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}
	var res []map[string]interface{}
	for _, u := range user.ReferenceUser {
		res = append(res, map[string]interface{}{
			"userID":       u.UserID,
			"userName":     u.UserName,
			"subsCount":    u.SubsCount,
			"activeStatus": u.ActiveStatus,
		})
	}
	result := map[string]interface{}{"employees": res}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *App) EmployeeTreeGetChild(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	var User model.User
	vars := mux.Vars(r)
	var err error
	User.UserID, err = strconv.Atoi(vars["userID"])
	User.ActiveOnly, err = strconv.ParseBool(vars["activeOnly"])
	user := model.User{
		UserID: User.UserID,
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}

	if err := user.EmployeeTreeSubs(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}
	var res []map[string]interface{}
	for _, u := range user.ReferenceUser {
		res = append(res, map[string]interface{}{
			"userID":       u.UserID,
			"userName":     u.UserName,
			"subsCount":    u.SubsCount,
			"activeStatus": u.ActiveStatus,
		})
	}
	result := map[string]interface{}{"employees": res}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *App) EmployeeTreeChangeSuperior(w http.ResponseWriter, r *http.Request) {
	var User model.User
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&User); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}
	user := model.User{
		UserID:     User.UserID,
		SuperiorID: User.SuperiorID,
	}

	if _, err := user.EmployeeTreeChangeSuperior(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}
	// if res != sql.Result. {
	// 	respondWithError(w, http.StatusInternalServerError, user.ResultDescription, user.ResultCode)
	// }
	result := map[string]interface{}{"status": 1}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *App) EmailValidation(w http.ResponseWriter, r *http.Request) {
	var User model.User
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&User); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}
	user := model.User{
		UserID:   User.UserID,
		ClientID: User.ClientID,
		Email:    User.Email,
	}

	if err := user.EmailValidation(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}
	// if res != sql.Result. {
	// 	respondWithError(w, http.StatusInternalServerError, user.ResultDescription, user.ResultCode)
	// }
	var status int
	if user.ResultDescription == "OK" {
		status = 1
	} else {
		status = -1
	}
	result := map[string]interface{}{"status": status, "description": user.ResultDescription}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *App) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	var User model.User
	vars := mux.Vars(r)
	var err error
	User.UserID, err = strconv.Atoi(vars["userID"])
	user := model.User{
		UserID: User.UserID,
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}

	if err := user.GetEmployees(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}
	var res []map[string]interface{}
	for _, u := range user.ReferenceUser {
		res = append(res, map[string]interface{}{
			"userID":         u.UserID,
			"employeeID":     u.EmployeeID,
			"userName":       u.UserName,
			"positionName":   u.PositionName,
			"departmentName": u.DepartmentName,
			"activeStatus":   u.ActiveStatus,
			"lastActivity":   u.LastActivity.String,
		})
	}
	result := map[string]interface{}{"employees": res}
	respondWithJSON(w, http.StatusOK, result)
}

//--Department
func (a *App) SaveDepartment(w http.ResponseWriter, r *http.Request) {
	var Department model.Department
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&Department); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}
	department := model.Department{
		ClientID:             Department.ClientID,
		DepartmentsSeparator: strings.Join(Department.DepartmentList, "|"),
		EntryBy:              Department.EntryBy,
	}

	if err := department.SaveDepartment(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}
	if department.ResultCode != 1 {
		respondWithError(w, http.StatusInternalServerError, department.ResultDescription, department.ResultCode)
		return
	}
	result := map[string]interface{}{"status": 1}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *App) GetActiveDepartment(w http.ResponseWriter, r *http.Request) {
	var Department model.Department
	vars := mux.Vars(r)
	var err error
	Department.ClientID, err = strconv.Atoi(vars["clientID"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}
	department := model.Department{
		ClientID: Department.ClientID,
	}
	var ds []model.Department
	if err, ds = department.GetActiveDepartments(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}

	var res []string
	for _, d := range ds {
		res = append(res, d.DepartmentName)
	}
	result := map[string]interface{}{"status": 1, "departments": res}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *App) GetAllDepartment(w http.ResponseWriter, r *http.Request) {
	var Department model.Department
	vars := mux.Vars(r)
	var err error
	Department.ClientID, err = strconv.Atoi(vars["clientID"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}
	department := model.Department{
		ClientID: Department.ClientID,
	}
	var ds []model.Department
	if err, ds = department.GetAllDepartments(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}

	var res []map[string]interface{}
	for _, d := range ds {
		res = append(res, map[string]interface{}{
			"selected":       d.Selected,
			"departmentName": d.DepartmentName,
		})
	}
	result := map[string]interface{}{"status": 1, "departments": res}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *App) GetAllFeatures(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	var User model.User
	vars := mux.Vars(r)
	var err error
	User.UserID, err = strconv.Atoi(vars["userID"])
	User.PositionName = vars["positionName"]
	User.SubscriptionID, err = strconv.Atoi(vars["subscriptionID"])
	user := model.User{
		UserID:         User.UserID,
		PositionName:   User.PositionName,
		SubscriptionID: User.SubscriptionID,
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}

	if err := user.GetFeatures(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}
	var res []map[string]interface{}
	for _, f := range user.Features {
		res = append(res, map[string]interface{}{
			"featureID":          f.FeatureID,
			"featureName":        f.FeatureName,
			"featureType":        f.FeatureType,
			"featureDescription": f.FeatureDescription,
		})
	}
	result := map[string]interface{}{"features": res}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *App) GetSubs(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	var User model.User
	vars := mux.Vars(r)
	var err error
	User.UserID, err = strconv.Atoi(vars["userID"])
	user := model.User{
		UserID: User.UserID,
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}

	if err := user.GetSubs(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}
	var res []map[string]interface{}
	for _, u := range user.ReferenceUser {
		res = append(res, map[string]interface{}{
			"userID":         u.UserID,
			"employeeID":     u.EmployeeID,
			"userName":       u.UserName,
			"positionName":   u.PositionName,
			"departmentName": u.DepartmentName,
			"activeStatus":   u.ActiveStatus,
			"lastActivity":   u.LastActivity.String,
		})
	}
	result := map[string]interface{}{"employees": res}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *App) GetAllActivities(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	var User model.User
	vars := mux.Vars(r)
	var err error
	User.UserID, err = strconv.Atoi(vars["userID"])
	user := model.User{
		UserID: User.UserID,
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}

	if err := user.GetAllActivities(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}
	var res []map[string]interface{}
	for _, a := range user.Activities {
		res = append(res, map[string]interface{}{
			"activityName":   a.ActivityName,
			"activityType":   a.ActivityType,
			"category":       a.Category,
			"classification": a.Classification,
			"utilization":    a.Utilization,
		})
	}
	result := map[string]interface{}{"activities": res}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *App) ChangePassword(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	var User model.User
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&User); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", -2)
		return
	}
	user := model.User{
		UserID:      User.UserID,
		OldPassword: User.OldPassword,
		NewPassword: User.NewPassword,
	}
	if err := user.ChangePassword(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), -1)
		return
	}

	if user.ResultCode == 1 {
		result := map[string]interface{}{"status": user.ResultCode, "description": user.ResultDescription}
		respondWithJSON(w, http.StatusOK, result)
	} else {
		respondWithError(w, http.StatusInternalServerError, user.ResultDescription, user.ResultCode)
		return
	}
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
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
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

func (a *App) UploadToGoogleCloud(file multipart.File, filepath string) (string, string, int) {
	scope := storage.DevstorageFullControlScope
	client, err := google.DefaultClient(context.Background(), scope)
	filepath = strings.Replace(filepath, "[ext]", "jpeg", -1)
	if err != nil {
		log.Fatalf("Unable to get default client: %v", err)
	}
	service, err := storage.New(client)
	if err != nil {
		log.Fatalf("Unable to create storage service: %v", err)
	}
	object := &storage.Object{
		Name: a.Config.Gcs.ScreenshotFolder + filepath,
	}
	//file, err := os.Open(*fileName)
	// if err != nil {
	// 	fatalf(service, "Error opening %q: %v", *fileName, err)
	// }
	if res, err := service.Objects.Insert(a.Config.Gcs.Bucket, object).Media(file).Do(); err == nil {
		fmt.Printf("Created object %v at location %v\n\n", res.Name, res.SelfLink)
	} else {
		return filepath, err.Error(), 0
		//fatalf(service, "Objects.Insert failed: %v", err)
	}
	return filepath, "", 1
}
func fieldSet(fields ...string) map[string]bool {
	set := make(map[string]bool, len(fields))
	for _, s := range fields {
		set[s] = true
	}
	return set
}
func SelectFields(s interface{}, fields ...string) map[string]interface{} {
	fs := fieldSet(fields...)
	rt, rv := reflect.TypeOf(&s), reflect.ValueOf(&s)
	out := make(map[string]interface{}, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		jsonKey := field.Tag.Get("json")
		if fs[jsonKey] {
			out[jsonKey] = rv.Field(i).Interface()
		}
	}
	return out
}

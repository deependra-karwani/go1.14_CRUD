package controllers

import (
	"CRUD/config"
	"CRUD/structs"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var db = config.GetDB()

func Register(w http.ResponseWriter, r *http.Request, done chan<- bool) {
	defer func() {
		done <- true
	}()
	w.Header().Set("Content-Type", "application/json")

	// E-mail Verification for Authentication

	r.ParseMultipartForm(50 << 20)

	profPic, handler, err := r.FormFile("prof")
	if err != nil {
		config.SendBadReqResponse(w, `{"message": "Invalid Request"}`)
		return
	}
	defer profPic.Close()

	if err := config.SaveFile(profPic, "../images/", handler.Filename); err != nil {
		config.SendBadReqResponse(w, `{"message": "Could not save image. Please Try Again."}`)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		config.SendBadReqResponse(w, `{"message": "Could not Secure Credentials. Please Try Again."}`)
		return
	}
	generatedPassword := string(hashedPassword)

	token, err := config.GenToken(r.FormValue("email"))
	if err != nil {
		config.SendBadReqResponse(w, `{"message": "Could not Complete Registration"}`)
		return
	}

	var userid int
	stmt := "INSERT INTO users(name, email, mobile, username, password, fcm, token, profPic) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	if err := db.QueryRow(stmt, r.FormValue("name"), r.FormValue("email"), r.FormValue("mobile"), r.FormValue("username"), generatedPassword, r.FormValue("fcm"), token, handler.Filename).Scan(&userid); err != nil {
		config.SendBadReqResponse(w, `{"message": "Could not Complete Registration. Please Try Again."}`)
		return
	}

	response := structs.SessSuccess{
		Message: "Registration Successful",
		UserId:  userid,
	}

	w.Header().Add("token", token)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		config.SendForbiddenResponse(w, `{"message": "Session could not be Established"}`)
	}
}

func Login(w http.ResponseWriter, r *http.Request, done chan<- bool) {
	defer func() {
		done <- true
	}()
	w.Header().Set("Content-Type", "application/json")

	var login structs.Login
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&login); err != nil {
		config.SendBadReqResponse(w, `{"message": "Invalid Request"}`)
		return
	}

	var forSess structs.ForSession
	// stmt := "SELECT email, password FROM users WHERE username = $1 AND token = NULL" // Single Session Only
	stmt := "SELECT id, email, password FROM users WHERE username = $1"
	if err := db.QueryRow(stmt, login.Username).Scan(&forSess.UserId, &forSess.Email, &forSess.HashedPassword); err != nil {
		config.SendBadReqResponse(w, `{"message": "Could not Verify Details. Please Try Again."}`)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(forSess.HashedPassword), []byte(login.Password)); err != nil {
		config.SendBadReqResponse(w, `{"message": "Incorrect Password"}`)
		return
	}

	token, err := config.GenToken(forSess.Email)
	if err != nil {
		config.SendBadReqResponse(w, `{"message": "Could not Complete Registration"}`)
		return
	}

	stmt = "UPDATE users SET token = $1, fcm = $2 WHERE id = $3"
	res, err := db.Exec(stmt, token, login.FCM, forSess.UserId)
	if err != nil {
		config.SendBadReqResponse(w, `{"message": "Could not Establish Session. Please Try Again."}`)
		return
	}

	count, err2 := res.RowsAffected()
	if count == 0 || err2 != nil {
		config.SendBadReqResponse(w, `{"message": "Could not Establish Session. Please Try Again."}`)
		return
	}

	response := structs.SessSuccess{
		Message: "Login Successful",
		UserId:  forSess.UserId,
	}

	w.Header().Add("token", token)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		config.SendForbiddenResponse(w, `{"message": "Session could not be Established"}`)
	}
}

func ForgotPassword(w http.ResponseWriter, r *http.Request, done chan<- bool) {
	defer func() {
		done <- true
	}()
	w.Header().Set("Content-Type", "application/json")

	// E-mail Verification for Authentication

	var forgot structs.Forgot
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&forgot); err != nil {
		config.SendBadReqResponse(w, `json:"Invalid Request"`)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(forgot.Password), bcrypt.DefaultCost)
	if err != nil {
		config.SendBadReqResponse(w, `{"message": "Could not Secure Credentials. Please Try Again."}`)
		return
	}
	generatedPassword := string(hashedPassword)

	stmt := "UPDATE users SET password = $1 WHERE email = $2"
	res, err := db.Exec(stmt, generatedPassword, forgot.Email)
	if err != nil {
		config.SendBadReqResponse(w, `{"message": "Could not Change Password. Please Try Again."}`)
		return
	}

	count, err2 := res.RowsAffected()
	if count == 0 || err2 != nil {
		config.SendBadReqResponse(w, `{"message": "Could not Change Password. Please Try Again."}`)
		return
	}

	config.SendSuccessResponse(w, `{"message": "Password Changed Successfully"}`)
}

func Logout(w http.ResponseWriter, r *http.Request, done chan<- bool) {
	defer func() {
		done <- true
	}()
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	if query["userid"] == nil {
		config.SendBadReqResponse(w, `{"message": "Invalid Request"}`)
		return
	}

	// stmt := UPDATE users SET token = NULL WHERE token = $1" // Single Session Only
	stmt := "UPDATE users SET token = NULL WHERE id = $1"
	res, err := db.Exec(stmt, query["userid"][0])
	if err != nil {
		config.SendBadReqResponse(w, `{"message": "Could not end Session. Please Try Again."}`)
		return
	}

	count, err2 := res.RowsAffected()
	if count == 0 || err2 != nil {
		config.SendBadReqResponse(w, `{"message": "Could not end Session"}`)
		return
	}

	config.SendSuccessResponse(w, `{"message": "Logout Successful"}`)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request, done chan<- bool) {
	defer func() {
		done <- true
	}()
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	if query["userid"] == nil {
		config.SendBadReqResponse(w, `{"message": "Invalid Request"}`)
		return
	}

	// stmt := "SELECT profPic, name, username FROM users" // All Users
	stmt := "SELECT id, profPic, name, username FROM users WHERE id <> $1"
	result, err := db.Query(stmt, query["userid"][0])
	if err != nil {
		config.SendBadReqResponse(w, `{"message": "Could not Fetch List. Please Try Again."}`)
		return
	}

	defer result.Close()
	var response structs.UserShortResponse
	for result.Next() {
		var user structs.UserShort

		if err := result.Scan(&user.Id, &user.ProfPic, &user.Name, &user.Username); err != nil {
			config.SendBadReqResponse(w, `{"message": "Could not Fetch List. Please Try Again."}`)
			return
		}

		response.Users = append(response.Users, user)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		config.SendBadReqResponse(w, `{"message": "Could not Fetch List. Please Try Again."}`)
	}
}

func GetUserDetails(w http.ResponseWriter, r *http.Request, done chan<- bool) {
	defer func() {
		done <- true
	}()
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	if query["userid"] == nil {
		config.SendBadReqResponse(w, `{"message": "Invalid Request"}`)
		return
	}

	var (
		response structs.UserDetailResponse
		details  structs.UserDetails
	)
	// stmt := "SELECT profPic, name, username, email, mobile FROM users WHERE id = $1 AND token = $2" // Only own Details, (r.Header.Get("token"))
	stmt := "SELECT profPic, name, username, email, mobile FROM users WHERE id = $1"
	if err := db.QueryRow(stmt, &details.ProfPic, &details.Name, &details.Username, &details.Email, &details.Mobile); err != nil {
		config.SendBadReqResponse(w, `{"message": "Could not Fetch Details. Please Try Again."}`)
		return
	}

	response.User = details

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		config.SendBadReqResponse(w, `{"message": "Could not Fetch Details. Please Try Again."}`)
	}
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request, done chan<- bool) {
	defer func() {
		done <- true
	}()
	w.Header().Set("Content-Type", "application/json")

	// E-Mail Update restricted as it is to be Verified

	r.ParseMultipartForm(50 << 20)

	var (
		res sql.Result
		err error
	)

	profPic, handler, err := r.FormFile("prof")
	if err != nil {
		stmt := "UPDATE users SET name = $1, username = $2, mobile = $3, WHERE id = $4 AND token = $5"
		res, err = db.Exec(stmt, r.FormValue("name"), r.FormValue("username"), r.FormValue("mobile"), r.FormValue("userid"), r.Header.Get("token"))
	} else {
		if err := config.SaveFile(profPic, "../images/", handler.Filename); err != nil {
			config.SendBadReqResponse(w, `{"message": "Could not save image. Please Try Again."}`)
			return
		}
		stmt := "UPDATE users SET name = $1, username = $2, mobile = $3, profPic = $4 WHERE id = $5 AND token = $6"
		res, err = db.Exec(stmt, r.FormValue("name"), r.FormValue("username"), r.FormValue("mobile"), handler.Filename, r.FormValue("userid"), r.Header.Get("token"))
	}

	if err != nil {
		config.SendBadReqResponse(w, `{"message": "Could not Update Profile. Please Try Again."}`)
		return
	}

	count, err2 := res.RowsAffected()
	if count == 0 || err2 != nil {
		config.SendBadReqResponse(w, `{"message": "Invalid Request"}`)
		return
	}

	config.SendSuccessResponse(w, `{"message": "Profile Updated Successfully"}`)
}

func DeleteUserAccount(w http.ResponseWriter, r *http.Request, done chan<- bool) {
	defer func() {
		done <- true
	}()
	w.Header().Set("Content-Type", "application/json")

	var user structs.UserId
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		config.SendBadReqResponse(w, `{"message": "Invalid Request"}`)
		return
	}

	stmt := "DELETE FROM users WHERE token = $1 AND id = $2"
	res, err := db.Exec(stmt, r.Header.Get("token"))
	if err != nil {
		config.SendBadReqResponse(w, `{"message": "Could not Delete Account. Please Try Again."}`)
		return
	}

	count, err2 := res.RowsAffected()
	if count == 0 || err2 != nil {
		config.SendBadReqResponse(w, `{"message": "Invalid Request"}`)
		return
	}

	config.SendSuccessResponse(w, `{"message": "Account Deleted Successfully"}`)
}

func RefreshToken(w http.ResponseWriter, r *http.Request, done chan<- bool) {
	defer func() {
		done <- true
	}()
	w.Header().Set("Content-Type", "application/json")

	email := fmt.Sprintf("%s", r.Context().Value("userid"))

	token, err := config.GenToken(email)
	if err != nil {
		config.SendBadReqResponse(w, `{"message": "Could Not Establish Session. Please Try Again."}`)
		return
	}

	stmt := "UPDATE handlers SET token = $1 WHERE email = $2"
	res, err := db.Exec(stmt, token, email)

	if err != nil {
		config.SendBadReqResponse(w, `{"message": "Could Not Re-Establish Session. Please Try Again."}`)
		return
	}

	count, err2 := res.RowsAffected()
	if count == 0 || err2 != nil {
		config.SendBadReqResponse(w, `{"message": "Could Not Re-Establish Session. Please Try Again."}`)
		return
	}

	w.Header().Add("token", token)
	config.SendSuccessResponse(w, `{"message": "Session Renewed"}`)
}

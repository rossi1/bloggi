package users

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

const hashRounds = 14

var validate = validator.New()

type userLoginParam struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// updateUserParam
type updateUserParam struct {
	FullName string `json:"fullname" validate:"required"`
	Username string `json:"username" validate:"required"`
}

func validateUsername(fl validator.FieldLevel) bool {
	var user User
	db := ReturnDBInstance()
	username := fl.Field().String()
	db.Where("username = ?", username).First(&user)

	if user.ID == 0 {
		return true
	}

	return false

}

var _ = validate.RegisterValidation("unique_username", validateUsername)

// RegisterUser Controller
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	defer r.Body.Close()

	if validateRequest := validate.Struct(user); validateRequest != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validateRequest.Error())
		return
	}
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), hashRounds)

	user.Password = string(hashPassword)
	user.CreatedAt = time.Now()
	db := ReturnDBInstance()
	db.Create(&user)
	response := map[string]interface{}{}
	if user.ID != 0 {
		w.WriteHeader(http.StatusCreated)
		token, _ := CreateToken(user)
	
		json.NewEncoder(w).Encode(user.requestResponse(token))

	} else {
		w.WriteHeader(http.StatusBadRequest)
		response["error"] = "Failed"
		response["message"] = "Failed to created user"
		json.NewEncoder(w).Encode(response)

	}

}

//LoginUser Controller
func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userLogin userLoginParam
	var user User
	json.NewDecoder(r.Body).Decode(&userLogin)
	defer r.Body.Close()
	if validateRequest := validate.Struct(userLogin); validateRequest != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validateRequest.Error())
		return
	}
	db := ReturnDBInstance()

	db.Where("username = ?", userLogin.Username).Find(&user)
	if user.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("user not found")

	} else {
		errChecking := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))

		if errChecking != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Invalid Credentails")

		} else {
			w.WriteHeader(http.StatusOK)
			token, _ := CreateToken(user)
			json.NewEncoder(w).Encode(user.requestResponse(token))

		}

	}

}

// UserProfile Controller
func UserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	param := mux.Vars(r)["id"]
	db := ReturnDBInstance()
	id, _ := strconv.Atoi(param)
	db.First(&user, id)
	if user.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("user not found")
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user.requestResponse(""))
	}

}

// UpdateProfile controller
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	var updateparam updateUserParam
	param := mux.Vars(r)["id"]
	db := ReturnDBInstance()
	id, _ := strconv.Atoi(param)

	db.First(&user, id)
	if user.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("user not found")

	} else {
		json.NewDecoder(r.Body).Decode(&updateparam)

		defer r.Body.Close()

		if validateRequest := validate.Struct(updateparam); validateRequest != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validateRequest.Error())
			return
		}
		user.FullName = updateparam.FullName
		user.Username = updateparam.Username
		db.Save(&user)
		w.WriteHeader(http.StatusOK)
		
		json.NewEncoder(w).Encode(user.requestResponse(""))

	}

}

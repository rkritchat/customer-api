package customer

import (
	"encoding/json"
	"errors"
	"ex_produce/internal/repository"
	"net/http"
)

type Service interface {
	GetUser(w http.ResponseWriter, r *http.Request)
	SaveUser(w http.ResponseWriter, r *http.Request)
}

type service struct {
	customerRepo repository.CustomerRepo
}

func NewService(customerRepo repository.CustomerRepo) Service {
	return &service{
		customerRepo: customerRepo,
	}
}

//request
type GetUserReq struct {
	Username string `json:"username"`
}

//response
type GetUserResp struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s service) GetUser(w http.ResponseWriter, r *http.Request) {
	var req GetUserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request message1", http.StatusBadRequest)
		return
	}
	if len(req.Username) == 0 {
		http.Error(w, "invalid request message2", http.StatusBadRequest)
		return
	}

	entity, err := s.customerRepo.FindByUsername(req.Username)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(&GetUserResp{
		Username: entity.Username,
		Password: entity.Password,
	})
}

type SaveReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

//http status
//200 OK!
//400 Bad request
//500 inter server error
func (s service) SaveUser(w http.ResponseWriter, r *http.Request) {
	//convert json to struct
	var req SaveReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invali request message", http.StatusBadRequest)
		return
	}

	//validate request
	err = validateRequest(req)
	if err != nil {
		http.Error(w, "invalid request message", http.StatusBadRequest)
		return
	}

	//save
	err = s.customerRepo.Save(repository.UserEntity{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		http.Error(w, "internal server err", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte("OK"))
}

func validateRequest(req SaveReq) error {
	if len(req.Username) == 0 {
		return errors.New("username is required")
	}

	if len(req.Password) == 0 {
		return errors.New("password is required")
	}

	if len(req.Email) == 0 {
		return errors.New("email is required")
	}

	return nil
}

package customer

import (
	"bytes"
	"encoding/json"
	"errors"
	"ex_produce/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mockErr = errors.New("mock err")

func Test_validateRequest(t *testing.T) {
	tt := []struct {
		name        string
		req         SaveReq
		expectedErr bool
	}{
		{
			name: "Should error when username is empty",
			req: SaveReq{
				Username: "",
				Password: "test_password",
				Email:    "test_email",
			},
			expectedErr: true,
		},
		{
			name: "Should error when password is empty",
			req: SaveReq{
				Username: "test_username",
				Password: "",
				Email:    "test_email",
			},
			expectedErr: true,
		},
		{
			name: "Should error when email is empty",
			req: SaveReq{
				Username: "test_username",
				Password: "test_password",
				Email:    "",
			},
			expectedErr: true,
		},
		{
			name: "Case success",
			req: SaveReq{
				Username: "test_username",
				Password: "test_password",
				Email:    "test_email",
			},
			expectedErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := validateRequest(tc.req)
			if tc.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_GetUser(t *testing.T) {
	//valid request
	req := GetUserReq{
		Username: "test_user",
	}
	reqJson, _ := json.Marshal(&req)

	//invalid request
	invalidReq := GetUserReq{
		Username: "",
	}
	invalidReqJson, _ := json.Marshal(invalidReq)

	tt := []struct {
		name             string
		w                *httptest.ResponseRecorder
		r                *http.Request
		customerRepo     repository.CustomerRepo
		expectHttpStatus int
	}{
		{
			name:             "Should return 200 when success",
			w:                httptest.NewRecorder(),
			r:                httptest.NewRequest(http.MethodPost, "http://127.0.0.9000/test", bytes.NewBuffer(reqJson)),
			customerRepo:     initCustRepoMock("OK"),
			expectHttpStatus: http.StatusOK,
		},
		{
			name:             "Should return 400 when username is empty",
			w:                httptest.NewRecorder(),
			r:                httptest.NewRequest(http.MethodPost, "http://127.0.0.1:9000/test", bytes.NewBuffer(invalidReqJson)),
			customerRepo:     initCustRepoMock("OK"),
			expectHttpStatus: http.StatusBadRequest,
		},
		{
			name:             "Should return 500 when error occurred",
			w:                httptest.NewRecorder(),
			r:                httptest.NewRequest(http.MethodPost, "http://127.0.0.1:9000/test", bytes.NewBuffer(reqJson)),
			customerRepo:     initCustRepoMock("!OK"),
			expectHttpStatus: http.StatusInternalServerError,
		},
		{
			name:             "Should return 400 when request in invalid json format",
			w:                httptest.NewRecorder(),
			r:                httptest.NewRequest(http.MethodPost, "http://127.0.0.1:9000/test", bytes.NewBuffer([]byte("afdsfa"))),
			customerRepo:     initCustRepoMock("OK"),
			expectHttpStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sv := NewService(tc.customerRepo)
			sv.GetUser(tc.w, tc.r)
			assert.Equal(t, tc.expectHttpStatus, tc.w.Code)
		})
	}
}

// == mock area ===
type custRepoMock struct {
	mock.Mock
}

func (m *custRepoMock) Save(entity repository.UserEntity) error {
	args := m.Called(entity)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (m *custRepoMock) FindByUsername(username string) (*repository.UserEntity, error) {
	args := m.Called(username)
	if args.Get(0) != nil { //success OK!
		return args.Get(0).(*repository.UserEntity), nil
	}
	return nil, args.Error(1)
}

//success
//failed
func initCustRepoMock(tc string) *custRepoMock {
	m := new(custRepoMock)
	switch tc {
	case "OK":
		m.On("Save", mock.AnythingOfType("repository.UserEntity")).Return(nil)
		m.On("FindByUsername", mock.AnythingOfType("string")).Return(&repository.UserEntity{ID: int64(1)}, nil)
	case "!OK":
		m.On("Save", mock.AnythingOfType("repository.UserEntity")).Return(mockErr)
		m.On("FindByUsername", mock.AnythingOfType("string")).Return(nil, mockErr)
	}

	return m
}

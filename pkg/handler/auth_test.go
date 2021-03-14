package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"github.com/psihachina/windfarms-backend/models"
	"github.com/psihachina/windfarms-backend/pkg/service"
	service_mocks "github.com/psihachina/windfarms-backend/pkg/service/mocks"
	"github.com/psihachina/windfarms-backend/utils"
)

func TestHandler_signUp(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *service_mocks.MockAuthorization, user models.User)

	user := randomUser(false)
	token := utils.RandomString(6)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            models.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: fmt.Sprintf(`{"email": "%s", "password": "%s"}`, user.Email, user.Password),
			inputUser: user,
			mockBehavior: func(r *service_mocks.MockAuthorization, user models.User) {
				r.EXPECT().CreateUser(user).Return(nil)
				r.EXPECT().GenerateToken(user.Email, user.Password).Return(token, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: fmt.Sprintf(`{"email":"%s","token":"%s"}`, user.Email, token),
		},
		{
			name:      "Wrong Input",
			inputBody: fmt.Sprintf(`{"email": "%s"}`, user.Email),
			inputUser: models.User{},
			mockBehavior: func(r *service_mocks.MockAuthorization, user models.User) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Error",
			inputBody: fmt.Sprintf(`{"email":"%s","password":"%s"}`, user.Email, user.Password),
			inputUser: user,
			mockBehavior: func(r *service_mocks.MockAuthorization, user models.User) {
				r.EXPECT().CreateUser(user).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			auth := service_mocks.NewMockAuthorization(c)
			test.mockBehavior(auth, test.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			// Init Endpoint
			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *service_mocks.MockAuthorization, user models.User)

	user := randomUser(false)
	token := utils.RandomString(6)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            models.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: fmt.Sprintf(`{"email": "%s", "password": "%s"}`, user.Email, user.Password),
			inputUser: user,
			mockBehavior: func(r *service_mocks.MockAuthorization, user models.User) {
				r.EXPECT().GenerateToken(user.Email, user.Password).Return(token, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: fmt.Sprintf(`{"email":"%s","token":"%s"}`, user.Email, token),
		},
		{
			name:      "Wrong Input",
			inputBody: fmt.Sprintf(`{"email": "%s"}`, user.Email),
			inputUser: models.User{},
			mockBehavior: func(r *service_mocks.MockAuthorization, user models.User) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Error",
			inputBody: fmt.Sprintf(`{"email":"%s","password":"%s"}`, user.Email, user.Password),
			inputUser: user,
			mockBehavior: func(r *service_mocks.MockAuthorization, user models.User) {
				r.EXPECT().GenerateToken(user.Email, user.Password).Return("", errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			auth := service_mocks.NewMockAuthorization(c)
			test.mockBehavior(auth, test.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			// Init Endpoint
			r := gin.New()
			r.POST("/sign-in", handler.signIn)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func randomUser(emptyField bool) models.User {
	var email string
	if emptyField {
		email = ""
	} else {
		email = utils.RandomEmail()
	}
	return models.User{
		Email:    email,
		Password: utils.RandomString(6),
	}
}

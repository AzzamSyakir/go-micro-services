package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-micro-services/src/auth-service/entity"
	model_request "go-micro-services/src/auth-service/model/request/controller"
	model_response "go-micro-services/src/auth-service/model/response"
	"net/http"
	"testing"

	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type UserWeb struct {
	Test *testing.T
	Path string
}

func NewUserWeb(test *testing.T) *UserWeb {
	userWeb := &UserWeb{
		Test: test,
		Path: "users",
	}
	return userWeb
}

func (userWeb *UserWeb) Start() {
	userWeb.Test.Run("UserWeb_Fetchuser_Succeed", userWeb.ListUser)
	userWeb.Test.Run("UserWeb_DeleteOneById_Succeed", userWeb.DeleteOneById)
	userWeb.Test.Run("UserWeb_UpdateUser_Succeed", userWeb.UpdateUser)
	userWeb.Test.Run("UserWeb_GetUserById_Succeed", userWeb.GetUserById)
	userWeb.Test.Run("UserWeb_GetUserByEmail_Succeed", userWeb.GetUserByEmail)
}
func (userWeb *UserWeb) ListUser(t *testing.T) {
	t.Parallel()

	testWeb := GetTestWeb()
	testWeb.AllSeeder.Up()
	defer testWeb.AllSeeder.Down()

	url := fmt.Sprintf("%s/%s", testWeb.Server.URL, userWeb.Path)
	request, newRequestErr := http.NewRequest(http.MethodGet, url, http.NoBody)
	if newRequestErr != nil {
		t.Fatal(newRequestErr)
	}
	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
	request.Header.Set("Authorization", "Bearer "+selectedSessionMock.AccessToken.String)
	response, doErr := http.DefaultClient.Do(request)
	if newRequestErr != nil {
		t.Fatal(newRequestErr)
	}
	if doErr != nil {
		t.Fatal(doErr)
	}

	bodyResponse := &model_response.Response[[]*entity.User]{}
	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(t, bodyResponse.Code, http.StatusOK)
}

func (userWeb *UserWeb) GetUserById(t *testing.T) {
	t.Parallel()

	testWeb := GetTestWeb()
	testWeb.AllSeeder.Up()
	defer testWeb.AllSeeder.Down()

	selectedUserMock := testWeb.AllSeeder.User.UserMock.Data[0]

	url := fmt.Sprintf("%s/%s/%s", testWeb.Server.URL, userWeb.Path, selectedUserMock.Id.String)
	request, newRequestErr := http.NewRequest(http.MethodGet, url, http.NoBody)
	if newRequestErr != nil {
		t.Fatal(newRequestErr)
	}
	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
	request.Header.Set("Authorization", "Bearer "+selectedSessionMock.AccessToken.String)
	response, doErr := http.DefaultClient.Do(request)
	if newRequestErr != nil {
		t.Fatal(newRequestErr)
	}
	if doErr != nil {
		t.Fatal(doErr)
	}

	bodyResponse := &model_response.Response[*entity.User]{}
	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(t, selectedUserMock.Id, bodyResponse.Data.Id)
	assert.Equal(t, selectedUserMock.Name, bodyResponse.Data.Name)
	assert.Equal(t, selectedUserMock.Email, bodyResponse.Data.Email)
	assert.Equal(t, selectedUserMock.Balance, bodyResponse.Data.Balance)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(bodyResponse.Data.Password.String), []byte(selectedUserMock.Password.String)))
}

func (userWeb *UserWeb) GetUserByEmail(t *testing.T) {
	t.Parallel()

	testWeb := GetTestWeb()
	testWeb.AllSeeder.Up()
	defer testWeb.AllSeeder.Down()

	selectedUserMock := testWeb.AllSeeder.User.UserMock.Data[0]

	url := fmt.Sprintf("%s/%s/email/%s", testWeb.Server.URL, userWeb.Path, selectedUserMock.Email.String)
	request, newRequestErr := http.NewRequest(http.MethodGet, url, http.NoBody)
	if newRequestErr != nil {
		t.Fatal(newRequestErr)
	}
	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
	request.Header.Set("Authorization", "Bearer "+selectedSessionMock.AccessToken.String)
	response, doErr := http.DefaultClient.Do(request)
	if doErr != nil {
		t.Fatal(doErr)
	}

	bodyResponse := &model_response.Response[*entity.User]{}
	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(t, selectedUserMock.Id, bodyResponse.Data.Id)
	assert.Equal(t, selectedUserMock.Name, bodyResponse.Data.Name)
	assert.Equal(t, selectedUserMock.Email, bodyResponse.Data.Email)
	assert.Equal(t, selectedUserMock.Balance, bodyResponse.Data.Balance)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(bodyResponse.Data.Password.String), []byte(selectedUserMock.Password.String)))
}

func (userWeb *UserWeb) UpdateUser(t *testing.T) {
	t.Parallel()

	testWeb := GetTestWeb()
	testWeb.AllSeeder.Up()
	defer testWeb.AllSeeder.Down()

	selectedUserMock := testWeb.AllSeeder.User.UserMock.Data[0]

	bodyRequest := &model_request.UserPatchOneByIdRequest{}
	bodyRequest.Name = null.NewString(selectedUserMock.Name.String+"patched", true)
	bodyRequest.Email = null.NewString(selectedUserMock.Email.String+"patched", true)
	bodyRequest.Balance = null.NewInt(selectedUserMock.Balance.Int64, true)
	bodyRequest.Password = null.NewString(selectedUserMock.Password.String+"patched", true)

	bodyRequestJsonByte, marshalErr := json.Marshal(bodyRequest)
	if marshalErr != nil {
		t.Fatal(marshalErr)
	}
	bodyRequestBuffer := bytes.NewBuffer(bodyRequestJsonByte)

	url := fmt.Sprintf("%s/%s/%s", testWeb.Server.URL, userWeb.Path, selectedUserMock.Id.String)
	request, newRequestErr := http.NewRequest(http.MethodPatch, url, bodyRequestBuffer)
	if newRequestErr != nil {
		t.Fatal(newRequestErr)
	}
	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
	request.Header.Set("Authorization", "Bearer "+selectedSessionMock.AccessToken.String)
	response, doErr := http.DefaultClient.Do(request)
	if doErr != nil {
		t.Fatal(doErr)
	}

	bodyResponse := &model_response.Response[*entity.User]{}
	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(t, selectedUserMock.Id, bodyResponse.Data.Id)
	assert.Equal(t, bodyRequest.Name, bodyResponse.Data.Name)
	assert.Equal(t, bodyRequest.Email, bodyResponse.Data.Email)
	assert.Equal(t, bodyRequest.Balance, bodyResponse.Data.Balance)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(bodyResponse.Data.Password.String), []byte(bodyRequest.Password.String)))
}

func (userWeb *UserWeb) DeleteOneById(t *testing.T) {
	t.Parallel()

	testWeb := GetTestWeb()
	testWeb.AllSeeder.Up()
	defer testWeb.AllSeeder.Down()

	selectedUserMock := testWeb.AllSeeder.User.UserMock.Data[0]

	url := fmt.Sprintf("%s/%s/%s", testWeb.Server.URL, userWeb.Path, selectedUserMock.Id.String)
	request, newRequestErr := http.NewRequest(http.MethodDelete, url, http.NoBody)
	if newRequestErr != nil {
		t.Fatal(newRequestErr)
	}
	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
	request.Header.Set("Authorization", "Bearer "+selectedSessionMock.AccessToken.String)
	response, doErr := http.DefaultClient.Do(request)
	if doErr != nil {
		t.Fatal(doErr)
	}

	bodyResponse := &model_response.Response[*entity.User]{}
	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(t, selectedUserMock.Id, bodyResponse.Data.Id)
	assert.Equal(t, selectedUserMock.Name, bodyResponse.Data.Name)
	assert.Equal(t, selectedUserMock.Email, bodyResponse.Data.Email)
	assert.Equal(t, selectedUserMock.Balance, bodyResponse.Data.Balance)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(bodyResponse.Data.Password.String), []byte(selectedUserMock.Password.String)))
}

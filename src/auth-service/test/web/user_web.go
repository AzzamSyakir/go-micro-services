package web

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"go-micro-services/internal/entity"
// 	model_request "go-micro-services/internal/model/request/controller"
// 	model_response "go-micro-services/internal/model/response"
// 	"net/http"
// 	"testing"

// 	"github.com/guregu/null"
// 	"github.com/stretchr/testify/assert"
// 	"golang.org/x/crypto/bcrypt"
// )

// type UserWeb struct {
// 	Test *testing.T
// 	Path string
// }

// func NewUserWeb(test *testing.T) *UserWeb {
// 	userWeb := &UserWeb{
// 		Test: test,
// 		Path: "users",
// 	}
// 	return userWeb
// }

// func (userWeb *UserWeb) Start() {
// 	userWeb.Test.Run("UserWeb_FindOneById_Succeed", userWeb.FindOneById)
// 	userWeb.Test.Run("UserWeb_FindOneByEmail_Succeed", userWeb.FindOneByEmail)
// 	userWeb.Test.Run("UserWeb_FindOneByUsername_Succeed", userWeb.FindOneByUsername)
// 	userWeb.Test.Run("UserWeb_UserPatchOneByIdRequest_Succeed", userWeb.PatchOneById)
// 	userWeb.Test.Run("UserWeb_DeleteOneById_Succeed", userWeb.DeleteOneById)
// }

// func (userWeb *UserWeb) FindOneById(t *testing.T) {
// 	t.Parallel()

// 	testWeb := GetTestWeb()
// 	testWeb.AllSeeder.Up()
// 	defer testWeb.AllSeeder.Down()

// 	selectedUserMock := testWeb.AllSeeder.User.UserMock.Data[0]

// 	url := fmt.Sprintf("%s/%s/%s", testWeb.UserServer.URL, userWeb.Path, selectedUserMock.Id.String)
// 	request, newRequestErr := http.NewRequest(http.MethodGet, url, http.NoBody)
// 	if newRequestErr != nil {
// 		t.Fatal(newRequestErr)
// 	}
// 	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
// 	request.Header.Set("Authorization", "Bearer "+selectedSessionMock.AccessToken.String)
// 	response, doErr := http.DefaultClient.Do(request)
// 	if newRequestErr != nil {
// 		t.Fatal(newRequestErr)
// 	}
// 	if doErr != nil {
// 		t.Fatal(doErr)
// 	}

// 	bodyResponse := &model_response.Response[*entity.User]{}
// 	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
// 	if decodeErr != nil {
// 		t.Fatal(decodeErr)
// 	}

// 	assert.Equal(t, http.StatusOK, response.StatusCode)
// 	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
// 	assert.Equal(t, selectedUserMock.Id, bodyResponse.Data.Id)
// 	assert.Equal(t, selectedUserMock.Name, bodyResponse.Data.Name)
// 	assert.Equal(t, selectedUserMock.Email, bodyResponse.Data.Email)
// 	assert.Equal(t, selectedUserMock.Username, bodyResponse.Data.Username)
// 	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(bodyResponse.Data.Password.String), []byte(selectedUserMock.Password.String)))
// 	assert.Equal(t, selectedUserMock.AvatarUrl, bodyResponse.Data.AvatarUrl)
// 	assert.Equal(t, selectedUserMock.Bio, bodyResponse.Data.Bio)
// }

// func (userWeb *UserWeb) FindOneByEmail(t *testing.T) {
// 	t.Parallel()

// 	testWeb := GetTestWeb()
// 	testWeb.AllSeeder.Up()
// 	defer testWeb.AllSeeder.Down()

// 	selectedUserMock := testWeb.AllSeeder.User.UserMock.Data[0]

// 	url := fmt.Sprintf("%s/%s?email=%s", testWeb.UserServer.URL, userWeb.Path, selectedUserMock.Email.String)
// 	request, newRequestErr := http.NewRequest(http.MethodGet, url, http.NoBody)
// 	if newRequestErr != nil {
// 		t.Fatal(newRequestErr)
// 	}
// 	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
// 	request.Header.Set("Authorization", "Bearer "+selectedSessionMock.AccessToken.String)
// 	response, doErr := http.DefaultClient.Do(request)
// 	if doErr != nil {
// 		t.Fatal(doErr)
// 	}

// 	bodyResponse := &model_response.Response[*entity.User]{}
// 	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
// 	if decodeErr != nil {
// 		t.Fatal(decodeErr)
// 	}

// 	assert.Equal(t, http.StatusOK, response.StatusCode)
// 	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
// 	assert.Equal(t, selectedUserMock.Id, bodyResponse.Data.Id)
// 	assert.Equal(t, selectedUserMock.Name, bodyResponse.Data.Name)
// 	assert.Equal(t, selectedUserMock.Email, bodyResponse.Data.Email)
// 	assert.Equal(t, selectedUserMock.Username, bodyResponse.Data.Username)
// 	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(bodyResponse.Data.Password.String), []byte(selectedUserMock.Password.String)))
// 	assert.Equal(t, selectedUserMock.AvatarUrl, bodyResponse.Data.AvatarUrl)
// 	assert.Equal(t, selectedUserMock.Bio, bodyResponse.Data.Bio)
// }

// func (userWeb *UserWeb) FindOneByUsername(t *testing.T) {
// 	t.Parallel()

// 	testWeb := GetTestWeb()
// 	testWeb.AllSeeder.Up()
// 	defer testWeb.AllSeeder.Down()

// 	selectedUserMock := testWeb.AllSeeder.User.UserMock.Data[0]

// 	url := fmt.Sprintf("%s/%s?username=%s", testWeb.UserServer.URL, userWeb.Path, selectedUserMock.Username.String)
// 	request, newRequestErr := http.NewRequest(http.MethodGet, url, http.NoBody)
// 	if newRequestErr != nil {
// 		t.Fatal(newRequestErr)
// 	}
// 	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
// 	request.Header.Set("Authorization", "Bearer "+selectedSessionMock.AccessToken.String)
// 	response, doErr := http.DefaultClient.Do(request)
// 	if doErr != nil {
// 		t.Fatal(doErr)
// 	}

// 	bodyResponse := &model_response.Response[*entity.User]{}
// 	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
// 	if decodeErr != nil {
// 		t.Fatal(decodeErr)
// 	}

// 	assert.Equal(t, http.StatusOK, response.StatusCode)
// 	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
// 	assert.Equal(t, selectedUserMock.Id, bodyResponse.Data.Id)
// 	assert.Equal(t, selectedUserMock.Name, bodyResponse.Data.Name)
// 	assert.Equal(t, selectedUserMock.Email, bodyResponse.Data.Email)
// 	assert.Equal(t, selectedUserMock.Username, bodyResponse.Data.Username)
// 	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(bodyResponse.Data.Password.String), []byte(selectedUserMock.Password.String)))
// 	assert.Equal(t, selectedUserMock.AvatarUrl, bodyResponse.Data.AvatarUrl)
// 	assert.Equal(t, selectedUserMock.Bio, bodyResponse.Data.Bio)
// }

// func (userWeb *UserWeb) PatchOneById(t *testing.T) {
// 	t.Parallel()

// 	testWeb := GetTestWeb()
// 	testWeb.AllSeeder.Up()
// 	defer testWeb.AllSeeder.Down()

// 	selectedUserMock := testWeb.AllSeeder.User.UserMock.Data[0]

// 	bodyRequest := &model_request.UserPatchOneByIdRequest{}
// 	bodyRequest.Name = null.NewString(selectedUserMock.Name.String+"patched", true)
// 	bodyRequest.Email = null.NewString(selectedUserMock.Email.String+"patched", true)
// 	bodyRequest.Username = null.NewString(selectedUserMock.Username.String+"patched", true)
// 	bodyRequest.Password = null.NewString(selectedUserMock.Password.String+"patched", true)
// 	bodyRequest.AvatarUrl = null.NewString(selectedUserMock.AvatarUrl.String+"patched", true)
// 	bodyRequest.Bio = null.NewString(selectedUserMock.Bio.String+"patched", true)

// 	bodyRequestJsonByte, marshalErr := json.Marshal(bodyRequest)
// 	if marshalErr != nil {
// 		t.Fatal(marshalErr)
// 	}
// 	bodyRequestBuffer := bytes.NewBuffer(bodyRequestJsonByte)

// 	url := fmt.Sprintf("%s/%s/%s", testWeb.UserServer.URL, userWeb.Path, selectedUserMock.Id.String)
// 	request, newRequestErr := http.NewRequest(http.MethodPatch, url, bodyRequestBuffer)
// 	if newRequestErr != nil {
// 		t.Fatal(newRequestErr)
// 	}
// 	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
// 	request.Header.Set("Authorization", "Bearer "+selectedSessionMock.AccessToken.String)
// 	response, doErr := http.DefaultClient.Do(request)
// 	if doErr != nil {
// 		t.Fatal(doErr)
// 	}

// 	bodyResponse := &model_response.Response[*entity.User]{}
// 	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
// 	if decodeErr != nil {
// 		t.Fatal(decodeErr)
// 	}

// 	assert.Equal(t, http.StatusOK, response.StatusCode)
// 	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
// 	assert.Equal(t, selectedUserMock.Id, bodyResponse.Data.Id)
// 	assert.Equal(t, bodyRequest.Name, bodyResponse.Data.Name)
// 	assert.Equal(t, bodyRequest.Email, bodyResponse.Data.Email)
// 	assert.Equal(t, bodyRequest.Username, bodyResponse.Data.Username)
// 	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(bodyResponse.Data.Password.String), []byte(bodyRequest.Password.String)))
// 	assert.Equal(t, bodyRequest.AvatarUrl, bodyResponse.Data.AvatarUrl)
// 	assert.Equal(t, bodyRequest.Bio, bodyResponse.Data.Bio)
// }

// func (userWeb *UserWeb) DeleteOneById(t *testing.T) {
// 	t.Parallel()

// 	testWeb := GetTestWeb()
// 	testWeb.AllSeeder.Up()
// 	defer testWeb.AllSeeder.Down()

// 	selectedUserMock := testWeb.AllSeeder.User.UserMock.Data[0]

// 	url := fmt.Sprintf("%s/%s/%s", testWeb.UserServer.URL, userWeb.Path, selectedUserMock.Id.String)
// 	request, newRequestErr := http.NewRequest(http.MethodDelete, url, http.NoBody)
// 	if newRequestErr != nil {
// 		t.Fatal(newRequestErr)
// 	}
// 	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
// 	request.Header.Set("Authorization", "Bearer "+selectedSessionMock.AccessToken.String)
// 	response, doErr := http.DefaultClient.Do(request)
// 	if doErr != nil {
// 		t.Fatal(doErr)
// 	}

// 	bodyResponse := &model_response.Response[*entity.User]{}
// 	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
// 	if decodeErr != nil {
// 		t.Fatal(decodeErr)
// 	}

// 	assert.Equal(t, http.StatusOK, response.StatusCode)
// 	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
// 	assert.Equal(t, selectedUserMock.Id, bodyResponse.Data.Id)
// 	assert.Equal(t, selectedUserMock.Name, bodyResponse.Data.Name)
// 	assert.Equal(t, selectedUserMock.Email, bodyResponse.Data.Email)
// 	assert.Equal(t, selectedUserMock.Username, bodyResponse.Data.Username)
// 	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(bodyResponse.Data.Password.String), []byte(selectedUserMock.Password.String)))
// 	assert.Equal(t, selectedUserMock.AvatarUrl, bodyResponse.Data.AvatarUrl)
// 	assert.Equal(t, selectedUserMock.Bio, bodyResponse.Data.Bio)
// }

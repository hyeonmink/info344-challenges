package handlers

import (
	"encoding/json"
	"testing"
	"time"

	"net/http"

	"net/http/httptest"

	"github.com/hyeonmink/challenges-hyeonmink/apiserver/models/users"
	"github.com/hyeonmink/challenges-hyeonmink/apiserver/sessions"
)

func TestUsersHandler(t *testing.T) {
	store := users.NewMemStore()
	testCase := &Context{
		SessionKey:   "test",
		SessionStore: sessions.NewMemStore(time.Hour),
		UserStore:    store,
	}

	user := &users.NewUser{
		Email:        "test@test.com",
		Password:     "Password",
		PasswordConf: "Password",
		UserName:     "UserName",
		FirstName:    "FirstName",
		LastName:     "LastName",
	}

	//UsersHandler
	handler := http.HandlerFunc(testCase.UsersHandler)
	resRec := httptest.NewRecorder()

	encoder := json.NewEncoder(resRec.Body)
	err := encoder.Encode(user)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/v1/users", resRec.Body)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	handler.ServeHTTP(resRec, req)
	if resRec.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusOK, resRec.Code)
	}

	contentType := resRec.Header().Get("Content-Type")
	expectedContentType := "application/json; charset=utf-8"
	if contentType != expectedContentType {
		t.Errorf("incorrect Content-Type response header: expected %s; got %s", expectedContentType, contentType)
	}

	if nil == resRec.Body || 0 == resRec.Body.Len() {
		t.Errorf("handler returned empty response body")
	}

	if _, err := store.GetByEmail(user.Email); err != nil {
		t.Errorf("Error finding email address")
	}

	if _, err := store.GetByUserName(user.UserName); err != nil {
		t.Errorf("Error finding userName")
	}

	req, err = http.NewRequest("GET", "/v1/users", resRec.Body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	handler.ServeHTTP(resRec, req)
	if resRec.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusOK, resRec.Code)
	}

	contentType = resRec.Header().Get("Content-Type")
	if contentType != expectedContentType {
		t.Errorf("incorrect Content-Type response header: expected %s; got %s", expectedContentType, contentType)
	}
	if nil == resRec.Body || 0 == resRec.Body.Len() {
		t.Errorf("handler returned empty response body")
	}

	actual := make(map[string]string)
	decoder := json.NewDecoder(resRec.Body)
	err = decoder.Decode(&actual)
	if nil != err {
		t.Errorf("error decoding returned JSON: %s", err.Error())
	}
	if actual["email"] != user.Email {
		t.Errorf("incorrect email: expected `%s` but got `%s`\n", actual["email"], user.Email)
	}
	if actual["userName"] != user.UserName {
		t.Errorf("incorrect username: expected `%s` but got `%s`\n", actual["username"], user.Email)
	}
	if actual["firstName"] != user.FirstName {
		t.Errorf("incorrect username: expected `%s` but got `%s`\n", actual["firstName"], user.Email)
	}
	if actual["lastName"] != user.LastName {
		t.Errorf("incorrect username: expected `%s` but got `%s`\n", actual["lastName"], user.Email)
	}

	req, err = http.NewRequest("SOMETHING ELSE", "/v1/users", resRec.Body)
	if err == nil {
		t.Errorf("Header should be either POST or GET")
	}

	//testing SessionsHandler
	cred := &users.Credentials{
		Email:    "test@test.com",
		Password: "Password",
	}

	handler = http.HandlerFunc(testCase.SessionsHandler)
	resRec = httptest.NewRecorder()

	encoder = json.NewEncoder(resRec.Body)
	err = encoder.Encode(cred)
	if err != nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest("SOMETING ELSE", "/v1/users", resRec.Body)
	if err == nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest("POST", "/v1/users", resRec.Body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	handler.ServeHTTP(resRec, req)
	if resRec.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusOK, resRec.Code)
	}

	contentType = resRec.Header().Get("Content-Type")
	expectedContentType = "application/json; charset=utf-8"
	if contentType != expectedContentType {
		t.Errorf("incorrect Content-Type response header: expected %s; got %s", expectedContentType, contentType)
	}

	if nil == resRec.Body || 0 == resRec.Body.Len() {
		t.Errorf("handler returned empty response body")
	}
	sid := resRec.Header().Get("Authorization")
	//testing UsersMeHanlder
	handler = http.HandlerFunc(testCase.UsersMeHanlder)
	resRec = httptest.NewRecorder()

	req, err = http.NewRequest("", "/v1/users", resRec.Body)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", sid)
	handler.ServeHTTP(resRec, req)
	if resRec.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusOK, resRec.Code)
	}

	//testing SessionsMineHandler
	handler = http.HandlerFunc(testCase.SessionsMineHandler)
	resRec = httptest.NewRecorder()

	req, err = http.NewRequest("DELETE", "/auth", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", sid)
	handler.ServeHTTP(resRec, req)
	if resRec.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusOK, resRec.Code)
	}
	if resRec.Body.String() != "user signed out" {
		t.Errorf("failed signing out")
	}
}

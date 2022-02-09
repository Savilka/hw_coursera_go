package main

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type xmlRow struct {
	Id        int    `xml:"id"`
	Guid      string `xml:"guid"`
	IsActive  bool   `xml:"isActive"`
	Balance   string `xml:"balance"`
	Picture   string `xml:"picture"`
	Age       int    `xml:"age"`
	EyeColor  string `xml:"eyeColor"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Gender    string `xml:"gender"`
	Company   string `xml:"company"`
	Email     string `xml:"email"`
	Phone     string `xml:"phone"`
	Address   string `xml:"address"`
	About     string `xml:"about"`
}

type xmlStructure struct {
	Version string   `xml:"version"`
	Row     []xmlRow `xml:"row"`
}

const pageSize = 25

func SearchServerSuccess(w http.ResponseWriter, r *http.Request) {
	dataFile, err := ioutil.ReadFile("dataset.xml")
	if err != nil {
		panic(err)
	}

	usersXml := &xmlStructure{}
	err = xml.Unmarshal(dataFile, &usersXml)
	if err != nil {
		return
	}

	var users []User

	for _, user := range usersXml.Row {
		users = append(users, User{
			Id:     user.Id,
			Name:   user.FirstName,
			Age:    user.Age,
			About:  user.About,
			Gender: user.Gender,
		})
	}

	offset, _ := strconv.Atoi(r.FormValue("offset"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))

	var startRow int
	if offset > 0 {
		startRow = offset * pageSize
	}

	endRow := startRow + limit
	users = users[startRow:endRow]

	jsonResponse, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		return
	}
}

func SearchServerBadField(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	jsonResponse, _ := json.Marshal(SearchErrorResponse{Error: "ErrorBadOrderField"})
	_, err := w.Write(jsonResponse)
	if err != nil {
		return
	}
}

func SearchServerUnknownBadRequest(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	jsonResponse, _ := json.Marshal(SearchErrorResponse{Error: "Unknown Bad Request"})
	_, err := w.Write(jsonResponse)
	if err != nil {
		return
	}
}

func SearchServerBadResponseJson(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte(""))
	if err != nil {
		return
	}
}

func SearchServerBadAccessToken(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
}

func SearchServerInternalServerError(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func SearchServerTimeoutError(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(time.Second * 2)
	w.WriteHeader(http.StatusOK)
}

func SearchServerJsonFail(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, `"err": "bad json"}`)
	if err != nil {
		return
	}
}

func SearchServerLimitFail(w http.ResponseWriter, _ *http.Request) {
	dataFile, err := ioutil.ReadFile("dataset.xml")
	if err != nil {
		panic(err)
	}

	usersXml := &xmlStructure{}
	err = xml.Unmarshal(dataFile, &usersXml)
	if err != nil {
		return
	}

	var users []User

	for _, user := range usersXml.Row {
		users = append(users, User{
			Id:     user.Id,
			Name:   user.FirstName,
			Age:    user.Age,
			About:  user.About,
			Gender: user.Gender,
		})
	}

	jsonResponse, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		return
	}
}

func TestSearchClient_FindUsers_BadField(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServerBadField))
	searchClient := &SearchClient{URL: ts.URL}
	_, err := searchClient.FindUsers(SearchRequest{})

	if err.Error() != "OrderFeld  invalid" {
		t.Error("ErrorBadOrderField is not done")
	}

	ts.Close()
}

func TestSearchClient_FindUsers_UnknownBadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServerUnknownBadRequest))
	searchClient := &SearchClient{URL: ts.URL}
	_, err := searchClient.FindUsers(SearchRequest{})

	if err.Error() != "unknown bad request error: Unknown Bad Request" {
		t.Error("unknown bad request error is not done")
	}

	ts.Close()
}

func TestSearchClient_FindUsers_BadResponseJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServerBadResponseJson))
	searchClient := &SearchClient{URL: ts.URL}
	_, err := searchClient.FindUsers(SearchRequest{})

	if err.Error() != "cant unpack error json: unexpected end of JSON input" {
		t.Error("BadResponseJson is not done")
	}

	ts.Close()
}

func TestSearchClient_FindUsers_BadAccessToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServerBadAccessToken))
	searchClient := &SearchClient{URL: ts.URL}
	_, err := searchClient.FindUsers(SearchRequest{})

	if err.Error() != "Bad AccessToken" {
		t.Error("Bad AccessToken is not done")
	}

	ts.Close()
}

func TestSearchClient_FindUsers_InternalServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServerInternalServerError))
	searchClient := &SearchClient{URL: ts.URL}
	_, err := searchClient.FindUsers(SearchRequest{})

	if err.Error() != "SearchServer fatal error" {
		t.Error("SearchServer fatal error is not done")
	}

	ts.Close()
}

func TestSearchClient_FindUsers_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServerSuccess))
	searchClient := &SearchClient{URL: ts.URL}
	_, err := searchClient.FindUsers(SearchRequest{Limit: 5, Offset: 0})

	if err != nil {
		t.Error("Doesn't work success request")
	}

	_, err = searchClient.FindUsers(SearchRequest{Limit: -1})

	if err.Error() != "limit must be > 0" {
		t.Error("limit must be > 0 is not done")
	}

	_, err = searchClient.FindUsers(SearchRequest{Offset: -1})

	if err.Error() != "offset must be > 0" {
		t.Error("offset must be > 0 is not done")
	}

	response, err := searchClient.FindUsers(SearchRequest{Limit: 26})
	if len(response.Users) != 25 {
		t.Error("limit must be < 25 is not done")
	}

	ts.Close()
}

func TestSearchClient_FindUsers_UnknownError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	searchClient := &SearchClient{URL: "bad url"}
	_, err := searchClient.FindUsers(SearchRequest{})

	if err == nil {
		t.Error("Unknown error is not done")
	}

	ts.Close()
}

func TestSearchClient_FindUsers_TimeLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServerTimeoutError))
	searchClient := &SearchClient{URL: ts.URL}
	_, err := searchClient.FindUsers(SearchRequest{})

	if err == nil {
		t.Error("Time Limit error is not done")
	}

	ts.Close()
}

func TestSearchClient_FindUsers_JsonFail(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServerJsonFail))
	searchClient := &SearchClient{URL: ts.URL}
	_, err := searchClient.FindUsers(SearchRequest{})

	if err.Error() != "cant unpack result json: invalid character ':' after top-level value" {
		t.Error("cant unpack result json error is not done")
	}

	ts.Close()
}

func TestSearchClient_FindUsers_LimitFailed(t *testing.T) {
	limit := 7
	ts := httptest.NewServer(http.HandlerFunc(SearchServerLimitFail))

	searchClient := &SearchClient{
		URL: ts.URL,
	}

	response, _ := searchClient.FindUsers(SearchRequest{Limit: limit})

	if limit == len(response.Users) {
		t.Error("Limit not true")
	}
	ts.Close()
}

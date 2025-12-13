package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

const (
	prodBaseURL  = ""
	localBaseURL = "http://localhost:8080/v1"
)

var baseURL string
var token string
var bookmarkID int
var reviewID int
var placeID = 1

type TestCase struct {
	name       string
	method     string
	path       string
	auth       bool
	body       any
	expectCode int
}

var tests = []TestCase{

	//users/register
	{"Register", "POST", "/users/register", false, map[string]string{
		"email":    "testEmail@gmail.com",
		"password": "SuperUnique!23",
		"name":     "testAcc",
	}, 204},
	{"Register_account_exists", "POST", "/users/register", false, map[string]string{
		"email":    "testEmail@gmail.com",
		"password": "SuperUnique!23",
		"name":     "testAcc",
	}, 409},
	{"Register_binding_fail", "POST", "/users/register", false, nil, 400},
	{"Register_short_name", "POST", "/users/register", false, map[string]string{
		"email":    "testEmail2@gmail.com",
		"password": "SuperUnique!23",
		"name":     "Danny",
	}, 400},
	{"Register_long_name", "POST", "/users/register", false, map[string]string{
		"email":    "testEmail2@gmail.com",
		"password": "SuperUnique!23",
		"name":     "DannyDannyDannyDannyDannyDannyDannyDanny",
	}, 400},
	{"Register_incorrect_email_format", "POST", "/users/register", false, map[string]string{
		"email":    "testEmail",
		"password": "SuperUnique!23",
		"name":     "testAcc2",
	}, 400},
	{"Register_short_password", "POST", "/users/register", false, map[string]string{
		"email":    "testEmail2@gmail.com",
		"password": "Un!23",
		"name":     "testAcc2",
	}, 400},
	{"RegisterFail_InvalidAscii", "POST", "/users/register", false, map[string]string{
		"email":    "testEmail2@gmail.com",
		"password": "SuperUnique!23👍",
		"name":     "testAcc2",
	}, 400},

	//users/login
	{"Login_binding_fail", "POST", "/users/login", false, nil, 400},
	{"Login_incorrect_name", "POST", "/users/login", false, map[string]string{
		"name":     "DannyDevit",
		"password": "SuperUnique!23",
	}, 401},
	{"Login_incorrect_password", "POST", "/users/login", false, map[string]string{
		"name":     "testAcc",
		"password": "SuperUnique!24",
	}, 401},
	{"Login_invalid-ascii", "POST", "/users/login", false, map[string]string{
		"name":     "testAcc",
		"password": "SuperUnique!23👍",
	}, 400},
	{"Login", "POST", "/users/login", false, map[string]string{
		"name":     "testAcc",
		"password": "SuperUnique!23",
	}, 200},

	//bookmarks, /bookmarks/{id}, & /bookmarks/place/{id}
	{"Get_bookmarks_dont_exist", "GET", "/bookmarks", true, nil, 200},
	{"Create_bookmark_invalid_jwt", "POST", "/bookmarks/places/1", false, nil, 401},
	{"Create_bookmark_invalid_id", "POST", "/bookmarks/places/99999999", true, nil, 500},
	{"Create_bookmark", "POST", "/bookmarks/places/1", true, nil, 204},
	{"Create_bookmark_already_exists", "POST", "/bookmarks/places/1", true, nil, 409},
	{"Get_bookmarks", "GET", "/bookmarks", true, nil, 200},
	{"Get_bookmarks_invalid_jwt", "GET", "/bookmarks", false, nil, 401},
	{"Get_my_bookmarks_user_not_found", "GET", "/bookmarks", false, nil, 401},
	{"Delete_bookmark_invalid_id", "DELETE", "/bookmarks/99", true, nil, 404},
	{"Delete_bookmark_invalid_jwt", "DELETE", "/bookmarks/1", false, nil, 401},
	{"Delete_bookmark_invalid_id", "DELETE", "/bookmarks/99999999", true, nil, 404},
	{"Delete_bookmark_success", "DELETE", "/bookmarks/0", true, nil, 204},

	//reviews/places/{id}
	{"Get_reviews_for_place", "GET", "/reviews/places/1?page=1&limit=1",
		false, nil, 200},
	{"Get_reviews_invalid_pagination", "GET", "/reviews/places/1?page=0&limit=-1", false, nil, 400},
	{"Get_reviews_invalidID", "GET", "places/0/reviews?page=1&limit=1", false, nil, 404},
	{"Get_reviews_for_place_invalid_id", "GET", "/reviews/places/99999999", false, nil, 400},
	{"Get_reviews_for_place_invalid_query", "GET", "/reviews/places/1?page=0&limit=-1", false, nil, 400},
	{"Create_review_invalid_rating", "POST", "/reviews/places/1", true, map[string]any{
		"rating": 6,
		"text":   "rating too high",
	}, 400},
	{"Create_review_long_text", "POST", "/reviews/places/1", true, map[string]any{
		"rating": 4,
		"text":   strings.Repeat("a", 201),
	}, 400},
	{"Create_review_invalid_jwt", "POST", "/reviews/places/1", false, map[string]any{
		"rating": 5, "text": "text"}, 401},
	{"Create_review_invalid_id", "POST", "/reviews/places/99999999", true, map[string]any{
		"rating": 5, "text": "text"}, 500},
	{"Create_review_long_body", "POST", "/reviews/places/1", true, map[string]any{
		"rating": 6, "text": strings.Repeat("a", 201)}, 400},
	{"Create_review_success", "POST", "/reviews/places/1", true, map[string]any{
		"rating": 5,
		"text":   "rating text",
	}, 201},
	{"Create_review_already_exists", "POST", "/reviews/places/1", true, map[string]any{
		"rating": 5, "text": "text"}, 409},

	//reviews{id}
	{"Update_review_invalid_jwt", "PATCH", "/reviews/1", false, map[string]any{
		"rating": 5, "text": "text"}, 401},
	{"Update_review_invalid_id", "PATCH", "/reviews/99999999", true, map[string]any{
		"rating": 5, "text": "text"}, 404},
	{"Update_review_long_body", "PATCH", "/reviews/1", true, map[string]any{
		"rating": 6, "text": strings.Repeat("a", 201)}, 400},
	{"Update_review", "PATCH", "/reviews/1", true, map[string]any{
		"rating": 5,
		"text":   "updated text",
	}, 200},
	{"Delete_review_invalid_jwt", "DELETE", "/reviews/1", false, nil, 401},
	{"Delete_review_invalid_id", "DELETE", "/reviews/99999999", true, nil, 404},
	{"Delete_review", "DELETE", "/reviews/1", true, nil, 204},

	{"Places_envelope", "GET", "/places/map?lat=39.95&lng=-75.165&latDelta=0.1&lngDelta=0.1", false, nil, 200},
	//recommendations
	{"Create_recommendation_invalid_jwt", "POST", "/recommendations", false, map[string]string{
		"name": "Place", "address": "123"}, 401},
	{"Create_recommendation_validation_fail", "POST", "/recommendations", true, map[string]string{
		"name": "", "address": ""}, 400},
	//users
	{"Delete_account_invalid_jwt", "DELETE", "/users", false, nil, 401},
	{"Delete_account", "DELETE", "/users", true, nil, 204},
	{"Delete_account_not_found", "DELETE", "/users", true, nil, 204},

	//reports
	{"Create_report_invalid_jwt", "POST", "/reports/places/1", false, map[string]string{
		"reason": "Test"}, 401},
	{"Create_report_invalid_id", "POST", "/reports/places/99999999", true, map[string]string{
		"reason": "Test"}, 500},
	{"Create_report_validation_fail", "POST", "/reports/places/1", true, map[string]string{
		"reason": ""}, 400},
}

func TestMain(m *testing.M) {
	env := flag.String("env", "prod", "Select environment: local or prod")
	flag.Parse()

	switch strings.ToLower(*env) {
	case "prod":
		baseURL = prodBaseURL
		fmt.Println("Running tests against production environment.")
	case "local":
		baseURL = localBaseURL
		fmt.Println("Running tests against local environment.")
	default:
		fmt.Println("Invalid environment. Use -env=local or -env=prod")
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func makeRequest(t *testing.T, tc TestCase) *http.Response {
	var bodyReader io.Reader
	if tc.body != nil {
		data, err := json.Marshal(tc.body)
		if err != nil {
			t.Fatalf("failed to marshal body for %s: %v", tc.name, err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(tc.method, baseURL+tc.path, bodyReader)
	if err != nil {
		t.Fatalf("failed to create request for %s: %v", tc.name, err)
	}

	req.Header.Set("Content-Type", "application/json")
	if tc.auth && token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed for %s: %v", tc.name, err)
	}

	return resp
}

func getToken(t *testing.T, resp *http.Response) {
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode login response: %v", err)
	}

	tkn, ok := result["token"]
	if !ok {
		t.Fatalf("no token in login response")
	}
	token = tkn.(string)
}

func getBookmarkID(t *testing.T, resp *http.Response) {
	var bookmarks []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&bookmarks); err != nil {
		t.Fatalf("failed to decode bookmarks list: %v", err)
	}

	if len(bookmarks) == 0 {
		t.Fatalf("no bookmarks returned")
	}

	idFloat, ok := bookmarks[0]["id"].(float64)
	if !ok {
		t.Fatalf("bookmark ID missing or not a number")
	}
	bookmarkID = int(idFloat)

	for i := range tests {
		if strings.HasPrefix(tests[i].name, "Delete_bookmark_success") {
			tests[i].path = "/bookmarks/" + strconv.Itoa(bookmarkID)
		}
	}
}
func getReviewID(t *testing.T, resp *http.Response) {
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode create review response: %v", err)
	}

	idVal, ok := result["id"]
	if !ok {
		t.Fatalf("no id in create review response")
	}

	idFloat, ok := idVal.(float64)
	if !ok {
		t.Fatalf("review ID is not a number")
	}

	reviewID = int(idFloat)

	for i := range tests {
		if strings.HasPrefix(tests[i].name, "Update_review") && !strings.HasSuffix(tests[i].name, "invalid_id") {
			tests[i].path = "/reviews/" + strconv.Itoa(reviewID)
		}
		if strings.HasPrefix(tests[i].name, "Delete_review") && !strings.HasSuffix(tests[i].name, "invalid_id") {
			tests[i].path = "/reviews/" + strconv.Itoa(reviewID)
		}
	}
}

func TestAPIEndpoints(t *testing.T) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := makeRequest(t, tc)
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectCode {
				t.Errorf("%s → got %d, expected %d", tc.name, resp.StatusCode, tc.expectCode)
			}

			//get token
			if tc.name == "Login" {
				getToken(t, resp)
			}
			//get bookmark id to ensure we delete same one
			if tc.name == "Get_bookmarks" && resp.StatusCode == 200 {
				getBookmarkID(t, resp)
			}
			//get review id to ensure we delete/edit same one
			if tc.name == "Create_review_success" && resp.StatusCode == 201 {
				getReviewID(t, resp)
			}
		})
	}
}

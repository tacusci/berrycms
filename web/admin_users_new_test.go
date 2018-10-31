package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/tacusci/berrycms/db"
)

const handlerRouteNewUser string = "/admin/users/new"
const handlerRouteNewRootUser string = "/admin/users/root/new"

func init() {
	//this is here so if a 500 or otherwise occurs the whole test doesn't crash because of a nil pointer
	db.Connect(db.SQLITE, "./berrycmstesting.db", "")
	db.Wipe()
	db.Setup()
}

func TestGet(t *testing.T) {
	//need this to force working directory contain /res folder
	os.Chdir("../")
	//will need to handle both new user and new root user routes
	aunh := AdminUsersNewHandler{}
	req := httptest.NewRequest("GET", handlerRouteNewRootUser, nil)
	responseRecorder := httptest.NewRecorder()

	aunh.Get(responseRecorder, req)

	resp := responseRecorder.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Test get retrieved a response which is not OK...")
	}

	//test handle get request for new user
	aunh = AdminUsersNewHandler{}
	req = httptest.NewRequest("GET", handlerRouteNewUser, nil)
	responseRecorder = httptest.NewRecorder()

	aunh.Get(responseRecorder, req)

	resp = responseRecorder.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Test get retrieved a response which is not OK...")
	}
}

func TestPost(t *testing.T) {
	aunh := AdminUsersNewHandler{}
	req := httptest.NewRequest("POST", handlerRouteNewUser, nil)
	responseRecorder := httptest.NewRecorder()

	//testing new regular user form POST submission result
	formValues := url.Values{}
	formValues["authhash"] = []string{"thisisatestpassword"}
	formValues["repeatedauthhash"] = []string{"thisisatestpassword"}
	formValues["firstname"] = []string{"IAmATest"}
	formValues["lastname"] = []string{"Person"}
	formValues["email"] = []string{"someone@place.com"}
	formValues["username"] = []string{"testuser222"}

	req.PostForm = formValues

	aunh.Post(responseRecorder, req)

	resp := responseRecorder.Result()

	if resp.StatusCode != http.StatusFound {
		t.Errorf("Test new user post didn't redirect request, STATUS CODE: %d", resp.StatusCode)
	}

	//location header will have been set on http server redirect
	if len(resp.Header["Location"]) > 0 && resp.Header["Location"][0] != "/admin/users" {
		t.Errorf("Test post new user didn't set header to redirect to correct location")
	}

	if len(resp.Header["Location"]) == 0 {
		t.Errorf("Test post new user didn't set location in header")
	}

	ut := db.UsersTable{}
	if user, err := ut.SelectByUsername(db.Conn, "testuser222"); err != nil || user.Username != "testuser222" {
		t.Errorf("Test post new user, didn't actually create the new user")
	}

	//testing new root user form POST submission result
	responseRecorder = httptest.NewRecorder()

	formValues = url.Values{}
	formValues["authhash"] = []string{"thisisatestpassword"}
	formValues["repeatedauthhash"] = []string{"thisisatestpassword"}
	formValues["firstname"] = []string{"IAmATest"}
	formValues["lastname"] = []string{"Person"}
	formValues["email"] = []string{"root@place.com"}
	formValues["username"] = []string{"rootuser222"}

	req = httptest.NewRequest("POST", handlerRouteNewRootUser, nil)

	req.PostForm = formValues

	aunh.Post(responseRecorder, req)

	resp = responseRecorder.Result()

	if resp.StatusCode != http.StatusFound {
		t.Errorf("Test new user post didn't redirect request, STATUS CODE: %d", resp.StatusCode)
	}

	//location header will have been set on http server redirect
	if len(resp.Header["Location"]) > 0 && resp.Header["Location"][0] != "/login" {
		t.Errorf("Test post new root user didn't set header to redirect to correct location")
	}

	if len(resp.Header["Location"]) == 0 {
		t.Errorf("Test post new root user didn't set location in header")
	}

	if !ut.RootUserExists() {
		t.Errorf("Test new root user post didn't create root user")
	}
}

func TestRoute(t *testing.T) {
	aunh := AdminUsersNewHandler{
		route: handlerRouteNewRootUser,
	}
	if aunh.Route() != handlerRouteNewRootUser {
		t.Errorf("Test fetched route doesn't match with set route")
	}

	aunh = AdminUsersNewHandler{
		route: handlerRouteNewUser,
	}
	if aunh.Route() != handlerRouteNewUser {
		t.Errorf("Test fetched route doesn't match with set route")
	}
}

func TestHandlesGet(t *testing.T) {
	aunh := AdminUsersNewHandler{}
	if aunh.HandlesGet() == false {
		t.Errorf("Test admin users new handler should handle get requests")
	}
}

func TestHandlesPost(t *testing.T) {
	aunh := AdminUsersNewHandler{}
	if aunh.HandlesPost() == false {
		t.Errorf("Test admin users new handler should handle post requests")
	}
}

func TestValidatePostFormPass(t *testing.T) {
	req := httptest.NewRequest("POST", handlerRouteNewRootUser, nil)

	//data set which should pass correctly
	formValues := url.Values{}
	formValues["authhash"] = []string{"thisisatestpassword"}
	formValues["repeatedauthhash"] = []string{"thisisatestpassword"}
	formValues["firstname"] = []string{"IAmATest"}
	formValues["lastname"] = []string{"Person"}
	formValues["email"] = []string{"someone@place.com"}
	formValues["username"] = []string{"testuser222"}

	req.PostForm = formValues

	//expected to pass, since all form fields should be valid
	if validated, err := validatePostForm(req); err != nil || validated == false {
		t.Errorf("Test POST should have validated, it has not: %v\n", err)
	}
}

func TestValidatePostFormFail(t *testing.T) {

	req := httptest.NewRequest("POST", "/admin/users/root/new", nil)

	formValues := url.Values{}
	formValues["authhash"] = []string{"thispassworddoesnotmatch"}
	formValues["repeatedauthhash"] = []string{"withtheotherone"}

	req.PostForm = formValues

	//expected to fail due to mismatching passwords
	if validated, err := validatePostForm(req); err == nil || validated == true {
		t.Errorf("Test POST should have failed [MISMATCH PASSWORD], it has not: %v\n", err)
	}

	formValues = url.Values{}
	formValues["authhash"] = []string{"thisisatestpassword"}
	formValues["repeatedauthhash"] = []string{"thisisatestpassword"}
	//first name deliberately blank
	formValues["firstname"] = []string{""}
	formValues["lastname"] = []string{"Lastname"}
	formValues["email"] = []string{"test@somewhere.com"}
	formValues["username"] = []string{"testuser222"}

	req.PostForm = formValues

	//expected to fail due to blank first name
	if validated, err := validatePostForm(req); err == nil || validated == true {
		t.Errorf("Test POST should have failed [FIRST NAME BLANK], is has not: %v\n", err)
	}

	formValues = url.Values{}
	formValues["authhash"] = []string{"thisisatestpassword"}
	formValues["repeatedauthhash"] = []string{"thisisatestpassword"}
	formValues["firstname"] = []string{"Firstname"}
	// last name deliberately blank
	formValues["lastname"] = []string{""}
	formValues["email"] = []string{"test@somewhere.com"}
	formValues["username"] = []string{"testuser222"}

	req.PostForm = formValues

	//expected to fail due to blank last name
	if validated, err := validatePostForm(req); err == nil || validated == true {
		t.Errorf("Test POST should have failed [LAST NAME BLANK], is has not: %v\n", err)
	}

	formValues = url.Values{}
	formValues["authhash"] = []string{"thisisatestpassword"}
	formValues["repeatedauthhash"] = []string{"thisisatestpassword"}
	formValues["firstname"] = []string{"Firstname"}
	formValues["lastname"] = []string{"Lastname"}
	//last name name deliberately blank
	formValues["email"] = []string{""}
	formValues["username"] = []string{"testuser222"}

	req.PostForm = formValues

	//expected to fail due to blank email
	if validated, err := validatePostForm(req); err == nil || validated == true {
		t.Errorf("Test POST should have failed [EMAIL BLANK], is has not: %v\n", err)
	}

	formValues = url.Values{}
	formValues["authhash"] = []string{"thisisatestpassword"}
	formValues["repeatedauthhash"] = []string{"thisisatestpassword"}
	formValues["firstname"] = []string{"Firstname"}
	formValues["lastname"] = []string{"Lastname"}
	formValues["email"] = []string{"test@somewhere.com"}
	//username deliberately blank
	formValues["username"] = []string{""}

	req.PostForm = formValues

	//expected to fail due to blank username
	if validated, err := validatePostForm(req); err == nil || validated == true {
		t.Errorf("Test POST should have failed [USERNAME BLANK], is has not: %v\n", err)
	}

	formValues = url.Values{}
	formValues["authhash"] = []string{"thisisatestpassword"}
	formValues["repeatedauthhash"] = []string{"thisisatestpassword"}
	formValues["firstname"] = []string{"Firstname"}
	formValues["lastname"] = []string{"Lastname"}
	//email deliberately incorrect format
	formValues["email"] = []string{"exampleemailplace@.com"}
	formValues["username"] = []string{"testuser222"}

	req.PostForm = formValues

	//expected to fail due to incorrect email format
	if validated, err := validatePostForm(req); err == nil || validated == true {
		t.Errorf("Test POST should have failed [EMAIL INCORRECT FORMAT], is has not: %v\n", err)
	}

	formValues = url.Values{}
	formValues["authhash"] = []string{"thisisatestpassword"}
	formValues["repeatedauthhash"] = []string{"thisisatestpassword"}
	formValues["firstname"] = []string{"Firstname"}
	formValues["lastname"] = []string{"Lastname"}
	formValues["email"] = []string{"test@somewhere.com"}
	//username deliberately incorrect format
	formValues["username"] = []string{"test*&$^user"}

	req.PostForm = formValues

	//expected to fail due to incorrect username format
	if validated, err := validatePostForm(req); err == nil || validated == true {
		t.Errorf("Test POST should have failed [USERNAME INCORRECT FORMAT], is has not: %v\n", err)
	}
}

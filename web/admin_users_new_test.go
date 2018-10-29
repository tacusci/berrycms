package web

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

type FormTest struct {
	testName      string
	vals          url.Values
	toFail        bool
	onFailMessage string
}

func TestValidatePostFormPass(t *testing.T) {
	req := httptest.NewRequest("POST", "/admin/users/root/new", nil)

	// data set which should pass correctly
	formValues := url.Values{}
	formValues["authhash"] = []string{"thisisatestpassword"}
	formValues["repeatedauthhash"] = []string{"thisisatestpassword"}
	formValues["firstname"] = []string{"IAmATest"}
	formValues["lastname"] = []string{"Person"}
	formValues["email"] = []string{"someone@place.com"}
	formValues["username"] = []string{"testuser222"}

	req.PostForm = formValues

	// expected to pass, since all form fields should be valid
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

	// expected to fail due to mismatching passwords
	if validated, err := validatePostForm(req); err == nil || validated == true {
		t.Errorf("Test POST should have failed [MISMATCH PASSWORD], it has not: %v\n", err)
	}

	formValues = url.Values{}
	formValues["authhash"] = []string{"thisisatestpassword"}
	formValues["repeatedauthhash"] = []string{"thisisatestpassword"}
	// first name deliberately blank
	formValues["firstname"] = []string{""}
	formValues["lastname"] = []string{"Lastname"}
	formValues["email"] = []string{"test@somewhere.com"}
	formValues["username"] = []string{"testuser222"}

	req.PostForm = formValues

	// expected to fail due to blank first name
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

	// expected to fail due to blank last name
	if validated, err := validatePostForm(req); err == nil || validated == true {
		t.Errorf("Test POST should have failed [LAST NAME BLANK], is has not: %v\n", err)
	}

	formValues = url.Values{}
	formValues["authhash"] = []string{"thisisatestpassword"}
	formValues["repeatedauthhash"] = []string{"thisisatestpassword"}
	formValues["firstname"] = []string{"Firstname"}
	formValues["lastname"] = []string{"Lastname"}
	// last name name deliberately blank
	formValues["email"] = []string{""}
	formValues["username"] = []string{"testuser222"}

	req.PostForm = formValues

	// expected to fail due to blank email
	if validated, err := validatePostForm(req); err == nil || validated == true {
		t.Errorf("Test POST should have failed [EMAIL BLANK], is has not: %v\n", err)
	}

	formValues = url.Values{}
	formValues["authhash"] = []string{"thisisatestpassword"}
	formValues["repeatedauthhash"] = []string{"thisisatestpassword"}
	formValues["firstname"] = []string{"Firstname"}
	formValues["lastname"] = []string{"Lastname"}
	formValues["email"] = []string{"test@somewhere.com"}
	// username deliberately blank
	formValues["username"] = []string{""}

	req.PostForm = formValues

	// expected to fail due to blank username
	if validated, err := validatePostForm(req); err == nil || validated == true {
		t.Errorf("Test POST should have failed [USERNAME BLANK], is has not: %v\n", err)
	}

	formValues = url.Values{}
	formValues["authhash"] = []string{"thisisatestpassword"}
	formValues["repeatedauthhash"] = []string{"thisisatestpassword"}
	formValues["firstname"] = []string{"Firstname"}
	formValues["lastname"] = []string{"Lastname"}
	// email deliberately incorrect format
	formValues["email"] = []string{"exampleemailplace@.com"}
	formValues["username"] = []string{"testuser222"}

	req.PostForm = formValues

	// expected to fail due to incorrect email format
	if validated, err := validatePostForm(req); err == nil || validated == true {
		t.Errorf("Test POST should have failed [EMAIL INCORRECT FORMAT], is has not: %v\n", err)
	}

	formValues = url.Values{}
	formValues["authhash"] = []string{"thisisatestpassword"}
	formValues["repeatedauthhash"] = []string{"thisisatestpassword"}
	formValues["firstname"] = []string{"Firstname"}
	formValues["lastname"] = []string{"Lastname"}
	formValues["email"] = []string{"test@somewhere.com"}
	// username deliberately incorrect format
	formValues["username"] = []string{"test*&$^user"}

	req.PostForm = formValues

	// expected to fail due to incorrect username format
	if validated, err := validatePostForm(req); err == nil || validated == true {
		t.Errorf("Test POST should have failed [USERNAME INCORRECT FORMAT], is has not: %v\n", err)
	}
}

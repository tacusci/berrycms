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

var failTestCases = []FormTest{
	FormTest{
		testName: "[Passwords Mismatch]",
		vals: url.Values{
			"authhash":         []string{"thispassworddoesnotmatch"},
			"repeatedauthhash": []string{"withtheotherone"},
		},
		toFail:        true,
		onFailMessage: "Expected %s test to fail %b, %v\n",
	},
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

	t.Logf("%d\n", len(failTestCases))

	req := httptest.NewRequest("POST", "/admin/users/root/new", nil)

	for _, testCase := range failTestCases {
		req.PostForm = testCase.vals
		if validated, err := validatePostForm(req); err == nil || validated == true {
			t.Errorf(testCase.onFailMessage, testCase.testName, testCase.toFail, err)
		}
	}

	// formValues := url.Values{}
	// formValues["authhash"] = []string{"thispassworddoesnotmatch"}
	// formValues["repeatedauthhash"] = []string{"withtheotherone"}

	// req.PostForm = formValues

	// // expected to fail due to mismatching passwords
	// if validated, err := validatePostForm(req); err == nil || validated == true {
	// 	t.Errorf("Test POST should have failed validation, it has not: %v\n", err)
	// }

	// formValues = url.Values{}
	// formValues["authhash"] = []string{"thisisatestpassword"}
	// formValues["repeatedauthhash"] = []string{"thisisatestpassword"}
	// formValues["firstname"] = []string{""}
	// formValues["lastname"] = []string{"Lastname"}
	// formValues["email"] = []string{""}
	// formValues["firstname"] = []string{""}

}

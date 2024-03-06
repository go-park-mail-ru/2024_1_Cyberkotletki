package user

import (
	"testing"
)

func TestValidatePassword(t *testing.T) {
	var userObj = NewUserEmpty()
	var password = ""
	err := userObj.ValidatePassword(password)

	if err == nil {
		t.Errorf("Expected len(password) > 0 , got %d", len(password))
	}
	password = "123456789101234567891012345678910123456789101234567891012345678910" +
		"12345678910123456789101234567891012345678910123456789101234567891012345678910" +
		"12345678910123456789101234567891012345678910123456789101234567891012345678910" +
		"12345678910123456789101234567891012345678910123456789101234567891012345678910" +
		"12345678910123456789101234567891012345678910123456789101234567891012345678910"
	err = userObj.ValidatePassword(password)

	if err == nil {
		t.Errorf("Expected len(password) < 72 , got %d", len(password))
	}
	password = "12345678910"
	err = userObj.ValidatePassword(password)

	if err != nil {
		t.Errorf("Expected 0 < len(password) < 72 , got %d", len(password))
	}

}

func TestValidateEmail(t *testing.T) {
	var userObj = NewUserEmpty()
	var email = "email@gmailcom"
	err := userObj.ValidateEmail(email)

	if err != nil {
		t.Errorf("Expected email with @ , got %s", email)
	}
	email = "emailgmailcom"
	err = userObj.ValidateEmail(email)
	if err == nil {
		t.Errorf("Expected email without @ , got %s", email)
	}
}

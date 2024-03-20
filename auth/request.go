package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

type SignupForm struct {
	Name     string
	Email    string
	Password string
	Confirm  string
}

func (f *SignupForm) bindAndValidate(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	confirmPassword := r.PostForm.Get("confirm-password")

	var verr error

	if err := validation.Validate(name,
		validation.Required,
	); err != nil {
		err = fmt.Errorf("name %s", err.Error())
		verr = errors.Join(err)
	}

	if err := validation.Validate(email,
		validation.Required,
		is.EmailFormat,
	); err != nil {
		err = fmt.Errorf("email %s", err.Error())
		verr = errors.Join(err)
	}

	if err := validatePassword(password); err != nil {
		verr = errors.Join(err)
	}

	if password != confirmPassword {
		err := errors.New("confirm password does not match the password")
		verr = errors.Join(err)
	}

	if verr != nil {
		return verr
	}

	f.Name = name
	f.Email = email
	f.Password = password

	return nil
}

type SigninForm struct {
	Email    string
	Password string
}

func (f *SigninForm) bindAndValidate(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	password := r.PostForm.Get("password")
	email := r.PostForm.Get("email")

	var verr error
	if err := validation.Validate(password, validation.Required); err != nil {
		err = fmt.Errorf("password %s", err.Error())
		verr = errors.Join(err)
	}

	if err := validation.Validate(email,
		validation.Required,
		is.EmailFormat,
	); err != nil {
		err = fmt.Errorf("email %s", err.Error())
		verr = errors.Join(verr, err)
	}
	if verr != nil {
		return verr
	}

	f.Password = password
	f.Email = email

	return nil
}

type SignupConfirmForm struct {
	OTP   string
	Email string
}

func (f *SignupConfirmForm) bindAndValidate(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	otp := r.PostForm.Get("otp")
	email := r.PostForm.Get("email")

	var validationErrors error
	if err := validation.Validate(otp,
		validation.Required,
		validation.Length(6, 6),
		is.Digit,
	); err != nil {
		err = fmt.Errorf("otp %s", err.Error())
		validationErrors = errors.Join(err)
	}

	if err := validation.Validate(email,
		validation.Required,
		is.EmailFormat,
	); err != nil {
		err = fmt.Errorf("email %s", err.Error())
		validationErrors = errors.Join(err)
	}
	if validationErrors != nil {
		return validationErrors
	}

	f.OTP = otp
	f.Email = email

	return nil
}

func isHXRequest(r *http.Request) bool {
	if r.Header.Get("Hx-Request") == "true" {
		return true
	}

	return false
}

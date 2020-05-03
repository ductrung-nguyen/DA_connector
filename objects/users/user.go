package users

import (
	"da_connector/context"
	"da_connector/objects/base"
	"fmt"
	"net/url"
	"strconv"

	"github.com/golang/glog"
)

const (
	CMD_API_EMAIL_AUTH        = "CMD_API_EMAIL_AUTH"
	CMD_CHANGE_EMAIL_PASSWORD = "CMD_CHANGE_EMAIL_PASSWORD"
	CMD_API_POP               = "CMD_API_POP"
	CMD_EMAIL_ACCOUNT_QUOTA   = "CMD_EMAIL_ACCOUNT_QUOTA"
)

type User struct {
	base.BaseObject
	Context context.UserContext
}

// EmailUsage contains information of the usage in bytes
type EmailUsage struct {
	Imap               int64
	Inbox              int64
	Spam               int64
	Total              int64
	Webmail            int64
	Quota              int64
	LastPasswordChange int64
}

func (u User) GetUsername() string {
	return ""
}

func handleError(m url.Values) (hasError bool, errorMessage string) {
	message := ""
	if len(m["text"]) > 0 {
		message = m["text"][0]
	}

	detail := ""
	if len(m["details"]) > 0 {
		detail = ". " + m["details"][0]
	}

	if v, ok := m["error"]; ok {
		if len(v) == 0 {
			return true, message + detail
		}

		if v[0] == "1" {
			return true, message + detail
		}

		return false, message + detail
	}
	return false, message + detail
}

// CheckEmailLogin check if the email and password are valid
//
// ctx := context.BaseContext{BaseURL: "http://example.com:2222/"}
// ctx.SetLogin("account_username", "account_password") // // this account must own domain example.com
// user := users.User{Context: context.UserContext{BaseContext: ctx}}
// ok, err := user.CheckEmailLogin("testweb@example.com", "email_password")
func (u User) CheckEmailLogin(email, password string) (ok bool, message string, err error) {
	qry := map[string]string{
		"email":  email,
		"passwd": password,
	}
	resp, err := u.Context.RawRequest("POST", CMD_API_EMAIL_AUTH, qry, map[string]string{}, true)
	if err != nil {
		glog.Errorf("Error when checking email login. Err: %s\n")
		return false, "Cannot send request to login", err
	}
	m, err := url.ParseQuery(resp)
	if err != nil {
		return false, "Invalid DA account username or password", err
	}

	if hasError, message := handleError(m); hasError {
		return false, message, nil
	} else {
		return true, message, nil
	}
}

// ChangeEmailPassword changes email password
func (u User) ChangeEmailPassword(email, oldPassword, newPassword string) (ok bool, message string, err error) {
	qry := map[string]string{
		"email":       email,
		"oldpassword": oldPassword,
		"password1":   newPassword,
		"password2":   newPassword,
		"api":         "1",
	}
	resp, err := u.Context.RawRequest("POST", CMD_CHANGE_EMAIL_PASSWORD, qry, map[string]string{}, false)
	if err != nil {
		glog.Errorf("Error when changing email password. Err: %s\n")
		return false, "Cannot send request to change email password", err
	}
	m, err := url.ParseQuery(resp)
	if err != nil {
		return false, "Invalid DA account username or password", err
	}
	if hasError, message := handleError(m); hasError {
		return false, message, nil
	} else {
		return true, message, nil
	}
}

// ListAllEmailAccounts list all email accounts
func (u User) ListAllEmailAccounts(domain string) (emails []string, message string, err error) {
	qry := map[string]string{"action": "list", "domain": domain}
	resp, err := u.Context.RawRequest("POST", CMD_API_POP, qry, map[string]string{}, true)
	if err != nil {
		glog.Errorf("Error when listing email accounts. Err: %s\n")
		return nil, "Cannot send request to list email accounts", err
	}
	m, err := url.ParseQuery(resp)
	glog.Infoln(m)
	if err != nil {
		return nil, "Invalid DA account username or password", err
	}
	if hasError, message := handleError(m); hasError {
		return nil, message, nil
	} else {
		return m["list[]"], message, nil
	}
}

// CreateEmailAccount create new email account
// quota : integer in MBs. (0 = unlimited, >=1 = number MBs)
// sendLimit: 0 = unlimited, "" to use the system default
func (u User) CreateEmailAccount(domain, emailUser, password string, quota int, sendLimit string) (ok bool, message string, err error) {
	qry := map[string]string{
		"action":  "create",
		"domain":  domain,
		"user":    emailUser,
		"passwd":  password,
		"passwd2": password,
		"quota":   fmt.Sprintf("%d", quota),
		"limit":   sendLimit,
	}
	resp, err := u.Context.RawRequest("POST", CMD_API_POP, qry, map[string]string{}, true)
	if err != nil {
		glog.Errorf("Error when creating email account. Err: %s\n")
		return false, "Cannot send request to create email account", err
	}
	m, err := url.ParseQuery(resp)
	if err != nil {
		return false, "Invalid DA account username or password", err
	}
	if hasError, message := handleError(m); hasError {
		return false, message, nil
	} else {
		return true, message, nil
	}
}

// CreateEmailAccount create new email account
// quota : integer in MBs. (0 = unlimited, >=1 = number MBs)
// sendLimit: 0 = unlimited, "" to use the system default
func (u User) DeleteEmailAccount(domain, emailUser string) (ok bool, message string, err error) {
	qry := map[string]string{
		"action": "delete",
		"domain": domain,
		"user":   emailUser,
	}
	resp, err := u.Context.RawRequest("POST", CMD_API_POP, qry, map[string]string{}, true)
	if err != nil {
		glog.Errorf("Error when deleting email account. Err: %s\n")
		return false, "Cannot send request to delete email account", err
	}
	m, err := url.ParseQuery(resp)
	if err != nil {
		return false, "Invalid DA account username or password", err
	}
	if hasError, message := handleError(m); hasError {
		return false, message, nil
	} else {
		return true, message, nil
	}
}

// GetEmailUsageInfo returns quota information of an account (can be DA account or email account)
func (u User) GetEmailUsageInfo(domain, user, password string) (usage EmailUsage, message string, err error) {
	qry := map[string]string{
		"domain":   domain,
		"user":     user,
		"password": password,
		"api":      "1",
	}
	resp, err := u.Context.RawRequest("POST", CMD_EMAIL_ACCOUNT_QUOTA, qry, map[string]string{}, false)
	if err != nil {
		glog.Errorf("Error when getting quota information. Err: %s\n")
		return EmailUsage{}, "Cannot send request to get quota information", err
	}
	m, err := url.ParseQuery(resp)
	if err != nil {
		return EmailUsage{}, "Invalid DA account username or password", err
	}
	if hasError, message := handleError(m); hasError {
		return EmailUsage{}, message, nil
	} else {
		imap, _ := strconv.ParseInt(m["imap_bytes"][0], 10, 64)
		inbox, _ := strconv.ParseInt(m["inbox_bytes"][0], 10, 64)
		lastPasswordChange, _ := strconv.ParseInt(m["last_password_change"][0], 10, 64)
		spam, _ := strconv.ParseInt(m["spam_bytes"][0], 10, 64)
		totalBytes, _ := strconv.ParseInt(m["total_bytes"][0], 10, 64)
		webmailBytes, _ := strconv.ParseInt(m["webmail_bytes"][0], 10, 64)

		return EmailUsage{
			Imap:               imap,
			Inbox:              inbox,
			LastPasswordChange: lastPasswordChange,
			Spam:               spam,
			Total:              totalBytes,
			Webmail:            webmailBytes,
		}, message, nil
	}
}

// GetEmailQuotaInfo returns quota information of an account (can be DA account or email account)
func (u User) GetEmailQuotaInfo(domain, user string) (usage EmailUsage, message string, err error) {
	qry := map[string]string{
		"domain": domain,
		"user":   user,
		"type":   "quota",
		"api":    "1",
	}
	resp, err := u.Context.RawRequest("POST", CMD_API_POP, qry, map[string]string{}, true)
	if err != nil {
		glog.Errorf("Error when getting quota information. Err: %s\n")
		return EmailUsage{}, "Cannot send request to get quota information", err
	}
	m, err := url.ParseQuery(resp)
	if err != nil {
		return EmailUsage{}, "Invalid DA account username or password", err
	}
	if hasError, message := handleError(m); hasError {
		return EmailUsage{}, message, nil
	} else {
		imap, _ := strconv.ParseInt(m["imap_bytes"][0], 10, 64)
		inbox, _ := strconv.ParseInt(m["inbox_bytes"][0], 10, 64)
		spam, _ := strconv.ParseInt(m["spam_bytes"][0], 10, 64)
		totalBytes, _ := strconv.ParseInt(m["total_bytes"][0], 10, 64)
		webmailBytes, _ := strconv.ParseInt(m["webmail_bytes"][0], 10, 64)
		quota, _ := strconv.ParseInt(m["quota"][0], 10, 64)

		return EmailUsage{
			Imap:    imap,
			Inbox:   inbox,
			Spam:    spam,
			Total:   totalBytes,
			Webmail: webmailBytes,
			Quota:   quota,
		}, message, nil
	}
}

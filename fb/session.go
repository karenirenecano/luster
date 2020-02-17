// Package fb contain functions for access to Facebook.
//
// First you need to use Login function to get Session.
// Then you can use Get or Post functions to make requests to Facebook.
package fb

import (
	"fmt"
	"net/http"
	"regexp"

	r "github.com/levigross/grequests"
)

// Session holds all data needed for making authenticated requests to Facebook
type Session struct {
	UID     string         // user id
	Secret  string         // fb_dtsg attribute needed for Post requests
	Cookies []*http.Cookie // lig of cookies
}

// Login to Facebook using email address and password
func Login(email string, password string) (*Session, error) {
	// post login credentials
	response, err := r.Post("https://m.facebook.com/login.php", &r.RequestOptions{
		Data:          map[string]string{"email": email, "pass": password},
		RedirectLimit: -1,
	})
	if err != nil {
		return nil, fmt.Errorf("Login request failed for email %s", email, err)
	}

	// get cookies from login response
	cookies := response.RawResponse.Cookies()

	// find user id in cookies
	uid := ""
	for _, c := range cookies {
		if c.Name == "c_user" {
			uid = c.Value
		}
	}
	if uid == "" {
		return nil, fmt.Errorf("User id was not found in cookies", err)
	}

	// go to user home page to find 'fb_dtsg' "secret" needed for next post requests
	response, err = r.Get("https://m.facebook.com/"+uid, &r.RequestOptions{
		Cookies: response.RawResponse.Cookies(),
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot get user page with uid '%s'", uid, err)
	}

	// extract fb_dtsg from hidden input
	re, _ := regexp.Compile("name=\"fb_dtsg\" value=\"([^\"]+)\"")
	m := re.FindStringSubmatch(response.String())
	if len(m) == 0 {
		return nil, fmt.Errorf("Cannot find 'fb_dtsg' value")
	}
	usecret := m[1]

	// create session
	return &Session{
		UID:     uid,
		Secret:  usecret,
		Cookies: cookies,
	}, nil
}

// Get makes GET request to given path under https://www.facebook.com/ with given url params
func (s *Session) Get(path string, params map[string]string) (string, error) {
	url := "https://www.facebook.com/" + path
	response, err := r.Get(url, &r.RequestOptions{
		Cookies: s.Cookies,
		Params:  params,
	})
	if err != nil {
		return "", fmt.Errorf("Get '%s' failed", url, err)
	}
	return response.String(), nil
}

// Post sends POST form request to given path under https://www.facebook.com/ with given url params
func (s *Session) Post(path string, params map[string]string) (string, error) {
	url := "https://www.facebook.com/" + path
	ro := &r.RequestOptions{
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Cookies: s.Cookies,
		Params:  params,
		Data: map[string]string{
			"__user":  s.UID,
			"fb_dtsg": s.Secret,
			"__a":     "",
		},
	}
	response, err := r.Post(url, ro)
	if err != nil {
		return "", fmt.Errorf("Post to '%s' failed", url, err)
	}
	return response.String(), nil
}

package context

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/golang/glog"
)

type BaseContext struct {
	Username string
	Password string
	BaseURL  string
}

// InvokeApiGet
func (b *BaseContext) SetLogin(username string, password string) {
	b.Username = username
	b.Password = password
}

// RawRequest sends raw request
// needDAAcc: if we need the DirectAdmin (manager) account
func (b *BaseContext) RawRequest(method string, command string, payload map[string]string, header map[string]string, needDAAcc bool) (response string, err error) {
	if b.BaseURL[len(b.BaseURL)-1] != '/' {
		b.BaseURL += "/"
	}

	str := ""
	for k, v := range payload {
		if str != "" {
			str = str + "&"
		}
		str += k + "=" + v
	}
	glog.Infoln(str)

	payloadInBytes := strings.NewReader(str)
	req, err := http.NewRequest(strings.ToUpper(method), b.BaseURL+command, payloadInBytes)

	if err != nil {
		glog.Errorf("Error when creating request. Error: %s\n", err)
		return "", err
	}

	if needDAAcc {
		req.SetBasicAuth(b.Username, b.Password)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for k, v := range header {
		req.Header.Set(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	defer res.Body.Close()

	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(res.Body)
		if err != nil {
			return "", nil
		}
		// defer reader.Close()
	default:
		reader = res.Body
	}

	contentBytes, err := ioutil.ReadAll(reader)
	defer reader.Close()
	if err != nil {
		return "", err
	}

	return string(contentBytes), nil
}

// check my ip address
package ftunnel

import (
	"io/ioutil"
	"net/http"
)

func myip() (string, error) {

	response, err := http.Get("http://ifconfig.me/ip")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

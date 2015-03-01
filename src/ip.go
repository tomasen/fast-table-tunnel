// Check ip address of this node
package ftunnel

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// check my ip address
func ip() (s string) {
	for {
		s = real_check_my_ip()
		if len(s) > 0 {
			return
		}
		log.Println("N(ip.0):", "retrying")
		time.Sleep(3 * time.Second)
		// TODO: query with other nodes or ipinfo.io/ip
	}
	return
}

func real_check_my_ip() (s string) {
	// TODO: switch to other api if necessary
	response, err := http.Get("http://ifconfig.me/ip")
	if err != nil {
		log.Println("E(real_check_my_ip.0):", err)
		return
	}
	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("E(real_check_my_ip.1):", err)
		return
	}

	return string(b)
}

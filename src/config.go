package ftunnel

import (
	"bitbucket.org/kardianos/osext"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type Config struct {
	last_checksum []byte
	once          sync.Once
}

func (rc *Config) Check(b []byte) {
	type Binary struct {
		BinaryUrl      string
		BinaryCheckSum []byte
	}
	var c Binary
	err := json.Unmarshal(b, &c)
	if err != nil {
		log.Println("E(config.check.Unmarshal): ", err)
		return
	}

	execfile, err := osext.Executable()
	if err != nil {
		log.Println("E(config.check.Executable): ", err)
		return
	}

	// TODO: only checksum when binary updated
	f, err := os.OpenFile(execfile, os.O_RDWR, 0777)
	if err != nil {
		log.Println("E(config.check.Open): ", err)
		return
	}
	defer f.Close()
	h1 := md5.New()
	io.Copy(h1, f)
	checksum := h1.Sum(nil)

	if len(c.BinaryCheckSum) > 0 && bytes.Equal(checksum, c.BinaryCheckSum) {
		return
	}

	response, err := http.Get(c.BinaryUrl)
	if err != nil {
		log.Println("E(config.check.Get): ", err)
		return
	} else {
		defer response.Body.Close()
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println("E(config.check.ReadAll): ", err)
			return
		}

		h2 := md5.New()
		h2.Write(b)
		bin_chksum := h2.Sum(nil)
		if bytes.Equal(checksum, bin_chksum) {
			return
		}

		if len(c.BinaryCheckSum) > 0 && !bytes.Equal(checksum, bin_chksum) {
			log.Println("E(config.check.BinaryCheckSum): not equal")
			return
		}

		f.Seek(0, 0)
		f.Write(b)

		// TODO: restart current process
	}
}

func (rc *Config) Load(uri string, c chan []byte) error {
	u, err := url.Parse(uri)
	if err != nil {
		// assume config is a file
		b, err := ioutil.ReadFile(uri)
		if err != nil {
			return err
		}
		c <- b
		return nil
	}

	// if config file is url
	// start config update daemon
	var b []byte
	switch u.Scheme {
	case "http", "https":
		response, err := http.Get(uri)
		if err != nil {
			return err
		} else {
			defer response.Body.Close()
			b, err = ioutil.ReadAll(response.Body)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("Unsupported url scheme while fetching config")
	}

	h := md5.New()
	h.Write(b)
	checksum := h.Sum(nil)

	if !bytes.Equal(checksum, rc.last_checksum) {
		rc.last_checksum = checksum
		rc.Check(b)
		c <- b
		go rc.once.Do(func() {
			t := time.Tick(1 * time.Minute)
			for _ = range t {
				err := rc.Load(uri, c)
				if err != nil {
					log.Println("E(config.load.tick):", err)
				}
			}
		})
	}
	return nil
}

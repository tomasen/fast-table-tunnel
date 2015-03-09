// check config update, self-update
// start/stop core control center and services
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

type Supervisor struct {
	last_checksum []byte
	once          sync.Once
}

func (sp *Supervisor) SelfUpdate() {
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

	if len(_core.BinaryCheckSum) > 0 && bytes.Equal(checksum, _core.BinaryCheckSum) {
		return
	}

	response, err := http.Get(_core.BinaryUrl)
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

		if len(_core.BinaryCheckSum) > 0 && !bytes.Equal(checksum, bin_chksum) {
			log.Println("E(config.check.BinaryCheckSum): not equal")
			return
		}

		f.Seek(0, 0)
		f.Write(b)

		// restart current process
		wd, err := os.Getwd()
		if nil != err {
			return
		}

		_, err = os.StartProcess(execfile, os.Args, &os.ProcAttr{
			Dir:   wd,
			Env:   os.Environ(),
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		})

		os.Exit(0)
	}
}

func (sp *Supervisor) Load(uri string) error {
	var b []byte

	u, err := url.Parse(uri)
	if err != nil {
		// config is a file
		b, err = ioutil.ReadFile(uri)
		if err != nil {
			return err
		}
	} else {
		// config file is an url
		// start config update daemon
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
	}

	h := md5.New()
	h.Write(b)
	checksum := h.Sum(nil)

	if !bytes.Equal(checksum, sp.last_checksum) {
		sp.last_checksum = checksum

		_core.Stop()

		err = json.Unmarshal(b, &_core)
		if err != nil {
			log.Println("E(config.load.Unmarshal2): ", err)
			return err
		}

		sp.SelfUpdate()

		_core.Start()

		go sp.once.Do(func() {
			t := time.Tick(1 * time.Minute)
			for _ = range t {
				err := sp.Load(uri)
				if err != nil {
					log.Println("E(config.load.tick):", err)
				}
			}
		})
	}
	return nil
}

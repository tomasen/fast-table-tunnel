package ftunnel

import ()

func Encrypt(buff []byte) {
	for k, v := range buff {
		buff[k] = ^v
	}
}

func Decrypt(buff []byte) {
	for k, v := range buff {
		buff[k] = ^v
	}
}

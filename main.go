package main

import (
	"flag"
	"github.com/klauspost/compress/zstd"
	"log"
	"net"
	"sync"
)

const (
	BUFFER_MAXSIZE = 64 * 1024
)


func main() {
	addrListen := flag.String("s", "", "listen to ip:port")
	addrConn := flag.String("c", "", "connect to ip:port")
	role := flag.String("role", "flip", "as: server, client")

	// parse arguments
	flag.Parse()

	if len(*addrListen) <= 0 || len(*addrConn) <= 0 {
		flag.PrintDefaults()
		return
	}

	Serve(*addrListen, *addrConn, *role)
}

func Serve(addrListen, addrConn string, role string) {
	var comp Scrambler
	switch role {
	case "server":
		comp = NewCompresser()
	case "client":
		comp = &ClientCompressor{NewCompresser()}
	default:
		comp = &Flipper{}
	}

	l, err := net.Listen("tcp", addrListen)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		connFromClient, err := l.Accept()
		if err != nil {
			log.Println("accept error:", err)
			break
		}

		go func() {
			connToServer, err := net.Dial("tcp", addrConn)
			if err != nil {
				// handle error
				log.Println("connect error:", err)
			}
			var wg sync.WaitGroup

			go func() {
				wg.Add(1)
				defer wg.Done()

				defer connToServer.Close()
				defer connFromClient.Close()

				// from server to client
				buff := make([]byte, BUFFER_MAXSIZE)

				for {
					n, err := connToServer.Read(buff)
					if err != nil {
						log.Println("read from server error:", err)
						return
					}

					// encode if this is server
					dbuf, err := comp.Encode(buff[:n])
					// decode if this is client
					if err != nil {
						log.Println("scramble error from server:", err)
						return
					}

					_, err = connFromClient.Write(dbuf)
					if err != nil {
						log.Println("write to client error:", err)
						return
					}
				}
			}()

			go func() {
				wg.Add(1)
				defer wg.Done()

				defer connToServer.Close()
				defer connFromClient.Close()

				// from client to server
				buff := make([]byte, BUFFER_MAXSIZE)

				for {
					n, err := connFromClient.Read(buff)
					if err != nil {
						log.Println("read from client error", err)
						return
					}

					// encode if this is client

					// decode if this is sever
					dbuf, err := comp.Decode(buff[:n])
					if err != nil {
						log.Println("scramble error to client:", err)
						return
					}

					_, err = connToServer.Write(dbuf)
					if err != nil {
						log.Println("write to client error:", err)
						return
					}
				}
			}()

			wg.Wait()
		}()
	}
}

type Scrambler interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}

type Flipper struct {}

func (f Flipper) Encode(buff []byte) ([]byte, error) {
	Flip(buff)
	return buff, nil
}

func (f Flipper) Decode(buff []byte) ([]byte, error) {
	Flip(buff)
	return buff, nil
}

type Compressor struct {
	enc *zstd.Encoder
	dec *zstd.Decoder
	dstBuf []byte
}

/*
 should use
 // Compress input to output.
	func Compress(in io.Reader, out io.Writer) error {
		enc, err := zstd.NewWriter(out)
		if err != nil {
			return err
		}
		_, err = io.Copy(enc, in)
		if err != nil {
			enc.Close()
			return err
		}
		return enc.Close()
	}
 */
func NewCompresser() *Compressor {
	dec, _ := zstd.NewReader(nil)
	enc, _ := zstd.NewWriter(nil)
	return &Compressor{enc, dec, make([]byte, BUFFER_MAXSIZE)}
}

func (e *Compressor) Encode(buff []byte) ([]byte, error) {
	return e.enc.EncodeAll(buff, e.dstBuf), nil
}

func (e *Compressor) Decode(buff []byte) ([]byte, error) {
	return e.dec.DecodeAll(buff, e.dstBuf)
}

type ClientCompressor struct {
	*Compressor
}

func (e *ClientCompressor) Encode(buff []byte) ([]byte, error) {
	return e.dec.DecodeAll(buff, e.dstBuf)
}

func (e *ClientCompressor) Decode(buff []byte) ([]byte, error) {
	return e.enc.EncodeAll(buff, e.dstBuf), nil
}

func Flip(buff []byte) {
	for k, v := range buff {
		buff[k] = ^v
	}
}

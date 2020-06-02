package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	goHTTP "net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type httpParameter struct {
	Address         string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	MaxHeaderBytes  int
	Handler         goHTTP.Handler //nil for gorilla mux
	IsTLS           bool           // is TLS
	CertificateFile string         //if TLS , certificate file
	KeyFile         string         //if TLS , key file
}
type HTTPServer struct {
	server *goHTTP.Server
	router *mux.Router
	param  *httpParameter
}

//newHTTP make a new SocketManager
func newHTTP(params *httpParameter) (s *HTTPServer, err error) {

	defer func() {
		if r := recover(); r != nil {

			err = errors.New("httpAPI NewService panic")
		}
	}()

	if params.Port == "" {
		return nil, errors.New("port error ")
	}

	r := mux.NewRouter()

	sv := &HTTPServer{
		router: r, //router
		server: &goHTTP.Server{
			Handler:        r, //router
			Addr:           params.Address + ":" + params.Port,
			ReadTimeout:    params.ReadTimeout,
			WriteTimeout:   params.WriteTimeout,
			IdleTimeout:    params.IdleTimeout,
			MaxHeaderBytes: params.MaxHeaderBytes,
		},
		param: params,
	}

	return sv, err
}

//Start starts http server
func (s *HTTPServer) Start() (resErr error) {

	go func() {
		defer func() {
			if r := recover(); r != nil {
				errMsg := fmt.Sprintf("httpServer start panic=%v", r)
				log.Printf(errMsg)
				resErr = fmt.Errorf(errMsg)
			}
		}()

		if s.param.IsTLS {
			if err := s.server.ListenAndServeTLS(s.param.CertificateFile,s.param.KeyFile); err != nil {
				panic(fmt.Sprintf("httpServer listenAndServer TLS error =%s", err.Error()))
			}

		} else {

			if err := s.server.ListenAndServe(); err != nil {
				panic(fmt.Sprintf("httpServer listenAndServer error =%s", err.Error()))
			}
		}
	}()

	return resErr
}

//shutdown stops http server
func (s *HTTPServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		//todo
	}
}

func (s *HTTPServer) GetRouter() *mux.Router {
	return s.router
}

func (s *HTTPServer) GetRealAddr(r *goHTTP.Request) string {

	remoteIP := ""
	// the default is the originating ip. but we try to find better options because this is almost
	// never the right IP
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = parts[0]
	}
	// If we have a forwarded-for header, take the address from there
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip.String()
		}
		// parse X-Real-Ip header
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip.String()
		}
	}

	return remoteIP
}

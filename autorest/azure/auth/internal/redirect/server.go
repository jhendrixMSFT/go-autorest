package redirect

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
)

const okPage = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <meta http-equiv="refresh" content="10;url=https://docs.microsoft.com/en-us/cli/azure/">
    <title>Login successfully</title>
</head>
<body>
    <h4>You have logged into Microsoft Azure!</h4>
    <p>You can close this window, or we will redirect you to the <a href="https://docs.microsoft.com/en-us/cli/azure/">Azure CLI documents</a> in 10 seconds.</p>
</body>
</html>
`

const failPage = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <title>Login failed</title>
</head>
<body>
    <h4>Some failures occurred during the authentication</h4>
    <p>You can log an issue at <a href="https://github.com/azure/azure-cli/issues">Azure CLI GitHub Repository</a> and we will assist you in resolving it.</p>
</body>
</html>
`

type Server interface {
	Start(reqState string) string
	Stop()
	Wait()
	AuthorizationCode() (string, error)
}

type server struct {
	wg   *sync.WaitGroup
	s    *http.Server
	code string
	err  error
}

func NewServer() Server {
	rs := &server{
		wg: &sync.WaitGroup{},
		s:  &http.Server{},
	}
	return rs
}

func (s *server) Start(reqState string) string {
	port := rand.Intn(600) + 8400
	s.s.Addr = fmt.Sprintf(":%d", port)
	s.s.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer s.wg.Done()
		qp := r.URL.Query()
		if respState, ok := qp["state"]; !ok {
			s.err = errors.New("missing OAuth state")
			return
		} else if respState[0] != reqState {
			s.err = errors.New("mismatched OAuth state")
			return
		}
		if err, ok := qp["error"]; ok {
			w.Write([]byte(failPage))
			s.err = fmt.Errorf("authentication error: %s; description: %s", err[0], qp.Get("error_description"))
			return
		}
		if code, ok := qp["code"]; ok {
			w.Write([]byte(okPage))
			s.code = code[0]
		} else {
			s.err = errors.New("authorization code missing in query string")
		}
	})
	s.wg.Add(1)
	go s.s.ListenAndServe()
	return fmt.Sprintf("http://localhost:%d", port)
}

func (s *server) Stop() {
	s.s.Shutdown(context.Background())
}

func (s *server) Wait() {
	s.wg.Wait()
}

func (s *server) AuthorizationCode() (string, error) {
	return s.code, s.err
}

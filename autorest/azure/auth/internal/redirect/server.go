package redirect

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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

type Server interface {
	Start() string
	Stop()
	Wait()
	QueryParams() url.Values
}

type server struct {
	wg *sync.WaitGroup
	s  *http.Server
	qp url.Values
}

func NewServer() Server {
	rs := &server{
		wg: &sync.WaitGroup{},
		s:  &http.Server{},
	}
	return rs
}

func (s *server) Start() string {
	port := "8735"
	s.s.Addr = fmt.Sprintf(":%s", port)
	s.s.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(okPage))
		qp := r.URL.Query()
		if len(qp) > 0 {
			s.qp = qp
		}
		s.wg.Done()
	})
	s.wg.Add(2)
	go s.s.ListenAndServe()
	return fmt.Sprintf("http://localhost:%s", port)
}

func (s *server) Stop() {
	s.s.Shutdown(context.Background())
}

func (s *server) Wait() {
	s.wg.Wait()
}

func (s *server) QueryParams() url.Values {
	return s.qp
}

package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/pkg/browser"
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

func main() {
	clientID := "04b07795-8ddb-461a-bbee-02f9e1bf7b46"
	redirectURL, wg := startLocalServer()
	state := getStateID()
	resource := "https://management.azure.com/"
	authURL := fmt.Sprintf("https://login.microsoftonline.com/common/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&resource=%s&prompt=select_account",
		clientID,
		redirectURL,
		state,
		resource)
	err := browser.OpenURL(authURL)
	if err != nil {
		panic(err)
	}
	wg.Wait()
	//req, _ := http.NewRequest()
	/*dfc := auth.NewDeviceFlowConfig("10feaaf4-4ae4-45c2-9e4a-544c4d831e0e", "72f988bf-86f1-41af-91ab-2d7cd011db47")
	_, err := dfc.Authorizer()
	if err != nil {
		panic(err)
	}*/
}

func getStateID() string {
	buff := make([]byte, 64)
	rand.Read(buff)
	return strings.ToLower(base64.StdEncoding.EncodeToString(buff)[:20])
}

func startLocalServer() (string, *sync.WaitGroup) {
	wg := &sync.WaitGroup{}
	portNumber := ":8735"
	s := &http.Server{Addr: portNumber}
	s.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(okPage))
		wg.Done()
		//s.Shutdown(context.Background())
	})
	wg.Add(2)
	go s.ListenAndServe()
	return fmt.Sprintf("http://localhost%s", portNumber), wg
}

type redirectServer struct {
	wg *sync.WaitGroup
	s  *http.Server
	qp url.Values
}

func NewRedirectServer() *redirectServer {
	rs := &redirectServer{
		wg: &sync.WaitGroup{},
		s:  &http.Server{},
	}
	return rs
}

func (rs *redirectServer) Start() string {
	port := "8735"
	rs.s.Addr = fmt.Sprintf(":%s", port)
	rs.s.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(okPage))
		rs.qp = r.URL.Query()
		rs.wg.Done()
	})
	rs.wg.Add(2)
	go rs.s.ListenAndServe()
	return port
}

func (rs *redirectServer) Stop() {
	rs.s.Shutdown(context.Background())
}

func (rs *redirectServer) Wait() {
	rs.wg.Wait()
}

func (rs *redirectServer) QueryParams() url.Values {
	return rs.qp
}

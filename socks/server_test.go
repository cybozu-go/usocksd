package socks

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/cybozu-go/cmd"
)

func hasIPv6() bool {
	if os.Getenv("TRAVIS") == "true" {
		return false
	}
	ln, err := net.Listen("tcp", "[::1]:29343")
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

func TestServerBasic(t *testing.T) {
	t.Parallel()

	_, err := exec.LookPath("curl")
	if err != nil {
		t.Skip("curl not found")
	}

	addr := "localhost:20080"
	env := cmd.NewEnvironment(context.Background())
	s := &Server{
		Env: env,
	}
	ln, err := net.Listen("tcp", ":20080")
	if err != nil {
		t.Skip(err)
	}
	s.Serve(ln)

	mux := http.NewServeMux()
	mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {})
	hs := &http.Server{
		Addr:    ":20081",
		Handler: mux,
	}
	go hs.ListenAndServe()

	time.Sleep(10 * time.Millisecond)

	url1 := "http://localhost:20081/ok"
	curl := exec.Command("curl", "-4", "-I", "--socks4", addr, url1)
	out, err := curl.Output()
	if err != nil {
		t.Error(err)
		t.Log(string(out))
	}

	curl = exec.Command("curl", "-4", "-I", "--socks4a", addr, url1)
	out, err = curl.Output()
	if err != nil {
		t.Error(err)
		t.Log(string(out))
	}

	curl = exec.Command("curl", "-4", "-I", "--socks5", addr, url1)
	out, err = curl.Output()
	if err != nil {
		t.Error(err)
		t.Log(string(out))
	}

	curl = exec.Command("curl", "-4", "-I", "--socks5-hostname", addr, url1)
	out, err = curl.Output()
	if err != nil {
		t.Error(err)
		t.Log(string(out))
	}

	if !hasIPv6() {
		goto DONE
	}

	curl = exec.Command("curl", "-6", "-I", "--socks5", addr, url1)
	out, err = curl.Output()
	if err != nil {
		t.Error(err)
		t.Log(string(out))
	}

	curl = exec.Command("curl", "-6", "-I", "--socks5-hostname", addr, url1)
	out, err = curl.Output()
	if err != nil {
		t.Error(err)
		t.Log(string(out))
	}

DONE:
	env.Cancel(nil)
	err = env.Wait()
	if err != nil {
		t.Error(err)
	}
}

type authenticator struct{}

func (a authenticator) Authenticate(r *Request) bool {
	switch r.Username {
	case "root":
		return true
	case "user":
		return r.Password == "pass"
	}
	return false
}

func TestServerAuth(t *testing.T) {
	t.Parallel()

	_, err := exec.LookPath("curl")
	if err != nil {
		t.Skip("curl not found")
	}

	addr := "localhost:20082"
	env := cmd.NewEnvironment(context.Background())
	s := &Server{
		Auth: authenticator{},
		Env:  env,
	}
	ln, err := net.Listen("tcp", ":20082")
	if err != nil {
		t.Skip(err)
	}
	s.Serve(ln)

	mux := http.NewServeMux()
	mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {})
	hs := &http.Server{
		Addr:    ":20083",
		Handler: mux,
	}
	go hs.ListenAndServe()

	time.Sleep(10 * time.Millisecond)

	url1 := "http://localhost:20083/ok"
	curl := exec.Command("curl", "-4", "-I", "--socks4", addr, url1)
	err = curl.Run()
	if err == nil {
		t.Error("authentication is necessary")
	}

	curl = exec.Command("curl", "-4", "-I", "-U", "root:", "--socks4", addr, url1)
	out, err := curl.CombinedOutput()
	if err != nil {
		t.Error(err)
		t.Log(string(out))
	}

	curl = exec.Command("curl", "-4", "-I", "-U", "root:", "--socks4a", addr, url1)
	out, err = curl.CombinedOutput()
	if err != nil {
		t.Error(err)
		t.Log(string(out))
	}

	curl = exec.Command("curl", "-4", "-I", "--socks5", addr, url1)
	err = curl.Run()
	if err == nil {
		t.Error("authenticatoin is necessary")
	}

	curl = exec.Command("curl", "-4", "-I", "-U", "root:", "--socks5", addr, url1)
	out, err = curl.CombinedOutput()
	if err != nil {
		t.Error(err)
		t.Log(string(out))
	}

	curl = exec.Command("curl", "-4", "-I", "-U", "user:pass", "--socks5", addr, url1)
	out, err = curl.CombinedOutput()
	if err != nil {
		t.Error(err)
		t.Log(string(out))
	}

	curl = exec.Command("curl", "-4", "-I", "-U", "user:", "--socks5", addr, url1)
	out, err = curl.CombinedOutput()
	if err == nil {
		t.Error("authentication should fail")
	}

	curl = exec.Command("curl", "-4", "-I", "-U", "user:bad", "--socks5", addr, url1)
	out, err = curl.CombinedOutput()
	if err == nil {
		t.Error("authentication should fail")
	}

	if !hasIPv6() {
		goto DONE
	}

	curl = exec.Command("curl", "-6", "-I", "-U", "user:pass", "--socks5", addr, url1)
	out, err = curl.CombinedOutput()
	if err != nil {
		t.Error(err)
		t.Log(string(out))
	}

DONE:
	env.Cancel(nil)
	err = env.Wait()
	if err != nil {
		t.Error(err)
	}
}

type rules struct{}

func (ru rules) Match(r *Request) bool {
	// Allow only request with hostname
	return len(r.Hostname) > 0
}

func TestServerRules(t *testing.T) {
	t.Parallel()

	_, err := exec.LookPath("curl")
	if err != nil {
		t.Skip("curl not found")
	}

	addr := "localhost:20084"
	env := cmd.NewEnvironment(context.Background())
	s := &Server{
		Rules: rules{},
		Env:   env,
	}
	ln, err := net.Listen("tcp", ":20084")
	if err != nil {
		t.Skip(err)
	}
	s.Serve(ln)

	mux := http.NewServeMux()
	mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {})
	hs := &http.Server{
		Addr:    ":20085",
		Handler: mux,
	}
	go hs.ListenAndServe()

	time.Sleep(10 * time.Millisecond)

	url1 := "http://localhost:20085/ok"
	curl := exec.Command("curl", "-4", "-I", "--socks4", addr, url1)
	err = curl.Run()
	if err == nil {
		t.Error("SOCKS4 should be denied")
	}

	curl = exec.Command("curl", "-4", "-I", "--socks4a", addr, url1)
	out, err := curl.CombinedOutput()
	if err != nil {
		t.Error(err)
		t.Log(string(out))
	}

	curl = exec.Command("curl", "-4", "-I", "--socks5", addr, url1)
	err = curl.Run()
	if err == nil {
		t.Error("SOCKS5 w/o hostname should be denied")
	}

	curl = exec.Command("curl", "-4", "-I", "--socks5-hostname", addr, url1)
	out, err = curl.CombinedOutput()
	if err != nil {
		t.Error(err)
		t.Log(string(out))
	}

	if !hasIPv6() {
		goto DONE
	}

	curl = exec.Command("curl", "-6", "-I", "--socks5-hostname", addr, url1)
	out, err = curl.CombinedOutput()
	if err != nil {
		t.Error(err)
		t.Log(string(out))
	}

DONE:
	env.Cancel(nil)
	err = env.Wait()
	if err != nil {
		t.Error(err)
	}
}

func TestServerError(t *testing.T) {
	t.Parallel()

	_, err := exec.LookPath("curl")
	if err != nil {
		t.Skip("curl not found")
	}

	addr := "localhost:20086"
	env := cmd.NewEnvironment(context.Background())
	s := &Server{
		Env: env,
	}
	ln, err := net.Listen("tcp", ":20086")
	if err != nil {
		t.Skip(err)
	}
	s.Serve(ln)

	time.Sleep(10 * time.Millisecond)

	url1 := "http://localhost:20087/ok"
	curl := exec.Command("curl", "-4", "-I", "--socks4", addr, url1)
	out, err := curl.CombinedOutput()
	if err == nil {
		t.Error("no listener on localhost:20087")
	}
	t.Log(string(out))

	curl = exec.Command("curl", "-4", "-I", "--socks5", addr, url1)
	out, err = curl.CombinedOutput()
	if err == nil {
		t.Error("no listener on localhost:20087")
	}
	t.Log(string(out))

	env.Cancel(nil)
	err = env.Wait()
	if err != nil {
		t.Error(err)
	}
}

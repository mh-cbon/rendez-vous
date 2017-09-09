package main_test

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	clean()
	defer clean()
	if err := build(); err != nil {
		t.Fatal(err)
	}

	t.Run("1", func(t *testing.T) {
		rv, err := runRendezVous("8070")
		if err != nil {
			t.Fatal(err)
		}
		defer rv.Process.Kill()
		defer rv.Process.Release()

		if err := runPing(":8070"); err != nil {
			t.Error(err)
		}
	})

	t.Run("2", func(t *testing.T) {
		rv, err := runRendezVous("8090")
		if err != nil {
			t.Fatal(err)
		}
		defer rv.Process.Kill()
		defer rv.Process.Release()

		pvk := "202d229c0f09f41c858066496b21c27e59266ec7c5b0933275518b351da5e92e"
		ws, err := runWebsite(":8090", "8091", "8092", pvk)
		if err != nil {
			t.Error(err)
		}
		defer ws.Process.Kill()
		defer ws.Process.Release()

		err = runHttpGet(":8090", "http://127.0.0.1:8091/index.html")
		if err != nil {
			t.Error(err)
		}

		err = runHttpGet(":8090", "http://b6b8113748fe0795658fa9d6ab3e36d27d72e97b7df407e7a8080d61ec405d74.me.com/index.html")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("3", func(t *testing.T) {
		rv, err := runRendezVous("8080")
		if err != nil {
			t.Fatal(err)
		}
		defer rv.Process.Kill()
		defer rv.Process.Release()

		pvk := "202d229c0f09f41c858066496b21c27e59266ec7c5b0933275518b351da5e92e"
		ws, err := runWebsite(":8080", "8081", "8082", pvk)
		if err != nil {
			t.Error(err)
		}
		defer ws.Process.Kill()
		defer ws.Process.Release()

		bw, err := runBrowser(":8080", "8083", "8084", "8085")
		if err != nil {
			t.Error(err)
		}
		defer bw.Process.Kill()
		defer bw.Process.Release()

		{
			// port 8005 runs the ws of the browser
			err := geturl(http.DefaultClient, "http://127.0.0.1:8085/index.html")
			if err != nil {
				t.Error(err)
			}
		}
		// port 8084 is a proxy to go either
		//- in the rendez-vous network
		//- in the browser ws
		//- in the regular internet
		proxyUrl, err := url.Parse("http://127.0.0.1:8084")
		if err != nil {
			t.Error(err)
		}
		client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
		{
			err := geturl(client, "http://b6b8113748fe0795658fa9d6ab3e36d27d72e97b7df407e7a8080d61ec405d74.me.com/index.html")
			if err != nil {
				t.Error(err)
			}
		}
	})
}

func geturl(client *http.Client, u string) error {
	fmt.Println("HTTP GET ", u)
	res, err := client.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = io.Copy(os.Stdout, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func build() error {
	cmd := exec.Command("go", "build", "-o", "t", "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}
func clean() error {
	exec.Command("killall", "t").Run()
	return os.Remove("t")
}

func makeCmd(b string, args ...string) *exec.Cmd {
	fmt.Println("go run main.go", strings.Join(args, " "))
	cmd := exec.Command(b, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd
}

func runRendezVous(port string) (*exec.Cmd, error) {
	cmd := makeCmd("./t", "serve", "-l", port)
	return cmd, timeout(cmd.Run, time.Second)
}

func runPing(remote string) error {
	cmd := makeCmd("./t", "client", "-q", "ping", "-r", remote)
	return cmd.Run()
}

func runWebsite(remote, listen, local, pvk string) (*exec.Cmd, error) {
	cmd := makeCmd("./t", "website", "-r", remote, "-l", listen, "--local", local, "--pvk", pvk, "--dir", "demows")
	return cmd, timeout(cmd.Run, time.Second)
}

func runBrowser(remote, listen, proxy, ws string) (*exec.Cmd, error) {
	cmd := makeCmd("./t", "browser", "-r", remote, "-l", listen, "--ws", ws, "--proxy", proxy, "--headless")
	return cmd, timeout(cmd.Run, time.Second)
}

func runHttpGet(remote, url string) error {
	cmd := makeCmd("./t", "http", "--url", url, "--remote", remote)
	return cmd.Run()
}

func timeout(do func() error, d time.Duration) error {
	rcv := make(chan error)
	go func() {
		rcv <- do()
	}()
	select {
	case err := <-rcv:
		close(rcv)
		if err != nil {
			return err
		}
	case <-time.After(d):
	}
	return nil
}

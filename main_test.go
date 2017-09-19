package main_test

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"
)

var srcIP = "0.0.0.0"
var dstIP = "127.0.0.1" // on windows, never write to 0.0.0.0:... or :..., it will fail with an error such
// ... write udp [::]:50929->:8070: wsasendto: The requested address is not valid in its context.

func Test1(t *testing.T) {
	clean()
	defer clean()
	if err := build(); err != nil {
		t.Fatal(err)
	}

	t.Run("1", func(t *testing.T) {
		rv, err := runRendezVous(srcIP + ":8070")
		if err != nil {
			t.Fatal(err)
		}
		defer rv.Process.Kill()
		defer rv.Process.Release()

		if err := runPing(dstIP + ":8070"); err != nil {
			t.Error(err)
		}
	})

	t.Run("2", func(t *testing.T) {
		rv, err := runRendezVous(srcIP + ":8090")
		if err != nil {
			t.Fatal(err)
		}
		defer rv.Process.Kill()
		defer rv.Process.Release()

		pvk := "504bc61393e5d7ea991dbfad4d5bb98093562d472fa22d425a35bcd46341d8f678e7d4c5aa13e3d9a538a5aa2a027cb5343931a48a6fd7b7b1ae699ec8125f12"
		ws, err := runWebsite(dstIP+":8090", srcIP+":8091", srcIP+":8092", pvk)
		if err != nil {
			t.Error(err)
		}
		defer ws.Process.Kill()
		defer ws.Process.Release()

		err = runHttpGet(dstIP+":8090", "http://127.0.0.1:8091/index.html")
		if err != nil {
			t.Error(err)
		}

		err = runHttpGet(dstIP+":8090", "http://78e7d4c5aa13e3d9a538a5aa2a027cb5343931a48a6fd7b7b1ae699ec8125f12.me.com/index.html")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("3", func(t *testing.T) {
		rv, err := runRendezVous(srcIP + ":8080")
		if err != nil {
			t.Fatal(err)
		}
		defer rv.Process.Kill()
		defer rv.Process.Release()

		pvk := "504bc61393e5d7ea991dbfad4d5bb98093562d472fa22d425a35bcd46341d8f678e7d4c5aa13e3d9a538a5aa2a027cb5343931a48a6fd7b7b1ae699ec8125f12"
		ws, err := runWebsite(dstIP+":8080", srcIP+":8081", srcIP+":8082", pvk)
		if err != nil {
			t.Error(err)
		}
		defer ws.Process.Kill()
		defer ws.Process.Release()

		bw, err := runBrowser(dstIP+":8080", srcIP+":8083", srcIP+":8084", srcIP+":8085")
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
			err := geturl(client, "http://78e7d4c5aa13e3d9a538a5aa2a027cb5343931a48a6fd7b7b1ae699ec8125f12.me.com/index.html")
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

var exeFile = "./t"

func init() {
	if runtime.GOOS == "windows" {
		exeFile = "./t.exe"
	}
}

func build() error {
	cmd := exec.Command("go", "build", "-o", exeFile, "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}
func clean() error {
	exec.Command("killall", exeFile).Run()
	exec.Command("cmd", "/C", "taskkill", "/T", "/F", "/IM", exeFile[2:]).Run()
	return os.Remove(exeFile)
}

func makeCmd(b string, args ...string) *exec.Cmd {
	fmt.Println("go run main.go", strings.Join(args, " "))
	cmd := exec.Command(b, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd
}

func runRendezVous(port string) (*exec.Cmd, error) {
	cmd := makeCmd(exeFile, "serve", "-l", port)
	return cmd, timeout(cmd.Run, time.Second)
}

func runPing(remote string) error {
	cmd := makeCmd(exeFile, "client", "-q", "ping", "-r", remote)
	return cmd.Run()
}

func runWebsite(remote, listen, local, pvk string) (*exec.Cmd, error) {
	cmd := makeCmd(exeFile, "website", "-r", remote, "-l", listen, "--local", local, "--pvk", pvk, "--dir", "_samples/static/assets")
	return cmd, timeout(cmd.Run, time.Second)
}

func runBrowser(remote, listen, proxy, ws string) (*exec.Cmd, error) {
	cmd := makeCmd(exeFile, "browser", "-r", remote, "-l", listen, "--ws", ws, "--proxy", proxy, "--headless")
	return cmd, timeout(cmd.Run, time.Second)
}

func runHttpGet(remote, url string) error {
	cmd := makeCmd(exeFile, "http", "--url", url, "--remote", remote)
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

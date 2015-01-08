// "THE BEER-WARE LICENSE" (Revision 42):
// <kevin.gillieron@gw-computing.net> wrote this file. As long as you retain
// this notice you can do whatever you want with this stuff. If we meet some
// day, and you think this stuff is worth it, you can buy me a beer in return
// Kevin Gillieron

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"

	"github.com/go-fsnotify/fsnotify"
)

// program version number
const version = "1.2.1"

const (
	doctype    = "<!DOCTYPE html>"
	headFormat = "<head><meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\">%s</head>"
)

type markdown struct {
	Text    string `json:"text"`
	Mode    string `json:"mode"`
	Context string `json:"context"`
}

func fatal(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func readBody(ioBody io.ReadCloser) string {
	body, err := ioutil.ReadAll(ioBody)
	if err != nil {
		fatal(err)
	}
	return string(body)
}

func createTempFile() *os.File {
	f, err := ioutil.TempFile(os.TempDir(), "ghmd")
	if err != nil {
		fatal(err)
	}

	return f
}

func render(path string, out *os.File, refresh bool) {
	md, err := ioutil.ReadFile(path)
	if err != nil {
		fatal(err)
	}

	buf, err := json.Marshal(markdown{
		Text: string(md),
		Mode: "markdown",
	})
	if err != nil {
		fatal(err)
	}

	resp, err := http.Post("https://api.github.com/markdown",
		"application/json", bytes.NewBuffer(buf))
	if err != nil {
		fatal(err)
		return
	}
	defer resp.Body.Close()

	err = out.Truncate(0)
	if err != nil {
		fatal(err)
	}

	_, err = out.Seek(0, 0)
	if err != nil {
		fatal(err)
	}

	var head string

	if refresh {
		head = fmt.Sprintf(headFormat, "<meta http-equiv=\"refresh\" content=\"2\">")
	} else {
		head = fmt.Sprintf(headFormat, "")
	}

	fmt.Fprintln(out, doctype, head, `<body class="markdown-body"><style>`, githubCSS, `</style>`)
	fmt.Fprintln(out, readBody(resp.Body), "</body>")
}

func watch(path string, out *os.File) {
	fileName := filepath.Base(path)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatal(err)
	}
	defer watcher.Close()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if filepath.Base(ev.Name) == fileName && ev.Op == fsnotify.Write {
					render(path, out, true)
				}
			case err := <-watcher.Errors:
				fmt.Fprintln(os.Stderr, "error:", err)
			}
		}
	}()

	err = watcher.Add(filepath.Dir(path))
	if err != nil {
		fatal(err)
	}

	<-done
}

// defaultCmd tries to determine the user's default web browser.
func defaultCmd() (string, error) {
	switch runtime.GOOS {
	case "linux", "freebsd", "dragonfly", "openbsd", "netbsd":
		return "xdg-open", nil
	case "darwin":
		return "open -a Safari", nil
	default:
		return "", errors.New("Not yet implemented for " + runtime.GOOS)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [markdown file]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	var versionFlag bool
	flag.BoolVar(&versionFlag, "v", false, "Display program version number.")

	var watchFlag bool
	flag.BoolVar(&watchFlag, "w", false, "Watch Markdown file change.")

	var outputFile string
	flag.StringVar(&outputFile, "o", "", "Output file. If not supplied, a temporary file will be created.")

	var runCmdFlag bool
	flag.BoolVar(&runCmdFlag, "r", false, "Run generated HTML file. It will try to open with your default web browser.")

	flag.Parse()

	if versionFlag {
		fmt.Println(os.Args[0], version)
		os.Exit(0)
	}

	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "Invalid # of arguments!")
		flag.Usage()
	}

	var err error

	path := flag.Arg(0)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		fatal("input markdown file does not exist")
	}

	// create output file
	var out *os.File
	if outputFile == "" {
		out = createTempFile()
		fmt.Println("Temporary file created!")
		fmt.Println(out.Name())
	} else {
		out, err = os.Create(outputFile)
		if err != nil {
			fatal(err)
		}
	}
	defer func() {
		out.Close()
		// only remove if it's a temporary file and if -w switch is enabled
		if outputFile == "" && watchFlag {
			err = os.Remove(out.Name())
			if err != nil {
				fatal(err)
			}
		}
	}()

	render(path, out, false)

	if runCmdFlag {
		runCmd, err := defaultCmd()
		if err != nil {
			fatal(err)
		}

		cmd := exec.Command(runCmd, out.Name())
		err = cmd.Start()
		if err != nil {
			fatal(err)
		}
		cmd.Wait()
	}

	if watchFlag {
		watch(path, out)
	}
}

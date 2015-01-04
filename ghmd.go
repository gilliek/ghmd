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

// CSS taken from https://gist.github.com/andyferra/2554919
const githubCSS = `body { font-family: Helvetica, arial, sans-serif; font-size: 14px; line-height: 1.6; padding-top: 10px; padding-bottom: 10px; background-color: white; padding: 30px; }
body > *:first-child { margin-top: 0 !important; }
body > *:last-child { margin-bottom: 0 !important; }
a { color: #4183C4; }
a.absent { color: #cc0000; }
a.anchor { display: block; padding-left: 30px; margin-left: -30px; cursor: pointer; position: absolute; top: 0; left: 0; bottom: 0; }
h1, h2, h3, h4, h5, h6 { margin: 20px 0 10px; padding: 0; font-weight: bold; -webkit-font-smoothing: antialiased; cursor: text; position: relative; }
h1 tt, h1 code { font-size: inherit; }
h2 tt, h2 code { font-size: inherit; }
h3 tt, h3 code { font-size: inherit; }
h4 tt, h4 code { font-size: inherit; }
h5 tt, h5 code { font-size: inherit; }
h6 tt, h6 code { font-size: inherit; }
h1 { font-size: 28px; color: black; }
h2 { font-size: 24px; border-bottom: 1px solid #cccccc; color: black; }
h3 { font-size: 18px; }
h4 { font-size: 16px; }
h5 { font-size: 14px; }
h6 { color: #777777; font-size: 14px; }
p, blockquote, ul, ol, dl, li, table, pre { margin: 15px 0; }
hr { border: 0 none; color: #cccccc; height: 4px; padding: 0; }
body > h2:first-child { margin-top: 0; padding-top: 0; }
body > h1:first-child { margin-top: 0; padding-top: 0; }
body > h1:first-child + h2 { margin-top: 0; padding-top: 0; }
body > h3:first-child, body > h4:first-child, body > h5:first-child, body > h6:first-child { margin-top: 0; padding-top: 0; }
a:first-child h1, a:first-child h2, a:first-child h3, a:first-child h4, a:first-child h5, a:first-child h6 { margin-top: 0; padding-top: 0; }
h1 p, h2 p, h3 p, h4 p, h5 p, h6 p { margin-top: 0; }
li p.first { display: inline-block; }
ul, ol { padding-left: 30px; }
ul :first-child, ol :first-child { margin-top: 0; }
ul :last-child, ol :last-child { margin-bottom: 0; }
dl { padding: 0; }
dl dt { font-size: 14px; font-weight: bold; font-style: italic; padding: 0;margin: 15px 0 5px; }
dl dt:first-child { padding: 0; }
dl dt > :first-child { margin-top: 0; }
dl dt > :last-child { margin-bottom: 0; }
dl dd { margin: 0 0 15px; padding: 0 15px; }
dl dd > :first-child { margin-top: 0; }
dl dd > :last-child { margin-bottom: 0; }
blockquote { border-left: 4px solid #dddddd; padding: 0 15px; color: #777777; }
blockquote > :first-child { margin-top: 0; }
blockquote > :last-child { margin-bottom: 0; }
table { padding: 0; border-collapse: collapse; }
table tr { border-top: 1px solid #cccccc; background-color: white; margin: 0; padding: 0; }
table tr:nth-child(2n) { background-color: #f8f8f8; }
table tr th { font-weight: bold; border: 1px solid #cccccc; text-align: left; margin: 0; padding: 6px 13px; }
table tr td { border: 1px solid #cccccc; text-align: left; margin: 0; padding: 6px 13px; }
table tr th :first-child, table tr td :first-child { margin-top: 0; }
table tr th :last-child, table tr td :last-child { margin-bottom: 0; }
img { max-width: 100%; }
span.frame { display: block; overflow: hidden; }
span.frame > span { border: 1px solid #dddddd; display: block; float: left;overflow: hidden; margin: 13px 0 0; padding: 7px; width: auto; }
span.frame span img { display: block; float: left; }
span.frame span span { clear: both; color: #333333; display: block; padding: 5px 0 0; }
span.align-center { display: block; overflow: hidden; clear: both; }
span.align-center > span { display: block; overflow: hidden; margin: 13px auto 0; text-align: center; }
span.align-center span img { margin: 0 auto; text-align: center; }
span.align-right { display: block; overflow: hidden; clear: both; }
span.align-right > span { display: block; overflow: hidden; margin: 13px 0 0; text-align: right; }
span.align-right span img { margin: 0; text-align: right; }
span.float-left { display: block; margin-right: 13px; overflow: hidden; float: left; }
span.float-left span { margin: 13px 0 0; }
span.float-right { display: block; margin-left: 13px; overflow: hidden; float: right; }
  span.float-right > span { display: block; overflow: hidden; margin: 13px auto 0; text-align: right; }
code, tt { margin: 0 2px; padding: 0 5px; white-space: nowrap; border: 1px solid #eaeaea; background-color: #f8f8f8; border-radius: 3px; } 
pre code { margin: 0; padding: 0; white-space: pre; border: none; background: transparent; }
.highlight pre { background-color: #f8f8f8; border: 1px solid #cccccc; font-size: 13px; line-height: 19px; overflow: auto; padding: 6px 10px; border-radius: 3px; }
pre { background-color: #f8f8f8; border: 1px solid #cccccc; font-size: 13px; line-height: 19px; overflow: auto; padding: 6px 10px; border-radius: 3px; }
  pre code, pre tt { background-color: transparent; border: none; }`

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

	fmt.Fprintln(out, doctype, head, "<body><style>", githubCSS, "</style>")
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
	default: // BSDs OS
		// TODO
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

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/urfave/cli"
	"github.com/zhanbei/serve-static"
)

const OptionNameEnableVirtualHosting = "enable-virtual-hosting"
const OptionNameNoTrailingSlash = "no-trailing-slash"
const OptionNameDirectoryListing = "enable-directory-listing"

var mUsingVirtualHost = false
var mNoTrailingSlash = true
var mDirectoryListing = false

// Disable directory listing with http.FileServer
func GetNoDirListingHandler(rootDir string) http.Handler {
	// @see https://stackoverflow.com/questions/26559557/how-do-you-serve-a-static-html-file-using-a-go-web-server
	// @see https://groups.google.com/forum/#!msg/golang-nuts/bStLPdIVM6w/hidTJgDZpHcJ
	// @see https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
	handler := http.FileServer(http.Dir(rootDir))
	if mDirectoryListing {
		fmt.Println("Enabled directory listing.")
		return handler
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqPath := r.URL.Path
		if strings.HasSuffix(reqPath, "/") && reqPath != "/" {
			exists, _ := servestatic.IsFileRegular(rootDir, reqPath, servestatic.IndexDotHtml)
			if !exists {
				http.NotFound(w, r)
				return
			}
		}
		handler.ServeHTTP(w, r)
	})
}

func Action(c *cli.Context) error {
	if !mNoTrailingSlash && mUsingVirtualHost {
		fmt.Println("ERROR: Sorry, currently virtual hosting is supported only in the " + OptionNameNoTrailingSlash + " mode.")
		log.Fatal("You may add the --" + OptionNameNoTrailingSlash + " option to use --" + OptionNameEnableVirtualHosting + " option.")
	}

	if c.NArg() <= 0 {
		log.Fatal("Please specify a port, like `static-server 8080`.")
	}
	address := c.Args().Get(0)
	port, err := strconv.Atoi(address)
	if err != nil {
		// Check the address.
	} else {
		// The address is only a port.
		if port < 1 || 65535 < port {
			log.Fatal("ERROR: unavailable port[" + strconv.Itoa(port) + "]; make sure http port is number and is limited to <0-65535>.")
		}
		if port <= 1024 {
			fmt.Println("WARNING: the port[" + strconv.Itoa(port) + "] specified is not bigger than 1024; root privileges may be needed!")
		}
		address = ":" + strconv.Itoa(port)
	}

	rootDir := "."
	if c.NArg() > 1 {
		rootDir = c.Args().Get(1)
	}
	rootDir, err = filepath.Abs(rootDir)
	if err != nil {
		fmt.Println("ERROR: The specified www-root-directory is invalid:" + rootDir)
		log.Fatal(err)
	}

	var handler http.Handler
	if !mNoTrailingSlash {
		// Hosting in the normal mode.
		handler = GetNoDirListingHandler(rootDir)
	} else {
		fmt.Println("Hosting static files in the " + OptionNameNoTrailingSlash + " mode.")
		if mUsingVirtualHost {
			fmt.Println("Enabled virtual hosting based on request.Host; @see https://en.wikipedia.org/wiki/Virtual_hosting.")
		}
		mStaticServer, err := servestatic.NewFileServer(rootDir, mUsingVirtualHost)
		if err != nil {
			fmt.Println("ERROR: The specified www-root-directory is invalid:" + rootDir)
			log.Fatal(err)
		}
		handler = mStaticServer
	}
	fmt.Println("Looking after directory:", rootDir)
	handler = handlers.CombinedLoggingHandler(os.Stdout, handler)
	fmt.Println("Server is running at:", address)
	http.ListenAndServe(address, handler)
	//fmt.Println("listening:", address, mUsingVirtualHost, mNoTrailingSlash)
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "static-server"
	app.Usage = "A static server in Go, supporting hosting static files in the no-trailing-slash version."
	app.Version = "0.9.0"
	app.Description = "A static server in Go, supporting hosting static files in the no-trailing-slash version."
	app.UsageText = "static-server [global options] [<http-address>:]<http-port> <www-root-directory>"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: OptionNameEnableVirtualHosting,

			Usage: "Whether to enable virtual hosting; @see https://en.wikipedia.org/wiki/Virtual_hosting",

			Destination: &mUsingVirtualHost,
		},
		cli.BoolFlag{
			Name: OptionNameNoTrailingSlash,

			Usage: "Hosting static files in the " + OptionNameNoTrailingSlash + " mode.",

			Destination: &mNoTrailingSlash,
		},
		cli.BoolFlag{
			Name: OptionNameDirectoryListing,

			Usage: "Listing files of a directory if the index.html is not found when in the normal mode.",

			Destination: &mDirectoryListing,
		},
	}
	app.Action = Action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

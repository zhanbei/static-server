package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
	. "github.com/zhanbei/static-server/conf"
	"github.com/zhanbei/static-server/configs"
	"github.com/zhanbei/static-server/core"
	"github.com/zhanbei/static-server/db"
	"github.com/zhanbei/static-server/helpers/terminator"
	"github.com/zhanbei/static-server/recorder"
	"github.com/zhanbei/static-server/utils"
)

var ops = NewDefaultServerOptions()

var OptionConfiguresFile = ""

// The primary program entrance.
// (Cli Arguments Receiver + Configuration File Parser + MongoDB Driver)
// Support more custom built, like for lite/medium/heavy programs, for cli/gui(with different themes) modes, and for linux/windows/mac platforms.
// @see [Support multiple entrances and keep the current one as the primary. · Issue #6 · zhanbei/static-server](https://github.com/zhanbei/static-server/issues/6)
func main() {
	app := cli.NewApp()
	app.Name = "static-server"
	app.Usage = "A static server in Go, supporting hosting static files in the no-trailing-slash version."
	app.Version = "0.9.1"
	app.Description = "A static server in Go, supporting hosting static files in the no-trailing-slash version."
	app.UsageText = "static-server [global options] [<http-address>:]<http-port> <www-root-directory>"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "configure",

			Usage: "The configuration file to be used.",

			Destination: &OptionConfiguresFile,
		},
		cli.BoolFlag{
			Name: OptionNameEnableVirtualHosting,

			Usage: "Whether to enable virtual hosting; @see https://en.wikipedia.org/wiki/Virtual_hosting",

			Destination: &ops.UsingVirtualHost,
		},
		cli.BoolFlag{
			Name: OptionNameNoTrailingSlash,

			Usage: "Hosting static files in the " + OptionNameNoTrailingSlash + " mode.",

			Destination: &ops.NoTrailingSlash,
		},
		cli.BoolFlag{
			Name: OptionNameDirectoryListing,

			Usage: "Listing files of a directory if the index.html is not found when in the normal mode.",

			Destination: &ops.DirectoryListing,
		},
	}
	app.Action = Action

	err := app.Run(os.Args)
	if err != nil {
		terminator.ExitWithPreLaunchServerError(err, "Loading configures from environment variables failed!")
	}
}

func Action(c *cli.Context) error {
	if utils.NotEmpty(OptionConfiguresFile) {
		return ActionConfigurationFile(c, ops, OptionConfiguresFile)
	} else {
		return ActionCliArguments(c, ops)
	}
}

// FIX-ME Use a default configuration file, like `vhss.(yaml|toml|json)`.
func ActionConfigurationFile(c *cli.Context, ops *ServerOptions, confFile string) error {
	// Prefer the cli arguments, over the configuration file.
	rawAddress, rawRootDir := "", ""
	if c.NArg() > 0 {
		rawAddress = c.Args().Get(0)
	}
	if c.NArg() > 1 {
		rawRootDir = c.Args().Get(1)
	}

	cfg, err := configs.LoadServerConfigures(confFile, ops, strings.TrimSpace(rawAddress), strings.TrimSpace(rawRootDir))
	if err != nil {
		terminator.ExitWithConfigError(err, "Loading and validating the configures failed!")
	}
	err = cfg.ValidateFile()
	if err != nil {
		terminator.ExitWithPreLaunchServerError(err, "Validating the required resources following configures failed!")
	}
	bts, err := json.Marshal(cfg)
	fmt.Println("Loading configures:", string(bts))
	fmt.Println(cfg, cfg.Server, cfg.Loggers, cfg.MongoDbOptions, cfg.GorillaOptions, confFile)

	loggers := recorder.GetActiveRecorders(cfg.Loggers)

	mon := cfg.MongoDbOptions
	if mon != nil && mon.Enabled {
		err = db.ConnectToMongoDb(cfg.MongoDbOptions)
		if err != nil {
			terminator.ExitWithPreLaunchServerError(err, "Connecting to mongodb failed!")
		}
		loggers = append(loggers, db.GetMongoRecorder(cfg.MongoDbOptions))
	}

	gor := cfg.GorillaOptions
	if (gor == nil || !gor.Enabled) && len(loggers) == 0 {
		// Add a default console(stdout) logger when there is no logger configured!
		loggers = append(loggers, recorder.GetDefaultRecorder())
	}

	return core.RealServer(cfg, loggers)
}

func ActionCliArguments(c *cli.Context, ops *ServerOptions) error {
	ops.ValidateOrExit()

	if c.NArg() <= 0 {
		terminator.ExitWithConfigError(nil, "Please specify a port, like `static-server 8080`.")
	}
	address := c.Args().Get(0)
	address, _ = ValidateArgAddressOrExit(address)

	rootDir := "."
	if c.NArg() > 1 {
		rootDir = c.Args().Get(1)
	}
	rootDir = ValidateArgRootDirOrExit(rootDir)

	fmt.Println("Loading arguments:", address, rootDir, ops)

	cfg := &Configure{rootDir, address, NewDefaultAppOptions(), ops, nil, nil, nil, nil}

	//fmt.Println("listening:", address, mUsingVirtualHost, mNoTrailingSlash)
	return core.RealServer(cfg, recorder.GetDefaultRecorders())
}

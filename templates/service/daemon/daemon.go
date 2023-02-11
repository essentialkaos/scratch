package daemon

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"os"
	"strings"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/knf"
	"github.com/essentialkaos/ek/v12/log"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/pid"
	"github.com/essentialkaos/ek/v12/signal"
	"github.com/essentialkaos/ek/v12/usage"

	knfv "github.com/essentialkaos/ek/v12/knf/validators"
	knff "github.com/essentialkaos/ek/v12/knf/validators/fs"

	"github.com/essentialkaos/{{SHORT_NAME}}/daemon/support"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Basic service info
const (
	APP  = "{{NAME}}"
	VER  = "{{VERSION}}"
	DESC = "{{DESC}}"
)

// Options
const (
	OPT_CONFIG   = "c:config"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_VERB_VER = "vv:verbose-version"
)

// Configuration file properties
const (
	LOG_DIR   = "log:dir"
	LOG_FILE  = "log:file"
	LOG_PERMS = "log:perms"
	LOG_LEVEL = "log:level"
)

// Pid file info
const (
	PID_DIR  = "/var/run/{{SHORT_NAME}}"
	PID_FILE = "{{SHORT_NAME}}.pid"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// optMap contains information about all supported options
var optMap = options.Map{
	OPT_CONFIG:   {Value: "/etc/{{SHORT_NAME}}.knf"},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:      {Type: options.BOOL, Alias: "ver"},

	OPT_VERB_VER: {Type: options.BOOL},
}

// ////////////////////////////////////////////////////////////////////////////////// //

func Init(gitRev string, gomod []byte) {
	_, errs := options.Parse(optMap)

	if len(errs) != 0 {
		for _, err := range errs {
			printError(err.Error())
		}

		os.Exit(1)
	}

	configureUI()

	switch {
	case options.GetB(OPT_VER):
		os.Exit(showAbout(gitRev))
	case options.GetB(OPT_HELP):
		os.Exit(showUsage())
	case options.GetB(OPT_VERB_VER):
		os.Exit(support.Show(APP, VER, gitRev, gomod))
	}

	loadConfig()
	validateConfig()
	registerSignalHandlers()
	setupLogger()
	createPidFile()

	log.Aux(strings.Repeat("-", 80))
	log.Aux("%s %s starting…", APP, VER)

	start()
}

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	term := os.Getenv("TERM")

	fmtc.DisableColors = true

	if term != "" {
		switch {
		case strings.Contains(term, "xterm"),
			strings.Contains(term, "color"),
			term == "screen":
			fmtc.DisableColors = false
		}
	}

	if os.Getenv("NO_COLOR") != "" {
		fmtc.DisableColors = true
	}
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}
}

// loadConfig reads and parses configuration file
func loadConfig() {
	err := knf.Global(options.GetS(OPT_CONFIG))

	if err != nil {
		printErrorAndExit(err.Error())
	}
}

// validateConfig validates configuration file values
func validateConfig() {
	errs := knf.Validate([]*knf.Validator{
		{LOG_DIR, knff.Perms, "DW"},
		{LOG_DIR, knff.Perms, "DX"},

		{LOG_LEVEL, knfv.NotContains, []string{"debug", "info", "warn", "error", "crit"}},
	})

	if len(errs) != 0 {
		printError("Error while configuration file validation:")

		for _, err := range errs {
			printError("  %v", err)
		}

		os.Exit(1)
	}
}

// registerSignalHandlers registers signal handlers
func registerSignalHandlers() {
	signal.Handlers{
		signal.TERM: termSignalHandler,
		signal.INT:  intSignalHandler,
		signal.HUP:  hupSignalHandler,
	}.TrackAsync()
}

// setupLogger confugures logger subsystems
func setupLogger() {
	err := log.Set(knf.GetS(LOG_FILE), knf.GetM(LOG_PERMS, 644))

	if err != nil {
		printErrorAndExit(err.Error())
	}

	err = log.MinLevel(knf.GetS(LOG_LEVEL))

	if err != nil {
		printErrorAndExit(err.Error())
	}
}

// createPidFile creates PID file
func createPidFile() {
	pid.Dir = PID_DIR

	err := pid.Create(PID_FILE)

	if err != nil {
		printErrorAndExit(err.Error())
	}
}

func start() {
	// DO YOUR STUFF HERE
}

// intSignalHandler is INT signal handler
func intSignalHandler() {
	log.Aux("Received INT signal, shutdown…")
	shutdown(0)
}

// termSignalHandler is TERM signal handler
func termSignalHandler() {
	log.Aux("Received TERM signal, shutdown…")
	shutdown(0)
}

// hupSignalHandler is HUP signal handler
func hupSignalHandler() {
	log.Info("Received HUP signal, log will be reopened…")
	log.Reopen()
	log.Info("Log reopened by HUP signal")
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	if len(a) == 0 {
		fmtc.Fprintln(os.Stderr, "{r}"+f+"{!}")
	} else {
		fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
	}
}

// printError prints warning message to console
func printWarn(f string, a ...interface{}) {
	if len(a) == 0 {
		fmtc.Fprintln(os.Stderr, "{y}"+f+"{!}")
	} else {
		fmtc.Fprintf(os.Stderr, "{y}"+f+"{!}\n", a...)
	}
}

// printErrorAndExit print error mesage and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// shutdown stops deamon
func shutdown(code int) {
	pid.Remove(PID_FILE)
	os.Exit(code)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showUsage prints usage info
func showUsage() int {
	genUsage().Render()
	return 0
}

// showAbout prints info about version
func showAbout(gitRev string) int {
	genAbout(gitRev).Render()
	return 0
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo()

	info.AddOption(OPT_CONFIG, "Path to configuration file", "file")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2009,
		Owner:   "ESSENTIAL KAOS",
		License: "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}

package app

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"strings"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/knf"
	"github.com/essentialkaos/ek/v12/log"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/usage"
	"github.com/essentialkaos/ek/v12/usage/completion/bash"
	"github.com/essentialkaos/ek/v12/usage/completion/fish"
	"github.com/essentialkaos/ek/v12/usage/completion/zsh"
	"github.com/essentialkaos/ek/v12/usage/man"
	"github.com/essentialkaos/ek/v12/usage/update"

	knfv "github.com/essentialkaos/ek/v12/knf/validators"
	knff "github.com/essentialkaos/ek/v12/knf/validators/fs"

	"github.com/essentialkaos/{{SHORT_NAME}}/app/support"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Basic application info
const (
	APP  = "{{NAME}}"
	VER  = "{{VERSION}}"
	DESC = "{{DESC}}"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Options
const (
	OPT_CONFIG   = "c:config"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// Configuration file properties
const (
	LOG_DIR   = "log:dir"
	LOG_FILE  = "log:file"
	LOG_PERMS = "log:perms"
	LOG_LEVEL = "log:level"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// optMap contains information about all supported options
var optMap = options.Map{
	OPT_CONFIG:   {Value: "/etc/{{SHORT_NAME}}.knf"},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.BOOL},

	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// useRawOutput is raw output flag
var useRawOutput = false

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main utility function
func Run(gitRev string, gomod []byte) {
	preConfigureUI()

	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		printError(errs[0].Error())
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print()
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Print(APP, VER, gitRev, gomod)
		os.Exit(0)
	case options.GetB(OPT_HELP) || len(args) == 0:
		genUsage().Print()
		os.Exit(0)
	}

	err := prepare()

	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}

	err = process(args)

	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}
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

	if !fsutil.IsCharacterDevice("/dev/stdout") && os.Getenv("FAKETTY") == "" {
		fmtc.DisableColors = true
		useRawOutput = true
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

// prepare prepares application to run
func prepare() error {
	err := knf.Global(options.GetS(OPT_CONFIG))

	if err != nil {
		return err
	}

	err = validateConfig()

	if err != nil {
		return err
	}

	err = setupLogger()

	if err != nil {
		return err
	}

	return nil
}

// validateConfig validates configuration file values
func validateConfig() error {
	errs := knf.Validate([]*knf.Validator{
		{LOG_DIR, knff.Perms, "DW"},
		{LOG_DIR, knff.Perms, "DX"},

		{LOG_LEVEL, knfv.NotContains, []string{"debug", "info", "warn", "error", "crit"}},
	})

	if len(errs) != 0 {
		return fmt.Errorf("Configuration file validation error: %w", errs[0])
	}

	return nil
}

// setupLogger confugures logger subsystem
func setupLogger() error {
	err := log.Set(knf.GetS(LOG_FILE), knf.GetM(LOG_PERMS, 644))

	if err != nil {
		return err
	}

	err = log.MinLevel(knf.GetS(LOG_LEVEL))

	if err != nil {
		return err
	}

	return nil
}

// process starts arguments processing
func process(args options.Arguments) error {
	return nil
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

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(info, "{{SHORT_NAME}}"))
	case "fish":
		fmt.Printf(fish.Generate(info, "{{SHORT_NAME}}"))
	case "zsh":
		fmt.Printf(zsh.Generate(info, optMap, "{{SHORT_NAME}}"))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(
		man.Generate(
			genUsage(),
			genAbout(""),
		),
	)
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
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2009,
		Owner:         "ESSENTIAL KAOS",
		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/{{SHORT_NAME}}", update.GitHubChecker},
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}

// ////////////////////////////////////////////////////////////////////////////////// //

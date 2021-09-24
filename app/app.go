package app

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2021 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pkg.re/essentialkaos/ek.v12/fmtc"
	"pkg.re/essentialkaos/ek.v12/fmtutil"
	"pkg.re/essentialkaos/ek.v12/fsutil"
	"pkg.re/essentialkaos/ek.v12/options"
	"pkg.re/essentialkaos/ek.v12/pluralize"
	"pkg.re/essentialkaos/ek.v12/terminal"
	"pkg.re/essentialkaos/ek.v12/usage"
	"pkg.re/essentialkaos/ek.v12/usage/completion/bash"
	"pkg.re/essentialkaos/ek.v12/usage/completion/fish"
	"pkg.re/essentialkaos/ek.v12/usage/completion/zsh"
	"pkg.re/essentialkaos/ek.v12/usage/man"
	"pkg.re/essentialkaos/ek.v12/usage/update"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "scratch"
	VER  = "0.0.6"
	DESC = "Utility for generating blank files for apps and services"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:      {Type: options.BOOL, Alias: "ver"},

	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// useRawOutput is raw output flag (for cli command)
var useRawOutput = false

// ////////////////////////////////////////////////////////////////////////////////// //

// Init is main app func
func Init() {
	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		for _, err := range errs {
			printError(err.Error())
		}

		os.Exit(1)
	}

	preConfigureUI()

	if options.Has(OPT_COMPLETION) {
		os.Exit(genCompletion())
	}

	if options.Has(OPT_GENERATE_MAN) {
		os.Exit(genMan())
	}

	configureUI()

	if options.GetB(OPT_VER) {
		os.Exit(showAbout())
	}

	if options.GetB(OPT_HELP) {
		os.Exit(showUsage())
	}

	if len(args) < 2 {
		listTemplates()
	} else {
		generateApp(args[0], args[1])
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

// configureUI configures UI
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	terminal.Prompt = "› "
}

// generateApp generates app from template
func generateApp(templateName, dir string) {
	err := checkTargetDir(dir)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	if !hasTemplate(templateName) {
		printErrorAndExit("There is no template with name \"%s\"", templateName)
	}

	dir, _ = filepath.Abs(dir)

	template, err := getTemplate(templateName)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	err = readVariablesValues(template.Vars)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	printVariablesInfo(template.Vars)

	fmtc.Println("{*}Generating files…{!}\n")

	err = copyTemplateData(template, dir)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	fmtc.Println("{g}Files successfully generated!{!}")
}

// listTemplates renders list of all available templates
func listTemplates() {
	templates, err := getTemplates()

	if err != nil {
		printErrorAndExit(err.Error())
	}

	if len(templates) == 0 {
		fmtc.Println("{y}No templates found{!}")
		return
	}

	fmtc.NewLine()

	for _, t := range templates {
		if len(t.Data) == 0 {
			fmtc.Printf(" {s}•{!} %s {s-}(empty){!}\n", t.Name)
		} else {
			fmtc.Printf(
				" {s}•{!} %s {s-}(%s){!}\n",
				t.Name, pluralize.P("%d %s", len(t.Data), "file", "files"),
			)
		}
	}

	fmtc.NewLine()
}

// readVariablesValues reads values for variables from template
func readVariablesValues(vars Variables) error {
	var curVar, totalVar int

	fmtc.NewLine()

	totalVar = vars.Count()

	for _, v := range knownVars.List {
		if !vars.Has(v) {
			continue
		}

		curVar++

		for {
			fmtc.Printf("{s-}[%d/%d]{!} {c}%s:{!}\n", curVar, totalVar, knownVars.Info[v].Desc)
			value, err := terminal.ReadUI("", true)

			fmtc.NewLine()

			if err != nil {
				os.Exit(1)
			}

			if !knownVars.Info[v].IsValid(value) {
				terminal.PrintWarnMessage("\"%s\" is not a valid value for this variable\n", value)
				continue
			}

			vars[v] = value

			break
		}
	}

	return nil
}

// printVariablesInfo prints defined variables
func printVariablesInfo(vars Variables) {
	fmtutil.Separator(false)

	for _, v := range knownVars.List {
		if !vars.Has(v) {
			continue
		}

		fmtc.Printf("  {*}%-16s{!} %s\n", v+":", vars[v])
	}

	fmtutil.Separator(false)

	fmtc.NewLine()

	ok, err := terminal.ReadAnswer("Everything is ok?", "y")

	fmtc.NewLine()

	if err != nil || !ok {
		os.Exit(1)
	}
}

// checkTargetDir checks target dir
func checkTargetDir(dir string) error {
	err := fsutil.ValidatePerms("DRWX", dir)

	if err != nil {
		return err
	}

	objects := fsutil.List(dir, false, fsutil.ListingFilter{
		NotMatchPatterns: []string{".git", ".github", "README.md", "LICENSE"},
	})

	if len(objects) != 0 {
		return fmt.Errorf("Target directory is not empty!")
	}

	return nil
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// printErrorAndExit prints error and exit with non-zero exit code
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showUsage prints usage info
func showUsage() int {
	genUsage().Render()
	return 0
}

// showAbout prints info about version
func showAbout() int {
	genAbout().Render()
	return 0
}

// genCompletion generates completion for different shells
func genCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(info, APP))
	case "fish":
		fmt.Printf(fish.Generate(info, APP))
	case "zsh":
		fmt.Printf(zsh.Generate(info, optMap, APP))
	default:
		return 1
	}

	return 0
}

// genMan generates man page
func genMan() int {
	fmt.Println(
		man.Generate(
			genUsage(),
			genAbout(),
		),
	)

	return 0
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "template", "dir")

	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample("package .", "Generate package blank files in current directory")
	info.AddExample("service $GOPATH/src/github.com/essentialkaos/myapp", "Generate service blank files in sources directory")

	return info
}

// genAbout generates info about version
func genAbout() *usage.About {
	return &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2006,
		Owner:         "ESSENTIAL KAOS",
		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/" + APP, update.GitHubChecker},
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

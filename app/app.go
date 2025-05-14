package app

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/essentialkaos/ek/v13/fmtc"
	"github.com/essentialkaos/ek/v13/fmtutil"
	"github.com/essentialkaos/ek/v13/fsutil"
	"github.com/essentialkaos/ek/v13/lscolors"
	"github.com/essentialkaos/ek/v13/options"
	"github.com/essentialkaos/ek/v13/path"
	"github.com/essentialkaos/ek/v13/pluralize"
	"github.com/essentialkaos/ek/v13/sortutil"
	"github.com/essentialkaos/ek/v13/support"
	"github.com/essentialkaos/ek/v13/support/apps"
	"github.com/essentialkaos/ek/v13/support/deps"
	"github.com/essentialkaos/ek/v13/system"
	"github.com/essentialkaos/ek/v13/terminal"
	"github.com/essentialkaos/ek/v13/terminal/input"
	"github.com/essentialkaos/ek/v13/terminal/tty"
	"github.com/essentialkaos/ek/v13/usage"
	"github.com/essentialkaos/ek/v13/usage/completion/bash"
	"github.com/essentialkaos/ek/v13/usage/completion/fish"
	"github.com/essentialkaos/ek/v13/usage/completion/zsh"
	"github.com/essentialkaos/ek/v13/usage/man"
	"github.com/essentialkaos/ek/v13/usage/update"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "scratch"
	VER  = "0.3.4"
	DESC = "Utility for generating blank files for apps and services"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.MIXED},

	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// templatesDir is path to directory with templates
var templatesDir string

// color tags for app name and version
var colorTagApp, colorTagVer string

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main utility function
func Run(gitRev string, gomod []byte) {
	preConfigureUI()

	args, errs := options.Parse(optMap)

	if !errs.IsEmpty() {
		terminal.Error("Options parsing errors:")
		terminal.Error(errs.Error(" - "))
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
		genAbout(gitRev).Print(options.GetS(OPT_VER))
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Collect(APP, VER).
			WithRevision(gitRev).
			WithDeps(deps.Extract(gomod)).
			WithApps(apps.Golang()).
			Print()
		os.Exit(0)
	case options.GetB(OPT_HELP):
		genUsage().Print()
		os.Exit(0)
	}

	if !findTemplatesDir() {
		os.Exit(1)
	}

	var err error

	switch len(args) {
	case 0:
		err = listTemplates()
	case 1:
		err = listTemplateData(args.Get(0).String())
	default:
		err = generateApp(
			args.Get(0).String(),
			args.Get(1).Clean().String(),
		)
	}

	if err != nil {
		terminal.Error(err)
		os.Exit(1)
	}
}

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if !tty.IsTTY() {
		fmtc.DisableColors = true
	}

	switch {
	case fmtc.IsTrueColorSupported():
		colorTagApp, colorTagVer = "{*}{&}{#F4BA33}", "{#F4BA33}"
	case fmtc.Is256ColorsSupported():
		colorTagApp, colorTagVer = "{*}{&}{#220}", "{#220}"
	default:
		colorTagApp, colorTagVer = "{*}{&}{y}", "{y}"
	}
}

// configureUI configures UI
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	input.Prompt = "› "
	input.NewLine = true
}

// findTemplatesDir tries to find directory with templates
func findTemplatesDir() bool {
	user, err := system.CurrentUser()

	if err != nil {
		terminal.Error("Can't get current user info: %v", err)
		return false
	}

	templatesDir = path.Clean(path.Join(user.HomeDir, ".config/scratch"))

	if !fsutil.IsExist(templatesDir) {
		terminal.Warn("▲ Can't find directory with templates")
		terminal.Warn("  Create directory ~/.config/scratch and add your templates to it")
		return false
	}

	err = fsutil.ValidatePerms("DRX", templatesDir)

	if err != nil {
		terminal.Error(err.Error())
		return false
	}

	return true
}

// generateApp generates app from template
func generateApp(templateName, dir string) error {
	err := checkTargetDir(dir)

	if err != nil {
		return err
	}

	if !hasTemplate(templateName) {
		return fmt.Errorf("There is no template with name %q", templateName)
	}

	dir, _ = filepath.Abs(dir)

	template, err := getTemplate(templateName)

	if err != nil {
		return err
	}

	err = readVariablesValues(template.Vars)

	if err != nil {
		return err
	}

	if !printVariablesInfo(template.Vars) {
		return nil
	}

	fmtc.Println("{*}Generating files…{!}\n")

	err = copyTemplateData(template, dir)

	if err != nil {
		return err
	}

	fmtc.Println("{g}Files successfully generated!{!}")

	return nil
}

// listTemplates renders list of all available templates
func listTemplates() error {
	templates, err := getTemplates()

	if err != nil {
		return err
	}

	if len(templates) == 0 {
		fmtc.Println("{y}No templates found{!}")
		return nil
	}

	fmtc.NewLine()

	for _, t := range templates {
		if len(t.Data) == 0 {
			fmtc.Printfn(" {s}•{!} %s {s-}(empty){!}", t.Name)
		} else {
			fmtc.Printfn(
				" {s}•{!} %s {s-}(%s){!}",
				t.Name, pluralize.P("%d %s", len(t.Data), "file", "files"),
			)
		}
	}

	fmtc.NewLine()

	return nil
}

// listTemplateData show list of files in template
func listTemplateData(name string) error {
	if !hasTemplate(name) {
		return fmt.Errorf("There is no templates with name %q", name)
	}

	t, err := getTemplate(name)

	if err != nil {
		return err
	}

	sortutil.StringsNatural(t.Data)

	fmtc.Printfn(
		"\n {s-}┌{!} {*}%s{!} {s-}(%s){!}\n {s-}│{!}",
		t.Name, pluralize.P("%d %s", len(t.Data), "file", "files"),
	)

	for i, file := range t.Data {
		if i+1 != len(t.Data) {
			fmtc.Print(" {s-}├─{!}")
		} else {
			fmtc.Print(" {s-}└─{!}")
		}

		fileSize := fsutil.GetSize(path.Join(t.Path, file))

		fmtc.Printfn(
			" %s {s-}(%s){!}",
			lscolors.ColorizePath(file),
			fmtutil.PrettySize(fileSize),
		)
	}

	fmtc.NewLine()

	for _, v := range knownVars.List {
		_, ok := t.Vars[v]

		if ok {
			fmtc.Printfn(" {s-}•{!} {s}%s — {&}%s{!}", v, knownVars.Info[v].Desc)
		}
	}

	fmtc.NewLine()

	return nil
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
			fmtc.Printfn("{s-}[%d/%d]{!} {c}%s:{!}", curVar, totalVar, knownVars.Info[v].Desc)
			value, err := input.Read("", input.NotEmpty)

			if err != nil {
				os.Exit(1)
			}

			if !knownVars.Info[v].IsValid(value) {
				terminal.Warn("%q is not a valid value for this variable\n", value)
				continue
			}

			vars[v] = value

			break
		}
	}

	return nil
}

// printVariablesInfo prints defined variables
func printVariablesInfo(vars Variables) bool {
	fmtutil.Separator(false)

	for _, v := range knownVars.List {
		if !vars.Has(v) {
			continue
		}

		fmtc.Printfn("  {*}%-16s{!} %s", v+":", vars[v])
	}

	fmtutil.Separator(false)

	fmtc.NewLine()

	ok, err := input.ReadAnswer("Everything is ok?", "y")

	fmtc.NewLine()

	return err == nil && ok
}

// checkTargetDir checks target dir
func checkTargetDir(dir string) error {
	if !fsutil.IsExist(dir) {
		return os.Mkdir(dir, 0755)
	}

	err := fsutil.ValidatePerms("DRWX", dir)

	if err != nil {
		return err
	}

	objects := fsutil.List(dir, false, fsutil.ListingFilter{
		NotMatchPatterns: []string{".git", ".github", "README.md", "LICENSE"},
	})

	if len(objects) != 0 {
		return fmt.Errorf("Target directory is not empty")
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Print(bash.Generate(info, APP))
	case "fish":
		fmt.Print(fish.Generate(info, APP))
	case "zsh":
		fmt.Print(zsh.Generate(info, optMap, APP))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(man.Generate(genUsage(), genAbout("")))
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "template", "target-dir")

	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample("package", "List files in template \"package\"")
	info.AddExample(
		"package .",
		"Generate files based on template \"package\" in current directory",
	)
	info.AddExample(
		"service $GOPATH/src/github.com/essentialkaos/myapp",
		"Generate files based on template \"service\" in given directory",
	)

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

		AppNameColorTag: colorTagApp,
		VersionColorTag: colorTagVer,
		DescSeparator:   "{s}—{!}",

		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		BugTracker:    "https://github.com/essentialkaos/scratch/issues",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/scratch", update.GitHubChecker},
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}

// ////////////////////////////////////////////////////////////////////////////////// //

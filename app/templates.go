package app

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/sortutil"
	"github.com/essentialkaos/ek/v12/timeutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const SRC_DIR = "github.com/essentialkaos/scratch"

const (
	VAR_NAME        = "NAME"
	VAR_SHORT_NAME  = "SHORT_NAME"
	VAR_VERSION     = "VERSION"
	VAR_DESC        = "DESC"
	VAR_DESC_README = "DESC_README"

	VAR_CODEBEAT_UUID = "CODEBEAT_UUID"

	VAR_SHORT_NAME_TITLE    = "SHORT_NAME_TITLE"
	VAR_SHORT_NAME_LOWER    = "SHORT_NAME_LOWER"
	VAR_SHORT_NAME_UPPER    = "SHORT_NAME_UPPER"
	VAR_SPEC_CHANGELOG_DATE = "SPEC_CHANGELOG_DATE"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type Variables map[string]string // name → value

type Template struct {
	Name string // Name of tempate
	Path string // Path to directory with template data

	Vars Variables // Variables
	Data []string  // List of files and directories of template
}

type VariableInfoStore struct {
	Info map[string]VariableInfo
	List []string
}

type VariableInfo struct {
	Desc      string
	Validator string
	IsDynamic bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

var knownVars = &VariableInfoStore{
	// Info contains info about all supported variables
	Info: map[string]VariableInfo{
		VAR_NAME:        {"Name", `^[a-zA-Z0-9\_\-]{2,32}$`, false},
		VAR_SHORT_NAME:  {"Short name (binary name or repository name)", `^[a-z0-9\_\-]{2,32}$`, false},
		VAR_VERSION:     {"Version (in semver notation)", `^[0-9]+\.[0-9]*\.?[0-9]*$`, false},
		VAR_DESC:        {"Description", `^.{16,128}$`, false},
		VAR_DESC_README: {"Description for README file (part after 'app is… ')", `^.{16,128}$`, false},

		VAR_CODEBEAT_UUID: {"Codebeat project UUID", ``, false},

		VAR_SHORT_NAME_TITLE:    {"Short name in title case", ``, true},
		VAR_SHORT_NAME_LOWER:    {"Short name in lower case", ``, true},
		VAR_SHORT_NAME_UPPER:    {"Short name in upper case", ``, true},
		VAR_SPEC_CHANGELOG_DATE: {"Date in spec changelog", ``, true},
	},

	// List contains variables which requires user input in particular order
	List: []string{
		VAR_NAME,
		VAR_SHORT_NAME,
		VAR_VERSION,
		VAR_DESC,
		VAR_DESC_README,
		VAR_CODEBEAT_UUID,
	},
}

var varRegex = regexp.MustCompile(`\{\{([A-Z0-9_]+)\}\}`)

// ////////////////////////////////////////////////////////////////////////////////// //

// Has returns true if map contains variable with given name
func (v Variables) Has(name string) bool {
	_, ok := v[name]
	return ok
}

// Count returns number of variables which requires user input
func (v Variables) Count() int {
	var result int

	for n := range v {
		if !knownVars.Info[n].IsDynamic {
			result++
		}
	}

	return result
}

// IsValid validates value
func (vi VariableInfo) IsValid(value string) bool {
	if vi.Validator == "" {
		return true
	}

	return regexp.MustCompile(vi.Validator).MatchString(value)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getTemplates returns slice with info about all available templates
func getTemplates() ([]*Template, error) {
	templatesDir, err := getTemplatesDir()

	if err != nil {
		return nil, err
	}

	var result []*Template

	templates := fsutil.List(
		templatesDir, true,
		fsutil.ListingFilter{Perms: "DRX"},
	)

	if len(templates) == 0 {
		return result, nil
	}

	sortutil.StringsNatural(templates)

	for _, templateName := range templates {
		template, err := getTemplate(templateName)

		if err != nil {
			return nil, fmt.Errorf("Problem with template \"%s\": %w", templateName, err)
		}

		result = append(result, template)
	}

	return result, nil
}

// hasTemplate returns true if template with given name is present
func hasTemplate(templateName string) bool {
	templates, err := getTemplates()

	if err != nil {
		return false
	}

	for _, t := range templates {
		if t.Name == templateName {
			return true
		}
	}

	return false
}

// getTemplatesDir returns path to directory with templates
func getTemplatesDir() (string, error) {
	gopath := os.Getenv("GOPATH")
	srcDir := gopath + "/src/" + SRC_DIR

	if !fsutil.IsExist(srcDir) {
		return "", fmt.Errorf("Can't find directory with scratch sources")
	}

	templatesDir := srcDir + "/templates"

	return templatesDir, fsutil.ValidatePerms("DRX", templatesDir)
}

// getTemplate returns list of all files and directories in template
func getTemplate(templateName string) (*Template, error) {
	templatesDir, err := getTemplatesDir()

	if err != nil {
		return nil, err
	}

	templateDir := templatesDir + "/" + templateName

	if !fsutil.IsExist(templateDir) {
		return nil, fmt.Errorf("Can't find template with name \"%s\"", templateName)
	}

	files := fsutil.ListAllFiles(templateDir, false)
	vars, err := extractVariables(templateDir, files)

	if err != nil {
		return nil, err
	}

	return &Template{
		Name: templateName,
		Path: templateDir,
		Vars: vars,
		Data: files,
	}, nil
}

// copyTemplate copies all files from template and applies variables
func copyTemplateData(tmpl *Template, targetDir string) error {
	var err error

	applyDynamicVariables(tmpl.Vars)

	for _, file := range tmpl.Data {
		sourceFile := tmpl.Path + "/" + file
		targetFile := targetDir + "/" + formatFileName(file, tmpl.Vars)
		targetFileDir := path.Dir(targetFile)

		if !fsutil.IsExist(targetFileDir) {
			err = os.MkdirAll(targetFileDir, 0755)

			if err != nil {
				return err
			}
		}

		err = copyTemplateFile(sourceFile, targetFile, tmpl.Vars)

		if err != nil {
			return err
		}
	}

	return nil
}

// copyTemplateFile copies file and applies variables
func copyTemplateFile(sourceFile, targetFile string, vars Variables) error {
	sfd, err := os.OpenFile(sourceFile, os.O_RDONLY, 0)

	if err != nil {
		return err
	}

	defer sfd.Close()

	tfd, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		return err
	}

	defer tfd.Close()

	r := bufio.NewReader(sfd)
	s := bufio.NewScanner(r)
	w := bufio.NewWriter(tfd)

	return writeTemplateData(s, w, vars)
}

// writeTemplateData writes template data
func writeTemplateData(s *bufio.Scanner, w *bufio.Writer, vars Variables) error {
	for s.Scan() {
		line := s.Text()

		if varRegex.MatchString(line) {
			for _, fn := range varRegex.FindAllStringSubmatch(line, -1) {
				varDef, varValue := fn[0], vars[fn[1]]
				line = strings.ReplaceAll(line, varDef, varValue)
			}
		}

		_, err := w.WriteString(line + "\n")

		if err != nil {
			return err
		}
	}

	return w.Flush()
}

// extractVariables extracts all unique variables from all files in template
func extractVariables(dir string, files []string) (Variables, error) {
	vars := make(Variables)

	for _, dataFile := range files {
		dataFilePath := dir + "/" + dataFile
		fileVars, err := scanFileForVariables(dataFilePath)

		if err != nil {
			return nil, err
		}

		if len(fileVars) == 0 {
			continue
		}

		for _, fileVar := range fileVars {
			vars[fileVar] = ""
		}
	}

	return vars, validateVariables(vars)
}

// scanFileForVariables scans given file for variables
func scanFileForVariables(file string) ([]string, error) {
	fd, err := os.OpenFile(file, os.O_RDONLY, 0)

	if err != nil {
		return nil, err
	}

	defer fd.Close()

	var result []string

	r := bufio.NewReader(fd)
	s := bufio.NewScanner(r)

	for s.Scan() {
		line := s.Text()

		if !varRegex.MatchString(line) {
			continue
		}

		fn := varRegex.FindAllStringSubmatch(line, -1)

		for _, fns := range fn {
			result = append(result, fns[1])
		}
	}

	return result, nil
}

// applyDynamicVariables generates values for dynamic variables
func applyDynamicVariables(vars Variables) {
	for v := range vars {
		switch v {
		case VAR_SPEC_CHANGELOG_DATE:
			vars[v] = timeutil.Format(time.Now(), "%a %b %d %Y")

		case VAR_SHORT_NAME_TITLE:
			vars[v] = strings.Title(vars[VAR_SHORT_NAME])

		case VAR_SHORT_NAME_LOWER:
			vars[v] = strings.ToLower(vars[VAR_SHORT_NAME])

		case VAR_SHORT_NAME_UPPER:
			vars[v] = strings.ToUpper(vars[VAR_SHORT_NAME])
		}
	}
}

// validateVariables validates variable
func validateVariables(vars Variables) error {
	for v := range vars {
		_, ok := knownVars.Info[v]

		if !ok {
			return fmt.Errorf("Template contains unknown variable \"%s\"", v)
		}
	}

	return nil
}

// formatFileName formats file name
func formatFileName(name string, vars Variables) string {
	switch {
	case strings.Contains(name, "_name_"):
		name = strings.ReplaceAll(name, "_name_", vars["SHORT_NAME"])
	}

	return name
}

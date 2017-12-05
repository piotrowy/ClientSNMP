package mib

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"fmt"
)

const (
	mibDir     = "/mibs/"
	currentDir = "./"
)

const (
	imports_section_regex        = "(?s)IMPORTS.*?;"             //single line
	imports_groups_regex         = ".*?FROM.*?[\n;]"             //single line
	import_elements_regex        = ".*FROM"                      //single line
	object_identifier_line_regex = "^.*OBJECT IDENTIFIER ::=.*$" //single line
	object_identifier_data_regex = "{.*}"
	object_type_regex            = "\\S*\\s*OBJECT-TYPE.*?::= {.*?}" //single line
	syntax_regex                 = "SYNTAX.*"
	access_regex                 = "ACCESS.*"
	status_regex                 = "STATUS.*"
	description_regex            = "\".*\"" //single line
	class_data_regex             = "::= {.*"
)

//returns true if regex is found
type fnMibRegEx func(s string) bool

func Parse(t *Tree, mib string) bool {
	var importedElements []string

	getElements := func(data string) bool {
		importedElements = []string{}
		return false
	}

	return ParseFn(mib, func(data string) bool {
		importSection := regexp.MustCompile(imports_section_regex).FindString(data)
		importGroups := regexp.MustCompile(imports_groups_regex).FindString(importSection)
		getElements(importGroups)

		oidLines := regexp.MustCompile(object_identifier_data_regex).FindStringSubmatch(data)
		fmt.Print(oidLines)
		return false
	})
}

func ParseFn(mib string, fn fnMibRegEx) bool {
	projectDir, err := filepath.Abs(currentDir)
	if err != nil {
		panic(err)
	}
	data := readFile(projectDir + mibDir + mib)
	return fn(data)
}

func readFile(f string) string {
	dat, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}
	return string(dat)
}

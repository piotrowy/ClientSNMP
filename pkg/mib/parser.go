package mib

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"fmt"
)

const (
	mibDir     = "/mibs/"
	mibList    = "/pkg/mib/MIBLIST"
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
	if !fileExist(mib) {
		return false
	}

	return parse(mib, func(data string) bool {
		importSection := regexp.MustCompile(imports_section_regex).FindString(data)
		importGroups := regexp.MustCompile(imports_groups_regex).FindString(importSection)
		fmt.Print(importGroups)
		return false

	})
}

func parse(name string, fn fnMibRegEx) bool {
	projectDir, err := filepath.Abs(currentDir)
	if err != nil {
		return false
	}
	data := readFile(projectDir + mibDir + name)
	return fn(data)
}

func readFile(f string) string {
	dat, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}
	return string(dat)
}

func fileExist(name string) bool {
	if p, err := filepath.Abs(currentDir); err == nil {
		if b, err := regexp.MatchString(name, readFile(p+mibList)); err == nil {
			return b
		}
	}
	return false
}

package mib

import (
	"io/ioutil"
)

const (
	MIBDIR       = "../../mibs/"
	MIBLIST_FILE = "./MIBLIST"
)

const (
	IMPORTS_SECTION_REGEX        = "IMPORTS.*?;"                 //single line
	IMPORTS_GROUPS_REGEX         = ".*?FROM.*?[\n;]"             //single line
	IMPORT_ELEMENTS_REGEX        = ".*FROM"                      //single line
	OBJECT_IDENTIFIER_LINE_REGEX = "^.*OBJECT IDENTIFIER ::=.*$" //single line
	OBJECT_IDENTIFIER_DATA_REGEX = "{.*}"
	OBJECT_TYPE_REGEX            = "\\S*\\s*OBJECT-TYPE.*?::= {.*?}" //single line
	SYNTAX_REGEX                 = "SYNTAX.*"
	ACCESS_REGEX                 = "ACCESS.*"
	STATUS_REGEX                 = "STATUS.*"
	DESCRIPTION_REGEX            = "\".*\"" //single line
	CLASS_DATA_REGEX             = "::= {.*"
)

//returns true if regex is found
type fnMibRegEx func(s string) bool

func Parse(t *Tree, mib string) {

}

func ParseMibRegEx(path string, fn fnMibRegEx) bool {
	data := readFile(path)
	return fn(data)
}

func readFile(f string) string {
	if dat, err := ioutil.ReadFile(f); err != nil {
		panic(err)
	} else {
		return string(dat)
	}
}

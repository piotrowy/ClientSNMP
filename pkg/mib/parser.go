package mib

import (
	"regexp"
	"strings"
	"strconv"
)

const (
	all   = -1
	space = " "
	empty = ""
)

var (
	importsSectionRegex       = regexp.MustCompile("(?s)IMPORTS.*?;")
	importsGroupsRegex        = regexp.MustCompile("(?s).*?FROM.*?[\\n;]")
	importElementsRegex       = regexp.MustCompile("(?s).*FROM")
	objectIdentifierLineRegex = regexp.MustCompile("(?m)^.*OBJECT IDENTIFIER ::=.*$")
	objectIdentifierDataRegex = regexp.MustCompile("{.*}")
	objectTypeRegex           = regexp.MustCompile("(?s)\\S*\\s*OBJECT-TYPE.*?SYNTAX.*?ACCESS.*?.*?::= {.*?}")
	syntaxRegex               = regexp.MustCompile("SYNTAX.*")
	accessRegex               = regexp.MustCompile("ACCESS.*")
	statusRegex               = regexp.MustCompile("STATUS.*")
	descriptionRegex          = regexp.MustCompile("(?s)\".*\"")
	classDataRegex            = regexp.MustCompile("::= {.*")
	commentRegex              = regexp.MustCompile("(?s)--.*?\\n")
)

type emptyStruct struct {

}

func Parse(mib string) (*Tree, error) {
	var (
		oids        oids
		objectTypes objectTypes
		dataTypes   dataTypes
		elements    map[string][]string
	)

	parser, err := GetParser(mib)
	if err != nil {
		return &Tree{}, err
	}

	parser = parser.Map(func(data string) string {
		for _, v := range commentRegex.FindAllString(data, all) {
			data = strings.Replace(data, v, empty, all)
		}
		return data
	})

	parser.Parse(func(data string) {
		elements = getElements(data)
		oids = getOids(data, elements)
		objectTypes = getObjectTypes(data)
		dataTypes = getDataTypes(data, objectTypes, elements)
	})

	return createTree(oids, objectTypes, dataTypes)
}

func getElements(data string) map[string][]string {
	elements := make(map[string][]string, 0)
	importSection := importsSectionRegex.FindString(data)
	importGroups := importsGroupsRegex.FindAllString(importSection, all)
	for _, v := range importGroups {
		res := elementsFromGroupString(v)
		if len(res) > 0 {
			elements[getImportFileName(v)] = res
		}
	}
	return elements
}

func getOids(data string, elements map[string][]string) (o oids) {
	o = oidsFromString(data)
	for k := range elements {
		o = append(o, oidsFromFile(k)...)
	}
	return
}

func getObjectTypes(data string) (ots objectTypes) {
	withoutImports := data[importsSectionRegex.FindStringIndex(data)[1]:]
	objectTypeLines := objectTypeRegex.FindAllString(withoutImports, all)
	for _, v := range objectTypeLines {
		ots = append(ots, objectTypeLine(v).toObjectType())
	}
	return
}

func getDataTypes(data string, objectTypes objectTypes, elements map[string][]string) (dts dataTypes) {
	for file, v := range elements {
		dts = append(dts, dataTypesFromFile(file, v)...)
	}
	for _, v := range objectTypes {
		if t, ok := parseDataType(data, v.Syntax); ok {
			dts = append(dts, t)
		}
	}
	return
}

func createTree(o oids, ot objectTypes, dt dataTypes) (*Tree, error) {
	dtMap := make(map[string]DataType)
	for _, v := range dt {
		dtMap[v.Name] = v
	}

	o, root, err := o.next()
	if err != nil {
		return &Tree{}, err
	}
	t := New(root, ObjectType{})

	var oid Oid
	for ; err == nil; o, oid, err = o.next() {
		t.InsertOid(oid)
	}
	err = nil

	var objectType ObjectType
	for ; err == nil; ot, objectType, err = ot.next() {
		if v, ok := dtMap[objectType.Syntax]; ok {
			t.distinctTypes[v.BaseType] = emptyStruct{}
			objectType.Min, objectType.Max = restrictionsFromSyntax(v.Restrictions)

		} else {
			t.distinctTypes[objectType.Syntax] = emptyStruct{}
		}
		t.InsertObjectType(objectType)
	}

	return t, nil
}

func dataTypesFromFile(file string, types []string) (dts dataTypes) {
	parser, err := GetParser(file)
	if err != nil {
		return
	}

	parser.Parse(func(data string) {
		for _, v := range types {
			if t, ok := parseDataType(data, v); ok {
				dts = append(dts, t)
			}
		}
	})
	return
}

func parseDataType(data, t string) (dt DataType, ok bool) {
	match := regexp.MustCompile("(?s)" + t + "\\s*::=.*?\n\n").FindString(data)
	if match != empty {
		dt, ok = createDataType(match), true
	}
	return
}

func createDataType(dataTypeLine string) DataType {
	baseTypeMatch := regexp.MustCompile("\\n.*?[({]").FindString(dataTypeLine)
	baseType := strings.TrimSpace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(baseTypeMatch,
		"{", empty, all), "(", empty, all), "IMPLICIT", empty, all), "EXPLICIT", empty, all))
	name := strings.TrimSpace(strings.Split(dataTypeLine, space)[0])
	restrictionsMatch := regexp.MustCompile("[{(].*[})]").FindString(dataTypeLine)
	restrictions := strings.TrimSpace(strings.Trim(restrictionsMatch, "{()}"))
	codingValue, _ := getCodingValue(dataTypeLine)
	return DataType{
		Name:         name,
		BaseType:     baseType,
		Restrictions: restrictions,
		CodingType:   getCodingType(dataTypeLine),
		CodingValue:  codingValue,
	}
}

func getCodingValue(dataTypeLine string) (r int, err error) {
	codingMatch := regexp.MustCompile("\\[.*?\\]").FindString(dataTypeLine)
	splitted := strings.Split(strings.Trim(codingMatch, "[]"), space)
	return strconv.Atoi(splitted[len(splitted)-1])
}

func getCodingType(dataTypeLine string) (r string) {
	const (
		EXPLICIT = "EXPLICIT"
		IMPLICIT = "IMPLICIT"
	)

	if strings.Contains(dataTypeLine, EXPLICIT) {
		r = EXPLICIT
	} else if strings.Contains(dataTypeLine, IMPLICIT) {
		r = IMPLICIT
	} else {
		r = empty
	}
	return
}

func getImportFileName(v string) string {
	trimmed := strings.TrimSpace(v)
	splitted := strings.Split(trimmed, space)
	dirtyFileName := splitted[len(splitted)-1]
	return strings.Replace(dirtyFileName, ";", empty, -1)
}

func elementsFromGroupString(group string) (elements []string) {
	elementLine := importElementsRegex.FindString(group)
	groups := strings.Split(elementLine, ",")
	for _, v := range groups {
		for _, v2 := range []string{"FROM", "IMPORTS", "OBJECT_TYPE"} {
			v = strings.Replace(v, v2, empty, -1)
		}
		v = strings.TrimSpace(v)
		if v != empty {
			elements = append(elements, v)
		}
	}
	return
}

type objectTypeLine string

func (o objectTypeLine) toObjectType() ObjectType {
	trimmed := strings.TrimSpace(objectTypeRegex.FindString(string(o)))
	splitted := strings.Split(trimmed, space)

	classData := strings.Split(objectIdentifierDataRegex.FindString(classDataRegex.FindString(string(o))), space)
	num, _ := strconv.Atoi(classData[2])
	syntax := strings.TrimSpace(strings.Replace(syntaxRegex.FindString(string(o)), "SYNTAX", empty, all))
	min, max := restrictionsFromSyntax(syntax)

	return ObjectType{
		Name:        splitted[0],
		Syntax:      syntax,
		Access:      strings.Replace(accessRegex.FindString(trimmed), "ACCESS", empty, all),
		Status:      strings.Replace(statusRegex.FindString(trimmed), "STATUS", empty, all),
		Description: strings.Trim(descriptionRegex.FindString(string(o)), "\""),
		Class:       classData[1],
		Number:      num,
		Min:         min,
		Max:         max,
	}
}

func restrictionsFromSyntax(syntaxLine string) (min uint64, max uint64) {
	rangeMatch := regexp.MustCompile("[0-9]*\\.\\.[0-9]*").FindString(syntaxLine)
	if rangeMatch != empty {
		splitted := strings.Split(rangeMatch, ".")
		min, _ = strconv.ParseUint(splitted[0], 10, 64)
		max, _ = strconv.ParseUint(splitted[1], 10, 64)
	}
	return
}

func oidsFromFile(file string) (o oids) {
	parser, err := GetParser(file)
	if err != nil {
		return
	}
	parser.Parse(func(data string) {
		o = oidsFromString(data)
	})
	return
}

func oidsFromString(data string) (o oids) {
	for _, v := range objectIdentifierLineRegex.FindAllString(data, all) {
		o = append(o, oidLine(v).toOid())
	}
	return
}

type oidLine string

func (o oidLine) toOid() Oid {
	trimmed := strings.TrimSpace(objectIdentifierDataRegex.FindString(string(o)))
	splitted := strings.Split(trimmed, space)

	num, _ := strconv.Atoi(splitted[len(splitted)-2])
	return Oid{
		Name:   strings.Split(strings.Trim(string(o), "- "), space)[0],
		Class:  splitted[1],
		Number: num,
	}
}

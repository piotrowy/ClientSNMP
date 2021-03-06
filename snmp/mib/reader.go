package mib

import (
	"io/ioutil"
	"path/filepath"
)

const (
	mibDir  = "/mibs/"
	mainDir = "./"
)

type ParsingFn func(s string)

type Parser struct {
	data string
}

func (p Parser) Parse(f ParsingFn) {
	f(p.data)
}

func (p Parser) Map(f func(s string) string) Parser {
	return Parser{
		data: f(p.data),
	}
}

func GetParser(mib string) (Parser, error) {
	data, err := Read(mib)
	if err != nil {
		return Parser{}, err
	}
	return Parser{
		data: data,
	}, nil
}

func Read(f string) (string, error) {
	path, err := filepath.Abs(mainDir)
	if err != nil {
		return "", err
	}
	return read(path + mibDir + f)
}

func read(f string) (string, error) {
	bytes, err := ioutil.ReadFile(f)
	return string(bytes), err
}

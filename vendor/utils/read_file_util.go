package utils

import (
	"bufio"
)

type Reader struct {
	reader  *bufio.Reader
	line    int
	lineStr string
}

func NewReader(r *bufio.Reader) Reader {
	return Reader{reader: r}
}

func (r *Reader) ReadLine() (line int, lineStr string, e error) {
	l, _, e := r.reader.ReadLine()
	if l == nil {
		return -1, EMPTY_STRING, nil
	}
	r.line++
	r.lineStr = string(l)
	//fmt.Println(r.line, r.lineStr)
	return r.line, r.lineStr, e
}

func (r *Reader) CurrentLine() (line int, lineStr string) {
	return r.line, r.lineStr
}

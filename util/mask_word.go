package util

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type MaskWordString map[string]string

var MaskWord MaskWordString

func init() {
	MaskWord = make(MaskWordString)
	MaskWord.loadMaskWord("")
}

func (m MaskWordString) loadMaskWord(filename string) {
	fp, err := os.Open(filename)
	if err != nil {
		return
	}
	defer fp.Close()

	br := bufio.NewReader(fp)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		line := string(a)
		line = strings.Replace(line, "\r", "", -1)
		line = strings.Replace(line, "\n", "", -1)
		s := strings.Fields(line)
		if len(s) == 3 {
			m[s[1]] = s[2]
		}
	}
}

//屏蔽字替换
func ReplaceMaskWord(s string) string {
	for k, v := range MaskWord {
		s = strings.Replace(s, k, v, -1)
	}
	return s
}

//是否含有屏蔽字
func HaveMaskWord(s string) bool {
	for k := range MaskWord {
		if strings.Contains(s, k) {
			return true
		}
	}
	return false
}

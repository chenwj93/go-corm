package struct_utils

import (
	"regexp"
	"strings"
	"utils"
)

var reg =`delete[\s]+from[\s]+([-_a-zA-Z0-9]+)[\s]+|delete[\s]+([-_a-zA-Z0-9]+)[\s]+|update[\s]+([-_a-zA-Z0-9]+)[\s]+|insert[\s]+into[\s]+([-_a-zA-Z0-9]+)[\s]+|insert[\s]+([-_a-zA-Z0-9]+)[\s]+`

var r *regexp.Regexp

func init() {
	r = regexp.MustCompile(reg)
}

func SelectTable(sql string) string{
	index := r.FindStringSubmatchIndex(strings.ToLower(sql))
	if len(index) == 0 {
		return utils.EMPTY_STRING
	}
	//fmt.Println(index)
	for i := 2; i < len(index); i += 2 {
		if index[i] != -1 && index[i+1] != -1{
			return sql[index[i]:index[i+1]]
		}
	}
	return utils.EMPTY_STRING
}
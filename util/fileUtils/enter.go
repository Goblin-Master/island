package fileUtils

import (
	"errors"
	"strings"
)

func SuffixJudge(fileName string) (suffix string, err error) {
	_list := strings.Split(fileName, ".")
	if len(_list) > 1 {
		suffix = _list[len(_list)-1]
	} else {
		err = errors.New("文件名错误")
	}
	return
}

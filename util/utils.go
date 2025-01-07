package util

import (
	"encoding/json"
	"path/filepath"
	"runtime"
	"tgwp/log/zlog"
)

/*
GetRootPath 搜索项目的文件根目录, 并和 myPath 拼接起来
*/
func GetRootPath(myPath string) string {
	_, fileName, _, ok := runtime.Caller(0)
	if !ok {
		panic("Something wrong with getting root path")
	}
	absPath, err := filepath.Abs(fileName)
	rootPath := filepath.Dir(filepath.Dir(absPath))
	if err != nil {
		panic(any(err))
	}
	return filepath.Join(rootPath, myPath)
}

// StructToMap
//
//	@Description: struct to map
//	@param value
//	@return map[string]interface{}
func StructToMap(value interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	resJson, err := json.Marshal(value)
	if err != nil {
		zlog.Errorf("Json Marshal failed ,msg: %s", err.Error())
		return nil
	}
	err = json.Unmarshal(resJson, &m)
	if err != nil {
		zlog.Errorf("Json Unmarshal failed,msg : %s", err.Error())
		return nil
	}
	return m
}

// StuctToJson
//
//	@Description: struct to json
//	@param value
//	@return string
//	@return error
func StuctToJson(value interface{}) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(data), err
}

// JsonToStruct
//
//	@Description: json to struct
//	@param str
//	@param value
//	@return error
func JsonToStruct(str string, value interface{}) error {
	return json.Unmarshal([]byte(str), value)
}

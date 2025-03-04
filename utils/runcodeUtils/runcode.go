package runcodeUtils

import (
	"encoding/json"
	"github.com/levigross/grequests"
	"strings"
	"tgwp/log/zlog"
	"tgwp/response"
)

const (
	CPP_VERSION    = "10.2.0"
	PYTHON_VERSION = "3.10.0"
)

type PistonResp struct {
	Language string `json:"language"`
	Version  string `json:"version"`
	Run      struct {
		Output string `json:"output"`
		Code   int    `json:"code"`
		Signal string `json:"signal"`
	} `json:"run"`
}

// RunCode 通过Piston API运行代码
func RunCode(language string, code string, stdin string, time_limit int, memory_limit int) (output string, is_limit_out bool, err error) {
	// 选择编程语言和版本
	version, fileName := "", ""
	switch language {
	case "cpp":
		version = CPP_VERSION
		fileName = "main.cpp"
	case "python":
		version = PYTHON_VERSION
		fileName = "main.py"
	default:
		err = response.ErrResp(err, response.LANGUAGE_NOT_EXIST)
		return
	}
	// 构造请求数据
	postData := map[string]interface{}{
		"language": language,
		"version":  version,
		"files": []map[string]interface{}{
			{
				"name":    fileName,
				"content": code,
			},
		},
		"stdin":            stdin,
		"run_cpu_time":     time_limit,
		"run_memory_limit": memory_limit,
	}
	zlog.Debugf("Piston API request data: %v", postData)
	// 在构造 postData 后添加以下代码
	jsonData, _ := json.Marshal(postData)
	zlog.Debugf("Real Request JSON: %s", string(jsonData))
	// 发送请求
	geq := &grequests.RequestOptions{
		JSON: postData,
	}
	resp, err := grequests.Post("https://emkc.org/api/v2/piston/execute", geq)
	if err != nil {
		zlog.Errorf("Post to Piston API failed: %v", err)
		err = response.ErrResp(err, response.PISTON_API_ERROR)
		return
	}
	zlog.Debugf("Piston API response: %v", resp.String())
	// 解析返回结果
	var pistonResp PistonResp
	if err = json.Unmarshal([]byte(resp.String()), &pistonResp); err != nil {
		zlog.Errorf("Unable to parse JSON response: %v", err)
		return
	}
	// 处理返回结果
	if len(pistonResp.Run.Signal) > 0 {
		is_limit_out = true
	} else {
		output = pistonResp.Run.Output
	}

	return
}

// CompareOutput 比较运行结果与答案是否相同
func CompareOutput(output, answer string) bool {
	// 一行一行比较，忽略行尾空格
	Out := strings.Split(output, "\n")
	Ans := strings.Split(answer, "\n")
	// 去除结尾空行
	for len(Out) > 0 && Out[len(Out)-1] == "" {
		Out = Out[:len(Out)-1]
	}
	for len(Ans) > 0 && Ans[len(Ans)-1] == "" {
		Ans = Ans[:len(Ans)-1]
	}
	// 行数不同直接返回 false
	if len(Out) != len(Ans) {
		return false
	}
	// 开始逐行比较
	for i := 0; i < len(Out); i++ {
		// 忽略行尾空格
		Out[i] = strings.TrimRight(Out[i], " ")
		Ans[i] = strings.TrimRight(Ans[i], " ")
		// 比较行内容
		if Out[i] != Ans[i] {
			return false
		}
	}
	// 全部相同
	return true
}

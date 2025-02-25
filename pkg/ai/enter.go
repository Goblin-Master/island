package ai

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"tgwp/global"
	"tgwp/log/zlog"
	"tgwp/utils/cacheUtils"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type ChatRequest struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
	Stream   bool      `json:"stream"`
}
type ChatResponse struct {
	Id      string `json:"id"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Object  string `json:"object"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

const baseUrl = "https://api.chatanywhere.tech/v1/chat/completions"

//go:embed chat.prompt
var chatPrompt string

func BaseRequest(ctx context.Context, r ChatRequest) (res *http.Response, err error) {
	method := "POST"
	byteData, err := json.Marshal(r)
	if err != nil {
		zlog.CtxErrorf(ctx, "json序列化失败 %s", err)
		return
	}
	req, err := http.NewRequest(method, baseUrl, bytes.NewBuffer(byteData))
	if err != nil {
		zlog.CtxErrorf(ctx, "ai请求参数格式化失败 %s", err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", global.Config.AI.ApiKey))
	req.Header.Add("Content-Type", "application/json")

	// 设置http客户端,优化程序
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConnsPerHost = 100
	httpClient := &http.Client{
		Timeout:   100 * time.Second,
		Transport: t,
	}

	res, err = httpClient.Do(req)
	if err != nil {
		zlog.CtxErrorf(ctx, "请求ai失败 %s", err)
		return
	}
	return
}
func Chat(ctx context.Context, content string) (msg string, err error) {
	r := ChatRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: chatPrompt,
			},
			{
				Role:    "user",
				Content: content,
			},
		},
		Model:  global.Config.AI.Model,
		Stream: false,
	}
	res, err := BaseRequest(ctx, r)
	if err != nil {
		return
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			zlog.CtxErrorf(ctx, "关闭响应体失败 %s", err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		zlog.CtxErrorf(ctx, "读取ai响应失败 %s", err)
		return
	}
	var response ChatResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		zlog.CtxErrorf(ctx, "json反序列化失败 %s", err)
		return
	}
	if len(response.Choices) > 0 {
		msg = response.Choices[0].Message.Content
		return
	}
	zlog.CtxErrorf(ctx, "ai响应失败 %s", string(body))
	return
}

//go:embed chat_stream.prompt
var chatStreamPrompt string

type Choice struct {
	Index int `json:"index"`
	Delta struct {
		Content string `json:"content"`
	} `json:"delta"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason interface{} `json:"finish_reason"`
}
type StreamData struct {
	Id                string   `json:"id"`
	Choices           []Choice `json:"choices"`
	Created           int      `json:"created"`
	Model             string   `json:"model"`
	Object            string   `json:"object"`
	SystemFingerprint string   `json:"system_fingerprint"`
}

func ChatStream(ctx context.Context, content string) (msgChan chan string, err error) {
	msgChan = make(chan string)
	// TODO: 获取用户id
	userid := int64(793478004095)
	history, err := cacheUtils.GetContent(ctx, strconv.FormatInt(userid, 10))
	if err != nil {
		zlog.CtxErrorf(ctx, "获取历史记录失败 %s", err)
		return
	}
	r := ChatRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: chatStreamPrompt,
			},
			{
				Role:    "user",
				Content: content,
			},
			{
				Role:    "assistant",
				Content: history,
			},
		},
		Model:  global.Config.AI.Model,
		Stream: true,
	}
	res, err := BaseRequest(ctx, r)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(res.Body)
	// 设置分割函数,按行分割
	scanner.Split(bufio.ScanLines)
	var reply string
	go func() {
		defer close(msgChan) // 确保通道在协程结束时关闭
		defer func() {
			if err = res.Body.Close(); err != nil {
				zlog.CtxErrorf(ctx, "关闭响应体失败 %s", err)
			}
		}()
		for scanner.Scan() {
			text := scanner.Text()
			if text == "" {
				continue
			}
			data := text[6:]
			if data == "[DONE]" {
				return
			}
			var item StreamData
			if err = json.Unmarshal([]byte(data), &item); err != nil {
				zlog.CtxErrorf(ctx, "解析失败 %s %s\n", err, data)
				continue
			}
			if len(item.Choices) > 0 {
				msgChan <- item.Choices[0].Delta.Content
				reply += item.Choices[0].Delta.Content
			}
		}
	}()
	// 保存历史记录
	cacheUtils.SaveContent(ctx, strconv.FormatInt(userid, 10), history, content+reply)
	return
}

package ai_eino

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"strconv"
	"tgwp/global"
	"tgwp/log/zlog"
	"tgwp/types"
	"tgwp/utils/aiUtils"
)

func Chat(ctx context.Context, req types.AiReq) (msg string, err error) {
	// 创建模板，使用 FString 格式
	template := prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题。你的目标是生成合适的文章简介，不超过200字。"),

		// 用户消息模板
		schema.UserMessage("问题: {question}"),
	)

	// 使用模板生成消息
	messages, err := template.Format(context.Background(), map[string]any{
		"role":     "知识渊博的island的人工智能助手",
		"style":    "积极、温暖且专业",
		"question": req.Content,
	})
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: global.Config.AI[req.Type].ApiUrl,
		Model:   global.Config.AI[req.Type].Model,  // 使用的模型版本
		APIKey:  global.Config.AI[req.Type].ApiKey, // OpenAI API 密钥
	})
	if err != nil {
		zlog.CtxErrorf(ctx, "创建模型失败 %s", err)
		return
	}
	reader, err := chatModel.Generate(ctx, messages)
	if err != nil {
		zlog.CtxErrorf(ctx, "AI生成简介失败 %s", err)
		return
	}
	msg = reader.Content
	return
}
func StreamChat(ctx context.Context, req types.AiReq) (outStream *schema.StreamReader[*schema.Message], err error) {
	// 创建模板，使用 FString 格式
	template := prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题。你的目标回答用户提出的问题。"),

		// 插入需要的对话历史（新对话的话这里不填）
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("问题: {question}"),
	)

	// 使用模板生成消息
	messages, err := template.Format(context.Background(), map[string]any{
		"role":     "知识渊博的island的人工智能助手",
		"style":    "积极、温暖且专业",
		"question": req.Content,
		// 对话历史（这个例子里模拟两轮对话历史）
		"chat_history": loadHistory(ctx, req.UserID),
	})
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		ByAzure: false, // 是否使用 Azure OpenAI
		BaseURL: global.Config.AI[req.Type].ApiUrl,
		Model:   global.Config.AI[req.Type].Model,  // 使用的模型版本
		APIKey:  global.Config.AI[req.Type].ApiKey, // OpenAI API 密钥
	})
	if err != nil {
		zlog.CtxErrorf(ctx, "创建模型失败 %s", err)
		return
	}
	return chatModel.Stream(ctx, messages)
}
func loadHistory(ctx context.Context, user_id int64) (history []*schema.Message) {
	data, err := aiUtils.GetHistory(ctx, strconv.FormatInt(user_id, 10))
	if err != nil {
		zlog.CtxErrorf(ctx, "获取历史记录失败 %s", err)
		return
	}
	for i, v := range data {
		if i%2 == 0 {
			history = append(history, &schema.Message{
				Role:    schema.User,
				Content: v,
			})
		} else {
			history = append(history, &schema.Message{
				Role:    schema.Assistant,
				Content: v,
			})
		}
	}
	return history
}

package ai_eino

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"tgwp/global"
	"tgwp/log/zlog"
)

func Chat(ctx context.Context, content string) (msg string, err error) {
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
		"question": content,
	})
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: "https://api.chatanywhere.tech",
		Model:   global.Config.AI.Model,  // 使用的模型版本
		APIKey:  global.Config.AI.ApiKey, // OpenAI API 密钥
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

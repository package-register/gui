package sdk

import (
	"context"
	"fmt"

	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/model/openai"
	"trpc.group/trpc-go/trpc-agent-go/runner"
	"trpc.group/trpc-go/trpc-agent-go/session/inmemory"
)

// AIService 封装 AI 对话服务
type AIService struct {
	runner runner.Runner
	agent  agent.Agent
	ctx    context.Context
	userID string
}

// AIServiceConfig AI 服务配置
type AIServiceConfig struct {
	APIKey  string
	BaseURL string
	Model   string
	UserID  string
}

// NewAIService 创建新的 AI 服务
func NewAIService(config AIServiceConfig) *AIService {
	// 创建 OpenAI 模型
	modelInstance := openai.New(config.Model,
		openai.WithAPIKey(config.APIKey),
		openai.WithBaseURL(config.BaseURL),
	)

	// 生成配置
	genCfg := model.GenerationConfig{
		Stream:      true, // 启用流式响应
		MaxTokens:   intPtr(4000),
		Temperature: floatPtr(0.7),
	}

	// 创建 LLM Agent
	agentInstance := llmagent.New(
		"ai-assistant",
		llmagent.WithModel(modelInstance),
		llmagent.WithGenerationConfig(genCfg),
		llmagent.WithInstruction(`You are a helpful AI assistant for GUI application.
You can help users with various tasks and answer their questions.
Always respond in the same language the user writes in.`),
		llmagent.WithDescription("GUI Application AI Assistant"),
	)

	// 创建 Session Service
	sessionSvc := inmemory.NewSessionService()

	// 创建 Runner
	runnerInstance := runner.NewRunner(
		"gui-ai-app",
		agentInstance,
		runner.WithSessionService(sessionSvc),
	)

	userID := config.UserID
	if userID == "" {
		userID = "default-user"
	}

	return &AIService{
		runner: runnerInstance,
		agent:  agentInstance,
		ctx:    context.Background(),
		userID: userID,
	}
}

// Chat 发送消息并获取回复（同步）
func (a *AIService) Chat(message string) (string, error) {
	sessionID := generateSessionID()
	userMessage := model.NewUserMessage(message)

	// 运行 Runner
	eventCh, err := a.runner.Run(a.ctx, a.userID, sessionID, userMessage)
	if err != nil {
		return "", fmt.Errorf("AI 调用失败: %w", err)
	}

	// 收集事件，获取最终回复
	var finalContent string
	for event := range eventCh {
		if event.Error != nil {
			// 关闭 runner
			a.runner.Close()
			return "", fmt.Errorf("AI 返回错误: %s", event.Error.Message)
		}

		// 检查是否是完成事件
		if event.IsRunnerCompletion() {
			// 从响应中提取内容
			if event.Response != nil && len(event.Response.Choices) > 0 {
				finalContent = event.Response.Choices[0].Message.Content
			}
			break
		}
	}

	// 关闭 runner
	a.runner.Close()

	return finalContent, nil
}

// ChatStream 发送消息并使用流式回调接收回复
func (a *AIService) ChatStream(message string, callback func(chunk string)) error {
	sessionID := generateSessionID()
	userMessage := model.NewUserMessage(message)

	// 运行 Runner
	eventCh, err := a.runner.Run(a.ctx, a.userID, sessionID, userMessage)
	if err != nil {
		return fmt.Errorf("AI 调用失败: %w", err)
	}

	// 处理流式事件
	for event := range eventCh {
		if event.Error != nil {
			a.runner.Close()
			return fmt.Errorf("AI 返回错误: %s", event.Error.Message)
		}

		// 检查是否有新的内容块
		if event.Response != nil && len(event.Response.Choices) > 0 {
			// 从 Delta 获取流式内容
			chunk := event.Response.Choices[0].Delta.Content
			if chunk != "" && callback != nil {
				callback(chunk)
			}
		}

		// 检查是否完成
		if event.Done {
			break
		}
	}

	// 关闭 runner
	a.runner.Close()

	return nil
}

// Close 关闭 AI 服务
func (a *AIService) Close() error {
	return a.runner.Close()
}

// 辅助函数
func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}

// generateSessionID 生成会话ID（简化版）
func generateSessionID() string {
	return "session-" + fmt.Sprintf("%d", len("default-user"))
}

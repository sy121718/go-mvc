package task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-mvc/pkg/queue"

	"github.com/hibiken/asynq"
)

// ========== 任务类型 ==========

const TypeEmailSend = "email:send"

// ========== Payload 定义 ==========

type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// ========== 处理器 ==========

func HandleEmailSend(ctx context.Context, t *asynq.Task) error {
	var payload EmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("解析邮件载荷失败: %v", err)
	}

	// 这里写发送邮件的逻辑
	fmt.Printf("发送邮件: To=%s, Subject=%s\n", payload.To, payload.Subject)

	return nil
}

// ========== 业务调用方法 ==========

func SendEmail(to, subject, body string) error {
	return queue.Enqueue(TypeEmailSend, EmailPayload{
		To:      to,
		Subject: subject,
		Body:    body,
	})
}

func SendEmailDelay(to, subject, body string, delay time.Duration) error {
	return queue.EnqueueIn(TypeEmailSend, delay, EmailPayload{
		To:      to,
		Subject: subject,
		Body:    body,
	})
}

// ========== 自动注册 ==========

func init() {
	queue.Register(func() {
		queue.RegisterHandler(TypeEmailSend, HandleEmailSend)
	})
}
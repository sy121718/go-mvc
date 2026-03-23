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

const TypeOrderCancel = "order:cancel"

// ========== Payload 定义 ==========

type OrderPayload struct {
	OrderID int64 `json:"order_id"`
}

// ========== 处理器 ==========

func HandleOrderCancel(ctx context.Context, t *asynq.Task) error {
	var payload OrderPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("解析订单载荷失败: %v", err)
	}

	// 这里写订单取消的逻辑
	fmt.Printf("取消订单: OrderID=%d\n", payload.OrderID)

	return nil
}

// ========== 业务调用方法 ==========

func CancelOrder(orderID int64, delay time.Duration) error {
	return queue.EnqueueIn(TypeOrderCancel, delay, OrderPayload{
		OrderID: orderID,
	})
}

// ========== 自动注册 ==========

func init() {
	queue.Register(func() {
		queue.RegisterHandler(TypeOrderCancel, HandleOrderCancel)
	})
}
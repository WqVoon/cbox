package image

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"time"
)

const (
	defaultInterval = 30 * time.Second
	defaultTimeout  = 30 * time.Second
	defaultRetries  = 3
)

// 用来定义一个健康检查任务
type HealthCheckTaskType struct {
	// 检查周期，即多久检查一次，默认 30s
	Interval time.Duration `json:"interval"`

	// 每个 Interval 的每个 Retry 的执行都不会超过这个时间，否则记为超时，默认 30s
	Timeout time.Duration `json:"timeout"`

	// 每个 Interval 时最多进行多少次 Retry，如果其中的一次 Retry 成功，那么认为本次 Interval 成功
	// 默认为 3 次
	Retries int `json:"retries"`

	// 检查时执行的具体内容，如果执行发生错误那么认为这次 Retry 失败
	Cmd string `json:"cmd"`
}

// 检查任务是否有效，仅在有效时才应该进行 Start
func (task *HealthCheckTaskType) IsValid() bool {
	return (task.Interval >= 0 &&
		task.Timeout >= 0 &&
		task.Retries >= 0 &&
		len(task.Cmd) > 0)
}

// 开始检查，如果 Interval 中的所有 Retry 均失败，那么使用最后一次 Retry 失败的原因调用 onFailed
func (task *HealthCheckTaskType) Start(onFailed func(error, []byte)) {
	if task.Interval == 0 {
		task.Interval = defaultInterval
	}

	if task.Timeout == 0 {
		task.Timeout = defaultTimeout
	}

	if task.Retries == 0 {
		task.Retries = defaultRetries
	}

	log.Println("health check cmd:", task.Cmd)
	emptyReader := bytes.NewReader(nil)

	for range time.NewTicker(task.Interval).C {
		var content []byte
		var err error

		for i := 0; i < task.Retries; i++ {
			ctx, cancelFunc := context.WithTimeout(context.TODO(), task.Timeout)
			cmd := exec.CommandContext(ctx, "sh", "-c", task.Cmd)
			cmd.Stdin = emptyReader // 避免使用 os.DevNull，因为目前还没挂载 dev :-P
			content, err = cmd.CombinedOutput()
			cancelFunc()

			if err == nil {
				break
			}
		}

		if err != nil && onFailed != nil {
			onFailed(err, content)
		}
	}
}
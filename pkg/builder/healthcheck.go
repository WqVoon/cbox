package builder

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wqvoon/cbox/pkg/image"
	"github.com/wqvoon/cbox/pkg/log"
)

const (
	interval = 1
	timeout  = 2
	retries  = 3
)

var validOptions = map[string]int{
	"--interval": interval,
	"--timeout":  timeout,
	"--retries":  retries,
}

func buildHealthCheck(options, args []string) (*image.HealthCheckTaskType, error) {
	task := &image.HealthCheckTaskType{Cmd: args}

	for _, op := range options {
		splitedOp := strings.Split(op, "=")
		if len(splitedOp) != 2 {
			return nil, fmt.Errorf("invalid option %q, should be `--key=val`", op)
		}

		key := splitedOp[0]
		val, err := strconv.Atoi(splitedOp[1])
		if err != nil {
			return nil, fmt.Errorf("invalid option %q, value should be integer", op)
		}

		if opType, isValid := validOptions[key]; isValid {
			switch opType {
			case interval:
				task.Interval = time.Duration(val)
			case timeout:
				task.Timeout = time.Duration(val)
			case retries:
				task.Retries = val
			}
		} else {
			log.Errorf("unsupported option %q, should be `interval`, `timeout` or `retries`", key)
		}
	}

	return task, nil
}

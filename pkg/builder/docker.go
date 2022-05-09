package builder

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/utils"
)

// 用来按空白符号切分一行，忽略引号内的空白符
var reg = regexp.MustCompile(`[^\s"']+|"([^"\\]|\\.)*"|'([^'\\]|\\.)*'`)

// 解析一个 Dockerfile 文件，如果解析成功，那么生成一个 BuildTask 对象
func ParseDockerfile(filename string) *BuildTask {
	log.Println("start to parse dockerfile")

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bt := &BuildTask{}
	s := bufio.NewScanner(file)

	for lineNum := 1; s.Scan(); lineNum++ {
		tokens := reg.FindAllString(s.Text(), -1)

		// 跳过空行和注释
		if len(tokens) == 0 || tokens[0][0] == '#' {
			continue
		}

		tmpLineNum := lineNum
		for {
			lastToken := tokens[len(tokens)-1]
			lastIdx := len(lastToken) - 1

			if lastToken[lastIdx] != '\\' || !s.Scan() {
				break
			}

			tmpLineNum++
			tokens[len(tokens)-1] = lastToken[:lastIdx]
			tokens = append(tokens, reg.FindAllString(s.Text(), -1)...)
		}

		// 先传递再赋值，因为要保证在 handleOneline 的视角下，被 \ 符号连起来的行共享第一行的行号
		handleOneline(bt, lineNum, tokens)
		lineNum = tmpLineNum
	}

	log.Println("parse dockerfile done")
	return bt
}

// 能执行到这里说明 len(tokens) 一定不为 0
func handleOneline(bt *BuildTask, lineNum int, tokens []string) {
	cmd, args := tokens[0], tokens[1:]

	switch cmd {
	case "FROM":
		if len(args) != 1 {
			log.Errorf("line %d error: format error, should be `FROM <name>:<tag>`\n", lineNum)
		}

		if bt.BaseImageStr != "" {
			log.Errorf("line %d error: duplicated base image", lineNum)
		}

		bt.BaseImageStr = args[0]

	case "RUN":
		if len(args) == 0 {
			log.Errorf("line %d error: format error, should be `RUN <cmd> [...<cmd>]`\n", lineNum)
		}

		bt.Cmds = append(bt.Cmds, args)

	case "ENV":
		if len(args) != 2 {
			log.Errorf("line %d error: format error, should be `ENV <key> <val>`\n", lineNum)
		}

		if bt.Env == nil {
			bt.Env = make(map[string]string)
		}

		key, val := args[0], args[1]

		bt.Env[key] = strings.Trim(val, "\"'")

	case "COPY":
		if len(args) != 2 {
			log.Errorf("line %d error: format error, should be `COPY src dst`\n", lineNum)
		}

		bt.Copy = append(bt.Copy, CopyType{Src: args[0], Dst: args[1]})

	case "ENTRYPOINT":
		if len(args) == 0 {
			log.Errorf("line %d error: format error, should be `ENTRYPOINT <cmd> [...<cmd>]`\n", lineNum)
		}

		bt.Entrypint = args

	case "NAME":
		if len(args) != 1 {
			log.Errorf("line %d error: format error, should be `Name <name>:<tag>`\n", lineNum)
		}

		if bt.ImageNameStr != "" {
			log.Errorf("line %d error: duplicated image name", lineNum)
		}

		bt.ImageNameStr = args[0]

	case "HEALTHCHECK":
		if bt.HealthCheckTask != nil {
			log.Errorf("line %d error: duplicated health check", lineNum)
		}

		optionsBorder := utils.FindIdxInStringSlice(args, "CMD")
		if optionsBorder == -1 {
			log.Errorf("line %d error: format error, should be `HEALTHCHECK [options] CMD <cmd> [...<cmd>]`\n", lineNum)
		}

		options, args := args[:optionsBorder], args[optionsBorder+1:]
		task, err := buildHealthCheck(options, args)
		if err != nil {
			log.Errorf("line %d error: %v\n", lineNum, err)
		}

		bt.HealthCheckTask = task

	default:
		log.Errorf("line %d error: unsupported cmd\n", lineNum)
	}
}

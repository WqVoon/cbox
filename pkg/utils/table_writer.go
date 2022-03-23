package utils

import (
	"fmt"
	"strings"

	"github.com/wqvoon/cbox/pkg/log"
)

// TODO: 一个勉强能看的表格，后面可以美化下

type TableWriter struct {
	header    string
	width     int
	fixSpaces func(string) string
}

// header 是标题，width 是每列的宽度，如果字段的长度大于列宽-1，那么会被裁减，且后三位变成省略号
func NewTableWriter(header []string, width int) *TableWriter {
	if width < 5 {
		log.Errorln("width can not less then 5")
	}

	spaces := strings.Repeat(" ", width)
	fixSpaces := func(name string) string {
		length := len(name)
		if length > width-1 {
			return name[:width-4] + "... "
		}
		return name + spaces[length:]
	}

	for idx, field := range header {
		if len(field) > width-1 {
			log.Errorln("length of header field can not greater than width-1")
		}

		field = strings.ReplaceAll(strings.ToUpper(field), " ", "_")
		header[idx] = fixSpaces(field)
	}

	return &TableWriter{
		header:    strings.Join(header, ""),
		width:     width + 1,
		fixSpaces: fixSpaces,
	}
}

func (tw *TableWriter) PrintlnHeader() {
	fmt.Println(tw.header)
}

func (tw *TableWriter) PrintlnData(data ...string) {
	for idx, field := range data {
		data[idx] = tw.fixSpaces(field)
	}
	fmt.Println(strings.Join(data, ""))
}

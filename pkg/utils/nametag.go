package utils

import (
	"fmt"
	"strings"

	"github.com/wqvoon/cbox/pkg/log"
)

type NameTag struct {
	Name, Tag string
}

func GetNameTag(nt string) *NameTag {
	if len(strings.Trim(nt, " ")) == 0 {
		log.Errorln("error format of name tag, should be `name:tag`")
	}

	splitedNt := strings.Split(nt, ":")
	if len(splitedNt) > 2 {
		log.Errorln("error format of name tag, should be `name:tag`")
	}

	name := splitedNt[0]
	tag := "latest"
	if len(splitedNt) == 2 {
		tag = splitedNt[1]
	}

	return &NameTag{name, tag}
}

func (nt *NameTag) String() string {
	return fmt.Sprintf("%s:%s", nt.Name, nt.Tag)
}

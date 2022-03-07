package utils

import (
	"fmt"
	"log"
	"strings"
)

type NameTag struct {
	Name, Tag string
}

func GetNameTag(nt string) *NameTag {
	splitedNt := strings.Split(nt, ":")
	if len(splitedNt) > 2 {
		log.Fatalln("error format of name tag, should be `name:tag`")
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

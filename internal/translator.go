package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

type translator struct {
	targetLang string
}

func newTranslator() translator {
	return translator{targetLang: "en"}
}

func (t *translator) translate(text string) (string, error) {
	workdir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	transPath := path.Join(workdir, "/translator/trans")
	command := exec.Command(transPath, "-no-warn", "-b", fmt.Sprintf(":%s", t.targetLang), text)
	out, err := command.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

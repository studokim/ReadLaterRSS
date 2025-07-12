package internal

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

type translator struct {
	db *sql.DB
}

func newTranslator() (translator, error) {
	db, err := Sqlite{}.db()
	if err != nil {
		return translator{}, err
	}
	return translator{db: db}, nil
}

func (t translator) shouldTranslate(feedTitle string) bool {
	_, err := t.getTranslateSettings(feedTitle)
	return err == nil
}

func (t translator) getTranslateSettings(feedTitle string) (translate, error) {
	row := t.db.QueryRow("select langFrom, langTo from translate where feedTitle=?", feedTitle)
	var langFrom string
	var langTo string
	err := row.Scan(&langFrom, &langTo)
	if err != nil {
		return translate{}, err
	}
	return translate{FeedTitle: feedTitle, LangFrom: langFrom, LangTo: langTo}, nil
}

func (t translator) translate(feedTitle string, text string) (string, error) {
	workdir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	transPath := path.Join(workdir, "/translator/trans")
	ts, err := t.getTranslateSettings(feedTitle)
	if err != nil {
		return "", err
	}
	command := exec.Command(transPath, "-no-warn", "-b", fmt.Sprintf("%s:%s", ts.LangFrom, ts.LangTo), text)
	out, err := command.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

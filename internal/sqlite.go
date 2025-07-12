package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

const dbFile string = "ReadLaterRSS.db"

type Sqlite struct {
	_db *sql.DB
}

func (s Sqlite) db() (*sql.DB, error) {
	if s._db == nil {
		db, err := sql.Open("sqlite", dbFile)
		if err != nil {
			return nil, err
		}
		s._db = db
		err = s.migrate()
		if err != nil {
			return nil, err
		}
	}
	return s._db, nil
}

func (s Sqlite) migrate() error {
	var version int
	db, err := s.db()
	if err != nil {
		return err
	}
	err = db.QueryRow("select version from meta").Scan(&version)
	if err != nil && err.Error() != "SQL logic error: no such table: meta (1)" {
		return err
	}

	if version < 1 {
		log.Println("Migrating: creating new db from scratch")
		// https://stackoverflow.com/a/65743498
		_, err = db.Exec(`PRAGMA writable_schema = 1;
				          DELETE FROM sqlite_master;
				          PRAGMA writable_schema = 0;
				          VACUUM;
				          PRAGMA integrity_check;`)
		if err != nil {
			return err
		}

		_, err = db.Exec("create table feed(title unique not null, description not null, author not null, email not null, feedType not null check(feedType in ('url','text')) )")
		if err != nil {
			return err
		}
		fdef := feed{Title: "Default", Description: "Automatically created feed", Author: "ReadLaterRSS", Email: "-", FeedType: urlType}
		fdeu := feed{Title: "Deutsch", Description: "Auto-translated Deutsch texts", Author: "ReadLaterRSS", Email: "-", FeedType: textType}
		_, err = db.Exec("insert into feed(title, description, author, email, feedType) values (?,?,?,?,?), (?,?,?,?,?)",
			fdef.Title, fdef.Description, fdef.Author, fdef.Email, fdef.FeedType,
			fdeu.Title, fdeu.Description, fdeu.Author, fdeu.Email, fdeu.FeedType)
		if err != nil {
			return err
		}
		_, err = db.Exec("create table item(feedTitle not null, id unique not null, title not null, created timestamp not null, url not null, text not null, foreign key(feedTitle) references feed(title))")
		if err != nil {
			return err
		}

		_, err = db.Exec("create table translate(feedTitle not null, langFrom not null check(length(langFrom) = 2), langTo not null check(length(langTo) = 2), foreign key(feedTitle) references feed(title))")
		if err != nil {
			return err
		}
		_, err = db.Exec(`CREATE TRIGGER prevent_translate_when_not_textType
			BEFORE INSERT ON translate FOR EACH ROW WHEN NEW.feedTitle IS NOT NULL
			BEGIN
    			SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM feed WHERE title = NEW.feedTitle AND feedType = 'text') THEN
        			RAISE(ABORT, 'Cannot translate NEW.feedTitle. Can only translate feeds of type text.')
    			END;
			END;`)
		if err != nil {
			return err
		}
		tdeu := translate{FeedTitle: fdeu.Title, LangFrom: "de", LangTo: "en"}
		_, err = db.Exec("insert into translate(feedTitle, langFrom, langTo) values (?,?,?)", tdeu.FeedTitle, tdeu.LangFrom, tdeu.LangTo)
		if err != nil {
			return err
		}

		_, err = db.Exec("create table meta(version int); insert into meta(version) values(1)")
		return err

	} else if version == 1 {
		log.Println("Not migrating, as already on version=1")
		return nil
	} else {
		return errors.New(fmt.Sprint("Expected version <= 1, got", version))
	}
}

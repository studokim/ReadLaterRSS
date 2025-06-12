package internal

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/glebarez/go-sqlite"
)

type Sqlite struct {
	_db *sql.DB
}

func (s *Sqlite) db() *sql.DB {
	if s._db == nil {
		db, err := sql.Open("sqlite", "history.db")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		s._db = db
		s.migrate()
	}
	return s._db
}

func (s *Sqlite) migrate() {
	fmt.Println("Caution: migrating!")
	_, err := s.db().Exec("create table if not exists record(feed not null, title, url, text, [when] timestamp not null)")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (s *Sqlite) Test() {
	version := s.db().QueryRow("select sqlite_version();")
	if version == nil {
		fmt.Println("version is nil")
		os.Exit(1)
	}
	if version.Err() != nil {
		fmt.Println(version.Err())
		os.Exit(1)
	}

	var column1 string
	err := version.Scan(&column1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("begin", column1, "end")
}

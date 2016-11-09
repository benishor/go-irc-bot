package plugins

import (
	"time"
	"strings"
	"github.com/kardianos/osext"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"path"
	"log"
	"github.com/benishor/go-irc-bot/irc"
)

type lastSeenInfo struct {
	User        string
	Channel     string
	LeavingTime time.Time
	Message     string
}

type lastSeenTool struct {
	database *sql.DB
}

func NewLastSeenToolSql() (*lastSeenTool) {
	folder, _ := osext.ExecutableFolder()

	db, err := sql.Open("sqlite3", path.Join(folder, "botdb.sqlite"))
	checkError(err)

	createTableSql := `
	CREATE TABLE IF NOT EXISTS seen(
		Nickname TEXT NOT NULL PRIMARY KEY,
		User TEXT,
		Channel TEXT,
		LeavingTime BIGINT,
		Message TEXT
	);`

	_, err = db.Exec(createTableSql)
	checkError(err)

	return &lastSeenTool{database: db}
}

func checkError(err error) {
	if nil != err {
		panic(err)
	}
}

func (t*lastSeenTool) set(user, channel, message string) {
	nickname := strings.ToLower(irc.ParseIrcUser(user).Nickname)

	sqlUpdateSeenInfo := `INSERT OR REPLACE INTO seen(Nickname, User, Channel, LeavingTime, Message) VALUES(?, ?, ?, ?, ?)`
	stmt, err := t.database.Prepare(sqlUpdateSeenInfo)
	checkError(err)

	stmt.Exec(nickname, user, channel, time.Now().Unix(), message);
}

func (t*lastSeenTool) get(queriedNickname string) (*lastSeenInfo) {
	var nickname, user, channel, message string
	var leavingTime int64

	query := `SELECT Nickname, User, Channel, LeavingTime, Message FROM seen WHERE Nickname LIKE ?`
	err := t.database.QueryRow(query, queriedNickname).Scan(&nickname, &user, &channel, &leavingTime, &message)

	if err == nil {
		return &lastSeenInfo{
			User: user,
			Channel: channel,
			LeavingTime: time.Unix(leavingTime, 0),
			Message: message}
	} else if err != sql.ErrNoRows {
		log.Printf("Failed to retrieve last seen info. Reason: %s", err)
	}

	return nil
}


package notes

import (
	"goirc/model"
)

type Note struct {
	Id        int64
	CreatedAt string `db:"created_at"`
	Text      string
	Nick      string
	Kind      string
}

func Create(target string, nick string, kind string, text string) error {
	result, err := model.DB.Exec(`insert into notes(nick, text, kind) values(?, ?, ?) returning id`, nick, text, kind)
	if err != nil {
		return err
	} else {
		noteId, err := result.LastInsertId()
		if err != nil {
			return err
		}
		err = MarkAsSeen(noteId, target)
		if err != nil {
			return err
		}
	}
	return nil
}

func MarkAsSeen(noteId int64, target string) error {
	//db.Select(`select * from channel_nicks where channel = ?`, target)
	channelNicks, err := model.JoinedNicks(target)
	if err != nil {
		return err
	}
	// for each channelNick insert a seen_by record
	for _, nick := range channelNicks {
		_, err := model.DB.Exec(`insert into seen_by(note_id, nick) values(?, ?)`, noteId, nick.Nick)
		if err != nil {
			return err
		}
	}
	return nil
}

package web

import (
	"bytes"
	_ "embed"
	"fmt"
	"goirc/model"
	"goirc/util"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type NickWithNoteCount struct {
	Nick  string
	Count int
}

//go:embed "templates/index.gohtml"
var indexTemplate string

//go:embed "templates/rss.gohtml"
var rssTemplate string

func Serve(db *sqlx.DB) {
	r := gin.Default()
	//r.LoadHTMLGlob("templates/*")

	r.GET("/snapshot.db", func(c *gin.Context) {
		os.Remove("/tmp/snapshot.db")
		if _, err := db.Exec(`vacuum into '/tmp/snapshot.db'`); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
			return
		}
		c.File("/tmp/snapshot.db")
	})
	r.GET("/", func(c *gin.Context) {
		nick := c.Query("nick")

		notes, err := getNotes(db, nick)
		if err != nil {
			log.Fatal(err)
		}

		nicks, err := getNicks(db)
		if err != nil {
			log.Fatal(err)
		}

		tmpl, err := template.New("name").Parse(indexTemplate)
		if err != nil {
			log.Fatal("error parsing template")
		}

		out := new(bytes.Buffer)
		err = tmpl.Execute(out, gin.H{
			"nicks": nicks,
			"notes": notes,
		})
		if err != nil {
			log.Fatal("error executing template on data")
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", out.Bytes())
	})

	r.GET("/rss.xml", func(c *gin.Context) {
		nick := c.Query("nick")

		notes, err := getNotes(db, nick)
		if err != nil {
			log.Fatal(err)
		}

		tmpl, err := template.New("name").Parse(rssTemplate)
		if err != nil {
			log.Fatal("error parsing template")
		}

		fnotes, err := formatNotesDates(notes)
		if err != nil {
			log.Fatalf("error formatting notes: %v", err)
		}

		out := new(bytes.Buffer)
		err = tmpl.Execute(out, gin.H{
			"notes": fnotes,
		})
		if err != nil {
			log.Fatal("error executing template on data")
		}

		c.Data(http.StatusOK, "text/xml; charset=utf-8", out.Bytes())
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func getNotes(db *sqlx.DB, nick string) ([]model.Note, error) {
	notes := []model.Note{}
	var err error
	if nick == "" {
		err = db.Select(&notes, `select created_at, text, nick, kind from notes order by created_at desc limit 1000`)
	} else {
		err = db.Select(&notes, `select created_at, text, nick, kind from notes where nick = ? order by created_at desc limit 1000`, nick)
	}
	return notes, err
}

func getNicks(db *sqlx.DB) ([]NickWithNoteCount, error) {
	nicks := []NickWithNoteCount{}
	err := db.Select(&nicks, `select nick, count(nick) as count from notes group by nick`)
	return nicks, err
}

func formatNotesDates(notes []model.Note) ([]model.Note, error) {
	result := []model.Note{}
	for _, n := range notes {
		newNote := n

		createdAt, err := util.ParseTime(n.CreatedAt)
		if err != nil {
			return nil, err
		}

		newNote.CreatedAt = createdAt.Format("Mon, 02 Jan 2006 15:04:05 -0700")
		result = append(result, newNote)
	}
	return result, nil
}
/*
Exercise 4.12: The popular web comic xkcd has a JSON interface. For example, a request to
https://xkcd.com/571/info.0.json produces a detailed des cription of comic 571, one of
many favorites. Download each URL (once!) and build an offline index. Write a tool xkcd
that, using this index, prints the URL and transcript of each comic that matches a search term
provided on the command line.
*/
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	BaseURL      = "https://xkcd.com/"
	JsonFileName = "info.0.json"
	SQLiteFile   = "/home/harald/go/data/xkcd.sqlite"
)

var (
	Nums2Ignore = map[int]bool{404: true}
)

type Comic struct {
	Num        int
	Year       int `json:"year,string"`
	Month      int `json:"month,string"`
	Day        int `json:"day,string"`
	Link       string
	News       string
	SafeTitle  string `json:"safe_title"`
	Transcript string
	Alt        string
	Img        string
	Title      string
}

type ComicMap map[int]*Comic

func FetchComic(number int) (*Comic, error) {
	var Comic Comic
	if number < 1 {
		return nil, fmt.Errorf("Invalid Comic Number")
	}

	res, err := http.Get(BaseURL + strconv.Itoa(number) + "/" + JsonFileName)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("Item not found")
		} else {
			return nil, fmt.Errorf("An error occured: %s", res.Status)
		}
	}
	if err = json.NewDecoder(res.Body).Decode(&Comic); err != nil {
		return nil, err
	}
	return &Comic, nil
}

func FetchMissingComics(cm *ComicMap) (*ComicMap, error) {
	var latest int
	res, err := http.Get(BaseURL + JsonFileName)
	if err != nil {
		return nil, err
	}
	{
		var ComicLast Comic
		err := json.NewDecoder(res.Body).Decode(&ComicLast)
		res.Body.Close()
		if err != nil {
			return nil, err
		} else {
			latest = ComicLast.Num
		}
	}
	for i := 0; i < latest; i++ {
		comicNr := i + 1
		if Nums2Ignore[comicNr] {
			continue
		}
		if _, ok := (*cm)[comicNr]; ok {
			continue
		}
		if cn, err := FetchComic(comicNr); err == nil {
			(*cm)[i] = cn
		}
	}
	return cm, nil
}

func OpenAndInitSQLiteDB(file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	sql :=
		`CREATE TABLE IF NOT EXISTS comics(
		 Num  INTEGER PRIMARY KEY,
		 Year INTEGER,
		 Month INTEGER,
		 Day INTEGER,
		 Link TEXT,
		 News TEXT,
		 Title TEXT,
		 SafeTitle TEXT,
		 Transcript TEXT,
		 Alt TEXT,
		 Img TEXT
		 );
   CREATE VIRTUAL TABLE IF NOT EXISTS comics_idx
	  USING FTS5(
			Title, 
			Transcript, 
			Alt,
			content='comics',
			content_rowid='Num'
		);
		create trigger IF NOT EXISTS comics_ai AFTER INSERT ON comics BEGIN
		  INSERT INTO comics_idx(rowid,Title,Transcript,Alt) 
		    VALUES(new.Num,new.Title,new.Transcript,new.Alt);
		END;
		create trigger IF NOT EXISTS comics_ad AFTER DELETE ON comics BEGIN
		  INSERT INTO comics_idx(comics_idx,rowid,Title,Transcript,Alt) 
		    VALUES('delete',old.Num,old.Title,old.Transcript,old.Alt);
		END;
		create trigger IF NOT EXISTS comics_au AFTER UPDATE ON comics BEGIN
		  INSERT INTO comics_idx(comics_idx,rowid,Title,Transcript,Alt) 
		    VALUES('delete',old.Num,old.Title,old.Transcript,old.Alt);
		  INSERT INTO comics_idx(rowid,Title,Transcript,alt) 
		    VALUES(new.Num,new.Title,new.Transcript,new.Alt);
		END;
		`
	_, err = db.Exec(sql)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ReadComicFromSQLiteDB(db *sql.DB, comicNr int) (*Comic, error) {
	if comicNr < 1 {
		return nil, fmt.Errorf("Invalid Comic Number")
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	SQL := `SELECT Num,Year,Month,Day,Link,News,Title,SafeTitle,Transcript,Alt,Img
						FROM comics
					WHERE Num = ?`
	row := db.QueryRow(SQL, comicNr)
	var c Comic
	if err := row.Scan(&c.Num, &c.Year, &c.Month, &c.Day, &c.Link, &c.News, &c.Title, &c.SafeTitle, &c.Transcript, &c.Alt, &c.Img); err != nil {
		return nil, err
	}
	return &c, nil
}
func ReadComicsFromSQLiteDB(db *sql.DB, query *string) (*ComicMap, error) {
	if err := db.Ping(); err != nil {
		return nil, err
	}
	var count int
	comics := make(ComicMap)
	SQL := "SELECT count(*) FROM comics"
	row := db.QueryRow(SQL)
	err := row.Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return &comics, nil
		}
		return nil, err
	}
	if count == 0 {
		return &comics, nil
	}
	var rows *sql.Rows
	if query == nil || *query == "" {
		// Fetch all
		SQL = "SELECT Num,Year,Month,Day,Link,News,Title,SafeTitle,Transcript,Alt,Img FROM comics"
		rows, err = db.Query(SQL)
	} else {
		SQL = `
			SELECT c.Num,c.Year,c.Month,c.Day,c.Link,c.News,c.Title,c.SafeTitle,c.Transcript,c.Alt,c.Img
				FROM comics c JOIN comics_idx ci ON ci.rowid = c.Num 
			WHERE comics_idx MATCH ?
			`
		rows, err = db.Query(SQL, *query)
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var c Comic
		if err := rows.Scan(&c.Num, &c.Year, &c.Month, &c.Day, &c.Link, &c.News, &c.Title, &c.SafeTitle, &c.Transcript, &c.Alt, &c.Img); err != nil {
			rows.Close()
			return nil, err
		}
		comics[c.Num] = &c
	}
	rows.Close()
	return &comics, nil
}

func UpdateSQLiteDB(db *sql.DB, cm *ComicMap) error {
	if err := db.Ping(); err != nil {
		return err
	}
	SQL := "SELECT Num FROM comics"
	ids := make(map[int]bool)
	rows, err := db.Query(SQL)
	if err != nil {
		return err
	}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		ids[id] = true
	}
	rows.Close()

	SQL = `INSERT INTO comics(
		       Num,Year,Month,Day,Link,News,Title,SafeTitle,Transcript,Alt,Img
				 ) VALUES(?,?,?,?,?,?,?,?,?,?,?);`
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(SQL)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, c := range *cm {
		if c == nil || (ids[c.Num] || Nums2Ignore[c.Num]) {
			continue
		}
		_, err := stmt.Exec(
			c.Num, c.Year, c.Month, c.Day, c.Link, c.News,
			c.Title, c.SafeTitle, c.Transcript, c.Alt, c.Img,
		)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Fprintf(os.Stderr, "Could not roll back: %v\n", rbErr)
			}
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func DisplayComic(c *Comic) error {
	if !(strings.Index(c.Img, "http://") == 0 ||
		strings.Index(c.Img, "https://") == 0) {
		return fmt.Errorf("Invalid URL: %s", c.Img)
	}
	var f *os.File
	var err error
	for _, d := range []string{"", "/tmp"} {
		if f, err = ioutil.TempFile(d, ".dc"); err == nil {
			break
		}
	}
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	res, err := http.Get(c.Img)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Item not found")
		} else {
			return fmt.Errorf("An error occured: %s", res.Status)
		}
	}
	var r io.LimitedReader
	r.R = res.Body
	r.N = 50 * 1024 * 1024
	for buf := make([]byte, 64*1024); ; {
		n, err := r.Read(buf)
		if !(err == nil || err == io.EOF) {
			return err
		}
		if n <= 0 {
			break
		}
		f.Write(buf[:n])
	}
	f.Close()
	if r.N <= 0 {
		return fmt.Errorf("Image too big")
	}
	app, err := gtk.ApplicationNew("cc11.de.golang412", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		return nil
	}
	app.Connect("activate", func() {
		display, _ := gdk.DisplayGetDefault()
		monitor, _ := display.GetPrimaryMonitor()
		geo := monitor.GetGeometry()
		mh, mw := float64(geo.GetHeight()), float64(geo.GetWidth())
		pixbuf, _ := gdk.PixbufNewFromFile(f.Name())
		h, w := float64(pixbuf.GetHeight()), float64(pixbuf.GetWidth())
		var sf float64
		switch {
		case h > mh && w <= mw:
			sf = mh / h
		case w > mw && h <= mh:
			sf = mw / w
		case w > mw && h > mh:
			rescale_w := mw / w
			rescale_h := mh / h
			if rescale_w < rescale_h {
				sf = rescale_w
			} else {
				sf = rescale_h
			}
		default:
			sf = 1
		}
		if sf < 1 {
			sf *= 0.95
			pixbuf, _ = pixbuf.ScaleSimple(int(w*sf), int(h*sf), gdk.INTERP_BILINEAR)
		}
		img, _ := gtk.ImageNewFromPixbuf(pixbuf)
		label, _ := gtk.LabelNew("")
		label.SetMarkup("<b>" + c.Title + "</b>")
		box, _ := gtk.FlowBoxNew()
		eventbox, _ := gtk.EventBoxNew()
		eventbox.Add(img)
		box.Add(eventbox)
		box.Add(label)
		appWindow, _ := gtk.ApplicationWindowNew(app)
		appWindow.SetTitle(c.Title)
		appWindow.Add(box)
		appWindow.SetPosition(gtk.WIN_POS_CENTER)
		appWindow.Resize(100, 100)
		appWindow.ShowAll()
		box.Connect("button_press_event", func() {
			appWindow.Destroy()
		})
	})
	app.Run(make([]string, 0))
	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr,
		"\nUsage: %[1]s -s <search term>\n"+
			"       %[1]s -d <ComicNr> (Display comic number <ComicNr>)\n"+
			"       %[1]s -u (Updates local database)\n\n",
		os.Args[0])
}

/* Unused */
/*
func OpenOrCreateDB(file string) (*ComicList, error) {
	var comics ComicList
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err = json.NewDecoder(f).Decode(&comics); !(err == nil || err == io.EOF) {
		return nil, err
	}
	return &comics, nil
}
*/
/* Unused */
/*
func WriteDB(file string, comics *ComicList) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = json.NewEncoder(f).Encode(&comics); !(err == nil || err == io.EOF) {
		return err
	}
	return nil
}
*/

func main() {
	var q string
	var comicNr int
	var err error
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	} else if os.Args[1] == "-u" {
		if len(os.Args) != 2 {
			usage()
			os.Exit(1)
		}
	} else if os.Args[1] == "-s" {
		if len(os.Args) < 3 {
			usage()
			os.Exit(1)
		}
		q = strings.Join(os.Args[2:], " ")
	} else if os.Args[1] == "-d" {
		if len(os.Args) != 3 {
			usage()
			os.Exit(1)
		}
		if comicNr, err = strconv.Atoi(os.Args[2]); err != nil || comicNr < 1 {
			usage()
			os.Exit(1)
		}
	} else {
		usage()
		os.Exit(1)
	}
	db, err := OpenAndInitSQLiteDB(SQLiteFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: Error: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	defer db.Close()
	if os.Args[1] == "-s" {
		comics, err := ReadComicsFromSQLiteDB(db, &q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Oops: Error: %v\n", os.Args[0], err)
			os.Exit(1)
		}
		comicnums := make([]int, 0, len(*comics))
		for k, _ := range *comics {
			comicnums = append(comicnums, k)
		}
		sort.Ints(comicnums)
		fmt.Fprintf(os.Stderr, "DEBUG: %+v\n", comicnums)
		for _, k := range comicnums {
			fmt.Fprintf(os.Stdout, "Number %d: %s\n%s\n\n", k, (*comics)[k].Title, strings.Repeat("=", len("Number:  "+strconv.Itoa(k)+(*comics)[k].Title)))
			fmt.Fprintf(os.Stdout, "Transcript:\n-----------\n%s\n\nWebsite: %s%d/\nImage: %s\n\n", (*comics)[k].Transcript, BaseURL, k, (*comics)[k].Img)
			DisplayComic((*comics)[k])
		}
		os.Exit(0)
	} else if os.Args[1] == "-u" {
		comics, err := ReadComicsFromSQLiteDB(db, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Oops: Error: %v\n", os.Args[0], err)
			os.Exit(1)
		}
		if _, err = FetchMissingComics(comics); err != nil {
			fmt.Fprintf(os.Stderr, "%s: Error: %v\n", os.Args[0], err)
			os.Exit(1)
		}
		if err := UpdateSQLiteDB(db, comics); err != nil {
			fmt.Fprintf(os.Stderr, "%s: Error: %v\n", os.Args[0], err)
			os.Exit(1)
		}
	} else if os.Args[1] == "-d" {
		comic, err := ReadComicFromSQLiteDB(db, comicNr)
		if err != nil {
			comic, err = FetchComic(comicNr)
		}
		if comic == nil {
			fmt.Fprintf(os.Stderr, "%s: Oops: Comic Nr. %d not found:\n", os.Args[0], comicNr)
			os.Exit(1)
		}
		if err := DisplayComic(comic); err != nil {
			fmt.Fprintf(os.Stderr, "%s: Error: %v\n", os.Args[0], err)
			os.Exit(1)
		}
	} else {
		usage()
		os.Exit(1)
	}

	os.Exit(0)
}

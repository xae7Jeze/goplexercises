/*
Exercise 4.13: The JSON-based web service of the Open Movie Database lets you search
https://omdbapi.com/ for a movie by name and download its poster image . Write a tool
poster that downloads the poster image for the movie named on the command line.
*/
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// read from file <programName>.json in current directory
var (
	BaseURL = "https://www.omdbapi.com/"
	APIKey  = ""
)

type Config struct {
	BaseURL string
	APIKey  string
}

type Movie struct {
	Title    string
	Poster   string
	Response string
	Error    string `json:",omitempty"`
}

func readConfig(cfgFile string) (*Config, error) {
	var c Config
	buf := make([]byte, 4096)
	f, err := os.Open(cfgFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if l, err := f.Read(buf); err != nil {
		return nil, err
	} else {
		buf = buf[:l]
	}
	if err := json.Unmarshal(buf, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func FetchMovie(title string) (*Movie, error) {
	var m Movie
	if title == "" {
		return nil, fmt.Errorf("Empty Title")
	}
	t := url.QueryEscape(title)

	res, err := http.Get(BaseURL + "/?apikey=" + APIKey + "&t=" + t)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("Movie not found")
		} else {
			return nil, fmt.Errorf("An error occured: %s", res.Status)
		}
	}
	if err = json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, err
	}
	if m.Response != "True" {
		if m.Error == "" {
			m.Error = "Unknown Error"
		}
		return nil, fmt.Errorf(m.Error)
	}
	if !(strings.Index(m.Poster, "http://") == 0 ||
		strings.Index(m.Poster, "https://") == 0) {
		return nil, fmt.Errorf("Movie not found")
	}

	return &m, nil
}

func DisplayPoster(m *Movie) error {
	if !(strings.Index(m.Poster, "http://") == 0 ||
		strings.Index(m.Poster, "https://") == 0) {
		return fmt.Errorf("Invalid URL: %s", m.Poster)
	}
	var f *os.File
	var err error
	for _, d := range []string{"", "/tmp"} {
		if f, err = ioutil.TempFile(d, ".mv"); err == nil {
			break
		}
	}
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	res, err := http.Get(m.Poster)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Poster not found")
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
	app, err := gtk.ApplicationNew("cc11.de.golang413", glib.APPLICATION_FLAGS_NONE)
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
		label.SetMarkup("<b>" + m.Title + "</b>")
		box, _ := gtk.FlowBoxNew()
		eventbox, _ := gtk.EventBoxNew()
		eventbox.Add(img)
		box.Add(eventbox)
		box.Add(label)
		appWindow, _ := gtk.ApplicationWindowNew(app)
		appWindow.SetTitle(m.Title)
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
		"\nUsage: %[1]s -t <Movie Title>\n",
		os.Args[0])
}

func main() {
	var title string
	var cfgFile = path.Base(os.Args[0]) + ".json"
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	} else if os.Args[1] == "-t" {
		if len(os.Args) < 3 {
			usage()
			os.Exit(1)
		}
		title = strings.Join(os.Args[2:], " ")
	} else {
		usage()
		os.Exit(1)
	}
	if cf, err := readConfig(cfgFile); err != nil {
		fmt.Fprintf(os.Stderr, "%s: Error reading config from '%s': %v\n", os.Args[0], cfgFile, err)
		os.Exit(1)
	} else {
		BaseURL = cf.BaseURL
		APIKey = cf.APIKey
	}
	if os.Args[1] == "-t" {
		movie, err := FetchMovie(title)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Oops: Error: %v\n", os.Args[0], err)
			os.Exit(1)
		}
		DisplayPoster(movie)
		os.Exit(0)
	}
	os.Exit(0)
}

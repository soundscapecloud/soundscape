package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/streamlist/streamlist/internal/archiver"
	"github.com/streamlist/streamlist/internal/youtube"

	"github.com/disintegration/imaging"
	"github.com/eduncan911/podcast"
	"github.com/julienschmidt/httprouter"
	"github.com/rylio/ytdl"
)

type response struct {
	Config   Config
	Request  *http.Request
	Params   *httprouter.Params
	HTTPHost string
	Version  string
	Backlink string
	DiskInfo *diskInfo
	Archiver *archiver.Archiver

	Error   string
	User    string
	IsAdmin bool
	Section string

	// Paging
	Page       int64
	Pages      []int64
	Limit      int64
	Total      int64
	GrandTotal int64

	// Search
	Query string

	List   *List
	Lists  []*List
	Media  *Media
	Medias []*Media

	ActiveMedias []*Media
	QueuedMedias []*Media

	Youtubes []youtube.Video
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func newResponse(r *http.Request, ps httprouter.Params) *response {
	diskInfo, err := newDiskInfo(datadir)
	if err != nil {
		panic(err)
	}
	user, _, _ := r.BasicAuth()
	isAdmin := stringInSlice(user, httpAdminUsers)
	return &response{
		Config:   config.Get(),
		Request:  r,
		Params:   &ps,
		User:     ps.ByName("user"),
		IsAdmin:  isAdmin,
		HTTPHost: httpHost,
		Version:  version,
		Backlink: backlink,
		DiskInfo: diskInfo,
		Archiver: archive,
	}
}

func logs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, line := range logtail.Lines() {
		fmt.Fprintf(w, "%s\n", line)
	}
}

func index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	redirect(w, r, "/")
}

/*func createUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Create a user
	u1 := User{Username: "admin", Password: "admin"}
	db.Create(&u1)
	fmt.Fprintln(w, "user created")
}*/

func home(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lists, err := listLists()
	if err != nil {
		_error(w, err)
		return
	}
	res := newResponse(r, ps)
	res.Section = "home"
	res.Lists = lists
	html(w, "home.html", res)
}

func configHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	switch key {
	case "volume":
		n, err := strconv.ParseFloat(value, 32)
		if err != nil {
			_error(w, err)
			return
		}
		if err := config.SetVolume(float32(n)); err != nil {
			_error(w, err)
			return
		}
	}
	toJSON(w, "OK")
}

func importHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var youtubes []youtube.Video
	if query := strings.TrimSpace(r.FormValue("q")); query != "" {
		yt, err := youtube.Search(query)
		if err != nil {
			logger.Errorf("query %q failed: %s", query, err)
		} else {
			youtubes = append(youtubes, yt...)
		}
	}

	var filtered []youtube.Video
	for _, v := range youtubes {
		// Already exists in library, so filter it out.
		if m, err := loadMedia(v.ID); err == nil {
			if m.HasAudio() || archive.InProgress(m.ID) {
				continue
			}
		}
		filtered = append(filtered, v)
	}
	youtubes = filtered

	res := newResponse(r, ps)
	res.Section = "import"
	res.Youtubes = youtubes
	html(w, "import.html", res)
}

func library(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	medias, err := ListMedias()
	if err != nil {
		_error(w, err)
		return
	}

	grandTotal := int64(len(medias))

	query := r.FormValue("q")

	// Filter
	if query != "" {
		var filtered []*Media
		for _, m := range medias {
			content := m.Title
			content += m.Description
			content += m.Author
			content += m.Source
			if !strings.Contains(strings.ToLower(content), strings.ToLower(query)) {
				continue
			}
			filtered = append(filtered, m)
		}
		medias = filtered
	}

	// pagination
	var limit int64 = 10
	page, _ := strconv.ParseInt(r.FormValue("p"), 10, 64)
	if page < 1 {
		page = 1
	}

	total := int64(len(medias))
	switch {
	case total > 100:
		limit = 20
	case total > 500:
		limit = 50
	case total > 1000:
		limit = 100
	}
	pages := []int64{}
	var lastpage = (total / limit) + 1
	for i := int64(1); i <= lastpage; i++ {
		pages = append(pages, i)
	}
	if page > lastpage {
		page = lastpage
	}

	// chunk
	var begin = (page - 1) * limit
	var end = begin + limit
	if end > total {
		end = total
	}

	lists, err := listLists()
	if err != nil {
		_error(w, err)
		return
	}

	res := newResponse(r, ps)
	res.Section = "library"
	res.Medias = medias[begin:end]
	res.Lists = lists
	res.Page = page
	res.Pages = pages
	res.Query = query
	res.Limit = limit
	res.Total = total
	res.GrandTotal = grandTotal
	html(w, "library.html", res)
}

//
// Media
//

func thumbnailMedia(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	media, err := FindMedia(ps.ByName("media"))
	if err != nil {
		_error(w, err)
		return
	}

	img, err := imaging.Open(media.imageFile())
	if err != nil {
		_error(w, err)
		return
	}

	img = imaging.Resize(img, 320, 0, imaging.Lanczos)

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Vary", "Accept-Encoding")
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", 7*86400))
	if err := imaging.Encode(w, img, imaging.JPEG); err != nil {
		_error(w, err)
		return
	}
}

func viewMedia(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	media, err := FindMedia(ps.ByName("media"))
	if err != nil {
		_error(w, err)
		return
	}

	res := newResponse(r, ps)
	res.Media = media
	res.Section = "library"
	res.Section = "view"
	html(w, "view.html", res)
}

func deleteMedia(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := DeleteMedia(ps.ByName("media")); err != nil {
		_error(w, err)
		return
	}
	redirect(w, r, "/library?p=%s&q=%s&message=mediadeleted", r.FormValue("p"), r.FormValue("q"))
}

func downloadMedia(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filename := filepath.Join(datadir, ps.ByName("filename"))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filepath.Base(filename)))
	http.ServeFile(w, r, filename)
}

func streamMedia(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filename := filepath.Join(datadir, ps.ByName("filename"))
	if id := ps.ByName("list"); id != "" {
		if _, err := findList(id); err != nil {
			_error(w, err)
			return
		}
	}
	if strings.HasSuffix(filename, ".m4a") {
		w.Header().Set("Content-Type", "video/mp4")
	}
	http.ServeFile(w, r, filename)
}

//
// Archiver
//

func archiverJobs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res := newResponse(r, ps)
	res.ActiveMedias = ActiveMedias()
	res.QueuedMedias = QueuedMedias()
	html(w, "jobs.html", res)
}

func archiverSave(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	source := fmt.Sprintf("https://www.youtube.com/v?id=%s", id)

	vinfo, err := ytdl.GetVideoInfoFromID(id)
	if err != nil {
		_error(w, err)
		return
	}

	media, err := NewMedia(vinfo.ID, vinfo.Author, vinfo.Title, int64(vinfo.Duration.Seconds()), source)
	if err != nil {
		_error(w, err)
		return
	}
	logger.Infof("created new media %q %q", media.ID, media.Title)

	archive.Add(id, source)
	toJSON(w, "OK")
}

func archiverCancel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	archive.Remove(ps.ByName("id"))
	redirect(w, r, "/import?message=savecancelled")
}

func deleteList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list, err := findList(ps.ByName("id"))
	if err != nil {
		_error(w, err)
		return
	}
	if err := DeleteList(list.ID); err != nil {
		_error(w, err)
		return
	}
	redirect(w, r, "/?message=playlistdeleted")
}

func podcastList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list, err := findList(ps.ByName("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	proto := r.Header.Get("X-Forwarded-Proto")
	if proto == "" {
		proto = "https"
	}
	baseurl := fmt.Sprintf("%s://%s%s", proto, httpHost, httpPrefix)

	p := podcast.New(list.Title, baseurl, list.Title, &list.Created, &list.Modified)
	p.AddAuthor(httpHost, "streamlist@"+httpHost)
	p.AddImage(baseurl + "/logo.png")

	for _, media := range list.Medias {
		typ := podcast.M4V
		ext := "m4a"
		filename := media.audioFile()

		fileInfo, err := os.Stat(filename)
		if err != nil {
			logger.Error(err)
			continue
		}

		streamurl := fmt.Sprintf("%s/stream/%s/%s.%s", baseurl, list.ID, media.ID, ext)

		item := podcast.Item{
			Title:       fmt.Sprintf("%s - %s", media.Title, media.Author),
			Description: fmt.Sprintf("%s\n\n%s", media.Description, media.Created),
			PubDate:     &media.Created,
		}
		item.AddEnclosure(streamurl, typ, fileInfo.Size())
		if _, err := p.AddItem(item); err != nil {
			_error(w, err)
			return
		}
	}
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	if err := p.Encode(w); err != nil {
		_error(w, err)
	}
}

func m3uList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list, err := findList(ps.ByName("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	ext := ".m4a"

	w.Header().Set("Content-Type", "application/mpegurl")
	fmt.Fprintf(w, "#EXTM3U\n")
	for _, media := range list.Medias {
		fmt.Fprintf(w, "#EXTINF:%d,%s\n", media.Length, media.Title)
		proto := r.Header.Get("X-Forwarded-Proto")
		if proto == "" {
			proto = "https"
		}
		fmt.Fprintf(w, "%s://%s%s/stream/%s/%s%s\n", proto, httpHost, httpPrefix, list.ID, media.ID, ext)
	}
}

func playList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list, err := findList(ps.ByName("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	res := newResponse(r, ps)
	res.Section = "play"
	res.List = list
	var medias []*Media
	db.Model(&list).Related(&medias, "Medias")
	res.Medias = medias
	html(w, "play.html", res)
}

func createList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == "GET" {
		res := newResponse(r, ps)
		res.Section = "create"
		html(w, "create.html", res)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))

	if title == "" {
		redirect(w, r, "/create")
		return
	}

	_, err := newList(title)
	if err != nil {
		_error(w, err)
		return
	}
	redirect(w, r, "/library?message=playlistadded")
}

func removeMediaList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	media, err := FindMedia(ps.ByName("media"))
	if err != nil {
		_error(w, err)
		return
	}
	list, err := findList(ps.ByName("list"))
	if err != nil {
		_error(w, err)
		return
	}
	if err := list.removeMedia(media); err != nil {
		_error(w, err)
		return
	}
	redirect(w, r, "/edit/%s", list.ID)
}

func addMediaList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	media, err := FindMedia(ps.ByName("media"))
	if err != nil {
		_error(w, err)
		return
	}

	list, err := findList(ps.ByName("list"))
	if err != nil {
		_error(w, err)
		return
	}

	list.addMedia(media)
	toJSON(w, "OK")
}

func shuffleList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list, err := findList(ps.ByName("id"))
	if err != nil {
		_error(w, err)
		return
	}
	if err := list.shuffleMedia(); err != nil {
		_error(w, err)
		return
	}

	redirect(w, r, "/play/%s", list.ID)
}

func editList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list, err := findList(ps.ByName("id"))
	if err != nil {
		_error(w, err)
		return
	}

	res := newResponse(r, ps)
	res.Section = "edit"
	res.List = list
	var medias []*Media
	db.Model(&list).Related(&medias, "Medias")
	res.Medias = medias
	html(w, "edit.html", res)
}

func staticAsset(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	serveAsset(w, r, ps.ByName("path"))
}

func logo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	serveAsset(w, r, "/logo.png")
}

func serveAsset(w http.ResponseWriter, r *http.Request, filename string) {
	path := "static" + filename

	b, err := Asset(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fi, err := AssetInfo(path)
	if err != nil {
		_error(w, err)
		return
	}
	http.ServeContent(w, r, path, fi.ModTime(), bytes.NewReader(b))
}

//
// API
//
func v1status(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// localhost only.
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if ip != "::1" && ip != "127.0.0.1" {
		http.NotFound(w, r)
		return
	}
	status := "idle"
	if len(QueuedMedias()) > 0 || len(ActiveMedias()) > 0 {
		status = "busy"
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%s\n", status)
}

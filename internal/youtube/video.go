package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rylio/ytdl"
	//log "github.com/Sirupsen/logrus"
)

var (
	ytplayerRegexp = regexp.MustCompile(`ytplayer\.config\s*=\s*(.*?)\s*;\s*ytplayer\.load`)
	fixurlRegexp   = regexp.MustCompile(`\,[^=]+=.*$`)
)

// Video ...
type Video struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Thumbnail string    `json:"thumbnail"`
	Views     int64     `json:"views"`
	Length    int64     `json:"length"`
	Rating    float64   `json:"rating"`
	Timestamp time.Time `json:"timestamp"`
	Streams   []Stream
}

// Filename ...
func (v Video) Filename(dir string) string {
	return filepath.Join(dir, v.ID+".mp4")
}

// ThumbnailFilename ...
func (v Video) ThumbnailFilename(dir string) string {
	return filepath.Join(dir, v.ID+".jpg")
}

// jsonFilename ...
func (v Video) jsonFilename(dir string) string {
	return filepath.Join(dir, v.ID+".json")
}

// Transcode ...
func (v Video) Transcode(ctx context.Context, dir string) error {
	filename := v.Filename(dir)
	tmpname := filename + ".encoding"
	if err := os.Rename(filename, tmpname); err != nil {
		return err
	}
	output, err := exec.CommandContext(ctx,
		"/usr/bin/ffmpeg",
		"-i", tmpname,
		"-vn", "-c:a", "aac",
		"-strict", "experimental",
		"-movflags", "+faststart", filename,
	).CombinedOutput()
	if err != nil {
		os.Remove(tmpname)
		os.Remove(filename)
		return fmt.Errorf("ffmpeg failed: %s %s", err, string(output))
	}
	os.Remove(tmpname)
	return nil
}

// Download ...
func (v Video) Download(ctx context.Context, dir string) error {
	if len(v.Streams) == 0 {
		return fmt.Errorf("no streams")
	}
	// pick first stream (seems to always be best one)
	//stream := v.Streams[0]

	vi, err := ytdl.GetVideoInfoFromID(v.ID)
	if err != nil {
		return err
	}

	format := vi.Formats[0]
	//for _, f := range vi.Formats { log.Printf("availabe format is %+v", f) }
	dlurl, err := vi.GetDownloadURL(format)
	if err != nil {
		return err
	}

	// download thumbnail
	if err := download(ctx, v.Thumbnail, v.ThumbnailFilename(dir)); err != nil {
		return err
	}
	// download video
	if err := download(ctx, dlurl.String(), v.Filename(dir)); err != nil {
		return err
	}
	// transcode
	if err := v.Transcode(ctx, dir); err != nil {
		return err
	}
	return v.write(v.jsonFilename(dir))
}

func (v Video) write(filename string) error {
	// v.Title = strings.

	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}

func download(ctx context.Context, rawurl, filename string) error {
	// open file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	// request file
	res, err := GET(ctx, rawurl)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("download %s failed: %s", rawurl, http.StatusText(res.StatusCode))
	}

	// write to file
	_, err = io.Copy(f, res.Body)
	// remove file if we failed
	if err != nil {
		os.Remove(filename)
	}
	return err
}

// GetVideo ...
func GetVideo(rawid string) (Video, error) {
	id := rawid
	if strings.HasPrefix(rawid, "http") {
		u, err := url.Parse(rawid)
		if err != nil {
			return Video{}, err
		}
		switch strings.TrimPrefix(u.Host, "www.") {
		case "youtube.com":
			id = u.Query().Get("v")
			if id == "" {
				id = strings.TrimPrefix(u.Path, "/v/")
			}
		case "youtu.be":
			id = u.Path
		}
	}

	if id == "" {
		return Video{}, fmt.Errorf("invalid video ID %q", rawid)
	}

	res, err := GET(nil, "https://www.youtube.com/watch?v="+id)
	if err != nil {
		return Video{}, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Video{}, err
	}

	matches := ytplayerRegexp.FindSubmatch(b)
	if len(matches) != 2 {
		return Video{}, fmt.Errorf("failed to extract ytconfig")
	}
	ytconfig := make(map[string]interface{})

	if err := json.Unmarshal(matches[1], &ytconfig); err != nil {
		return Video{}, fmt.Errorf("failed to unmarshal ytconfig: %s", err)
	}

	args, ok := ytconfig["args"].(map[string]interface{})
	if !ok {
		return Video{}, fmt.Errorf("missing args in ytconfig")
	}

	for _, k := range []string{
		"video_id",
		"title",
		"author",
		"iurlmq",
		"view_count",
		"avg_rating",
		"length_seconds",
		"timestamp",
		"url_encoded_fmt_stream_map",
	} {
		if _, ok := args[k]; !ok {
			return Video{}, fmt.Errorf("failed to extract %q", k)
		}
	}

	// ID
	//id := args["video_id"].(string)

	// Title
	title := args["title"].(string)

	// Author
	author := args["author"].(string)

	// Thumbnail - thumbnail_url -> iurlhq -> iurlmq
	thumbnail := args["iurlmq"].(string)

	// Views
	views, err := strconv.ParseInt(args["view_count"].(string), 10, 64)
	if err != nil {
		return Video{}, err
	}

	// Length
	length, err := strconv.ParseInt(args["length_seconds"].(string), 10, 64)
	if err != nil {
		return Video{}, err
	}

	// Rating
	rating, err := strconv.ParseFloat(args["avg_rating"].(string), 64)
	if err != nil {
		return Video{}, err
	}

	// Timestamp
	// ts, err := strconv.ParseInt(args["timestamp"].(string), 10, 64)
	// if err != nil { log.Fatal(err) }
	// timestamp := time.Unix(ts, 0)

	// Streams
	streamMap, err := url.ParseQuery(args["url_encoded_fmt_stream_map"].(string))
	if err != nil {
		return Video{}, err
	}
	urls := make([]string, 0, 100)
	itags := make([]int, 0, 100)

	for key, values := range streamMap {
		for i, v := range values {
			// TODO: figure out why we're getting ",quality=medium" and ",itag=42" on end of values
			// probably because they use CSV AND & separated query values?
			v = fixurlRegexp.ReplaceAllString(v, "")
			switch key {
			case "url":
				urls = append(urls[:i], append([]string{v}, urls[i:]...)...)
			case "itag":
				itag, err := strconv.Atoi(v)
				if err != nil {
					return Video{}, fmt.Errorf("parsing itag failed: %s", err)
				}
				itags = append(itags[:i], append([]int{itag}, itags[i:]...)...)
			}
		}
	}

	//if len(urls) != len(itags) { return nil, fmt.Errorf("mismatched urls %+v and itags %+v", urls, itags) }
	count := len(urls)
	if len(itags) < count {
		count = len(itags)
	}

	var streams []Stream
	for i := 0; i < count; i++ {
		s, err := newStream(itags[i], urls[i])
		if err != nil {
			return Video{}, err
		}
		streams = append(streams, s)
	}

	return Video{
		ID:        id,
		Title:     title,
		Author:    author,
		Thumbnail: thumbnail,
		Views:     views,
		Length:    length,
		Rating:    rating,
		Streams:   streams,
	}, nil
}

// GET ...
func GET(ctx context.Context, rawurl string) (*http.Response, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	if ctx != nil {
		client.Timeout = 24 * time.Hour
	}

	req, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		return nil, err
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36")

	return client.Do(req)
}

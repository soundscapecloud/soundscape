package youtube

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
)

// SetDebug ...
func SetDebug() {
	log.SetLevel(log.DebugLevel)
}

// Search ...
func Search(query string) ([]Video, error) {
	u, err := url.Parse("https://www.youtube.com/results")
	if err != nil {
		return nil, err
	}
	q := &url.Values{}
	q.Add("search_query", query)
	u.RawQuery = q.Encode()

	res, err := GET(nil, u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	ytdataPattern := regexp.MustCompile(`window\["ytInitialData"\]\s*=\s*(?P<ydata>\{.*?\})\s*;\s*\n`)

	var ytd *ytdata
	var finderr error
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		matches := ytdataPattern.FindStringSubmatch(s.Text())
		if len(matches) < 2 {
			return
		}
		match := matches[1]

		y := ytdata{}
		if err := json.Unmarshal([]byte(match), &y); err != nil {
			finderr = fmt.Errorf("failed to extract ytdata: %s", err)
			return
		}
		ytd = &y
	})
	if finderr != nil {
		return nil, finderr
	}
	if ytd == nil {
		return nil, fmt.Errorf("failed to find ytdata")
	}

	// LOL
	container := ytd.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents
	if len(container) == 0 {
		return nil, fmt.Errorf("failed to find contents containter in ytdata")
	}
	contents := container[0].ItemSectionRenderer.Contents

	var videos []Video
	for _, c := range contents {
		vr := c.VideoRenderer

		// id
		id := vr.VideoID
		if id == "" {
			log.Debugf("failed to find VideoID in ytdata")
			continue
		}

		// title
		title := vr.Title.SimpleText
		if title == "" {
			return nil, fmt.Errorf("failed to find Title in ytdata")
		}

		// length
		length, err := func() (int64, error) {
			videotime := vr.LengthText.SimpleText
			f := strings.Split(videotime, ":")
			switch len(f) {
			case 2:
				videotime = fmt.Sprintf("%sm%ss", f[0], f[1])
			case 3:
				videotime = fmt.Sprintf("%sh%sm%ss", f[0], f[1], f[2])
			default:
				return 0, fmt.Errorf("invalid length text in ytdata")
			}
			d, err := time.ParseDuration(videotime)
			if err != nil {
				return 0, err
			}
			return int64(d.Seconds()), nil
		}()
		if err != nil {
			log.Debug(err)
			continue
		}
		if title == "" {
			log.Debugf("failed to find Length in ytdata")
			continue
		}

		// thumbnail
		thumbnail := fmt.Sprintf("https://i.ytimg.com/vi/%s/hqdefault.jpg", id)

		videos = append(videos, Video{
			ID:        id,
			Title:     title,
			Thumbnail: thumbnail,
			Length:    length,
		})
	}
	return videos, nil
}

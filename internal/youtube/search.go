package youtube

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	//log "github.com/Sirupsen/logrus"
)

func Search(query string, max int) ([]Video, error) {
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

	var videos []Video
	var finderr error
	doc.Find(".yt-lockup-dismissable").Each(func(i int, q *goquery.Selection) {
		// title
		title, ok := q.Find(".yt-lockup-title > a").Attr("title")
		if !ok {
			finderr = err
			return
		}

		// length
		length, err := func() (int64, error) {
			videotime := q.Find(".video-time").Text()
			f := strings.Split(videotime, ":")
			switch len(f) {
			case 2:
				videotime = fmt.Sprintf("%sm%ss", f[0], f[1])
			case 3:
				videotime = fmt.Sprintf("%sh%sm%ss", f[0], f[1], f[2])
			default:
				return 0, fmt.Errorf("invalid video-time")
			}
			d, err := time.ParseDuration(videotime)
			if err != nil {
				return 0, err
			}
			return int64(d.Seconds()), nil
		}()
		if err != nil {
			finderr = err
			return
		}

		// id
		link, ok := q.Find(".yt-lockup-title > a").Attr("href")
		if !ok {
			finderr = fmt.Errorf("missing link")
			return
		}
		id, err := func() (string, error) {
			l, err := url.ParseRequestURI(link)
			if err != nil {
				return "", err
			}
			v := l.Query().Get("v")
			if v == "" {
				return "", fmt.Errorf("missing ID in link %q", link)
			}
			return v, nil
		}()
		if err != nil {
			finderr = err
			return
		}

		// thumbnail
		thumbnail := fmt.Sprintf("https://i.ytimg.com/vi/%s/hqdefault.jpg", id)

		videos = append(videos, Video{
			ID:        id,
			Title:     title,
			Thumbnail: thumbnail,
			Length:    length,
		})
	})

	/*
		var videos []Video
		for _, id := range ids {
			if len(videos) >= max {
				break
			}
			v, err := GetVideo(id)
			if err != nil {
				log.Errorf("get %q failed: %s", id, err)
				continue
			}
			videos = append(videos, v)
		}
	*/

	return videos, nil
}

/*
   // title
   title, ok := q.Find(".yt-lockup-title > a").Attr("title")
   if !ok {
       errs = append(errs, fmt.Errorf("missing title"))
       return
   }

   // duration
   duration, err := func() (time.Duration, error) {
       viewcount := q.Find(".video-time").Text()
       f := strings.Split(viewcount, ":")
       switch len(f) {
       case 2:
           viewcount = fmt.Sprintf("%sm%ss", f[0], f[1])
       case 3:
           viewcount = fmt.Sprintf("%sh%sm%ss", f[0], f[1], f[2])
       }
       return time.ParseDuration(viewcount)
   }()
   if err != nil {
       errs = append(errs, err)
       return
   }

   // created and views
   ago := ""
   viewcount := ""
   q.Find(".yt-lockup-meta-info > li").Each(func(_ int, m *goquery.Selection) {
       if strings.Contains(m.Text(), "views") {
           viewcount = m.Text()
           viewcount = strings.Replace(viewcount, ",", "", -1)
           viewcount = strings.Replace(viewcount, " views", "", -1)
       }
       if strings.Contains(m.Text(), "ago") {
           ago = m.Text()
       }
   })
   created, err := ago2time(ago)
   if err != nil {
       errs = append(errs, err)
       return
   }

   views, err := strconv.ParseInt(viewcount, 10, 64)
   if err != nil {
       errs = append(errs, err)
       return
   }
*/

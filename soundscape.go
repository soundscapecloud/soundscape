package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var ErrMediaNotFound = errors.New("media not found")

//
// Media
//
type Media struct {
	ID          string    `json:"id"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Length      int64     `json:"length"` // In seconds
	Source      string    `json:"source"`
	Modified    time.Time `json:"modified"`
	Created     time.Time `json:"created"`
}

func mediaFile(id string) string {
	if id == "" {
		panic("invalid media id")
	}
	return filepath.Join(datadir, id+".media")
}

func NewMedia(id, author, title string, length int64, source string) (*Media, error) {
	media := &Media{
		ID:       id,
		Author:   author,
		Title:    title,
		Length:   length,
		Source:   source,
		Modified: time.Now(),
		Created:  time.Now(),
	}
	return media, media.Save()
}

func QueuedMedias() []*Media {
	var medias []*Media
	for _, id := range archive.QueuedJobs() {
		m, err := loadMedia(id)
		if err != nil {
			logger.Warnf("failed to find media for job %q", id)
			continue
		}
		medias = append(medias, m)
	}
	return medias
}

func ActiveMedias() []*Media {
	var medias []*Media
	for _, id := range archive.ActiveJobs() {
		m, err := loadMedia(id)
		if err != nil {
			logger.Warnf("failed to find media for job %q", id)
			continue
		}
		medias = append(medias, m)
	}
	return medias
}

func DeleteMedia(id string) error {
	media, err := FindMedia(id)
	if err != nil {
		return nil
	}

	// Remove all list references to this media.
	lists, err := ListLists()
	if err != nil {
		return err
	}
	for _, l := range lists {
		if err := l.RemoveMedia(media); err != nil {
			return err
		}
	}

	// Remove all media files.
	files := []string{
		media.ImageFile(),
		media.VideoFile(),
		media.AudioFile(),
		media.File(),
	}
	for _, f := range files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			continue
		}
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}

func DeleteList(id string) error {
	list, err := FindList(id)
	if err != nil {
		return err
	}
	return os.Remove(list.File())
}

func FindMedia(id string) (*Media, error) {
	medias, err := ListMedias()
	if err != nil {
		return nil, err
	}
	for _, m := range medias {
		if m.ID == id {
			return m, nil
		}
	}
	return nil, ErrMediaNotFound
}

func loadMedia(id string) (*Media, error) {
	b, err := ioutil.ReadFile(mediaFile(id))
	if err != nil {
		return nil, err
	}
	var media Media
	return &media, json.Unmarshal(b, &media)
}

func ListMedias() ([]*Media, error) {
	files, err := ioutil.ReadDir(datadir)
	if err != nil {
		return nil, err
	}
	sort.Slice(files, func(i, j int) bool {
		return files[j].ModTime().Before(files[i].ModTime())
	})

	var medias []*Media
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".media") {
			continue
		}
		// media must exist.
		m, err := loadMedia(strings.TrimSuffix(f.Name(), ".media"))
		if err != nil {
			return nil, err
		}
		// must have an image file.
		if !m.HasImage() {
			continue
		}
		// must have an audio file (otherwise it's not finished transcoding)
		if !m.HasAudio() {
			continue
		}
		medias = append(medias, m)
	}
	return medias, nil
}

func (m Media) Save() error {
	b, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return err
	}
	return Overwrite(m.File(), b, 0644)
}

func (m Media) File() string {
	return mediaFile(m.ID)
}

func (m Media) ImageFile() string {
	return filepath.Join(datadir, m.ID+".jpg")
}

func (m Media) VideoFile() string {
	return filepath.Join(datadir, m.ID+".mp4")
}

func (m Media) AudioFile() string {
	return filepath.Join(datadir, m.ID+".m4a")
}

func (m Media) HasImage() bool {
	_, err := os.Stat(m.ImageFile())
	return err == nil
}

func (m Media) HasVideo() bool {
	_, err := os.Stat(m.VideoFile())
	return err == nil
}

func (m Media) HasAudio() bool {
	_, err := os.Stat(m.AudioFile())
	return err == nil
}

//
// List
//
type List struct {
	ID    string `json:"id"`
	Title string `json:"title"`

	Medias []*Media `json:"medias"`

	Modified time.Time `json:"modified"`
	Created  time.Time `json:"created"`
}

func listFile(id string) string {
	if id == "" {
		panic("invalid list id")
	}
	return filepath.Join(datadir, id+".playlist")
}

func NewList(title string) (*List, error) {
	id, err := RandomNumber()
	if err != nil {
		return nil, err
	}
	list := &List{
		ID:       fmt.Sprintf("%d", id),
		Title:    title,
		Modified: time.Now(),
		Created:  time.Now(),
	}
	return list, list.Save()
}

func (l *List) File() string {
	return listFile(l.ID)
}

func (l *List) Save() error {
	b, err := json.MarshalIndent(l, "", "    ")
	if err != nil {
		return err
	}
	l.Modified = time.Now()
	return Overwrite(l.File(), b, 0644)
}

func (l *List) HasMedia(media *Media) bool {
	for _, m := range l.Medias {
		if m.ID == media.ID {
			return true
		}
	}
	return false
}

func (l *List) TotalLength() (total int64) {
	for _, m := range l.Medias {
		total += m.Length
	}
	return total
}

func (l *List) ShuffleMedia() error {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	var medias []*Media
	for _, i := range r.Perm(len(l.Medias)) {
		medias = append(medias, l.Medias[i])
	}
	l.Medias = medias
	return l.Save()
}

func (l *List) AddMedia(media *Media) error {
	l.Medias = append(l.Medias, media)
	return l.Save()
}

func (l *List) RemoveMedia(media *Media) error {
	if !l.HasMedia(media) {
		return nil
	}
	var medias []*Media
	for _, m := range l.Medias {
		if m.ID == media.ID {
			continue
		}
		medias = append(medias, m)
	}
	l.Medias = medias
	return l.Save()
}

func FindList(id string) (*List, error) {
	b, err := ioutil.ReadFile(listFile(id))
	if err != nil {
		return nil, err
	}
	var list List
	return &list, json.Unmarshal(b, &list)
}

func ListLists() ([]*List, error) {
	files, err := ioutil.ReadDir(datadir)
	if err != nil {
		return nil, err
	}
	// sort.Slice(files, func(i, j int) bool {
	// 	return files[j].ModTime().Before(files[i].ModTime())
	// })

	var lists []*List
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".playlist") {
			continue
		}
		l, err := FindList(strings.TrimSuffix(f.Name(), ".playlist"))
		if err != nil {
			return nil, err
		}
		lists = append(lists, l)
	}
	sort.Slice(lists, func(i, j int) bool {
		return lists[j].Created.Before(lists[i].Created)
	})
	return lists, nil
}

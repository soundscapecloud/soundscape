package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	//"sort"
	//"strings"
	"sync"
	"time"
)

var errMediaNotFound = errors.New("media not found")

// Config is the global config
type Config struct {
	sync.RWMutex
	filename string

	// Settings
	Volume float32 `json:"volume"`
}

// NewConfig returns a new Config
func NewConfig(filename string) (*Config, error) {
	filename = filepath.Join(datadir, filename)
	c := &Config{filename: filename}
	b, err := ioutil.ReadFile(filename)

	// Default for new config
	if os.IsNotExist(err) {
		c.Volume = 0.2
		return c, c.Save()
	}
	if err != nil {
		return nil, err
	}

	// Open existing config
	if err := json.Unmarshal(b, c); err != nil {
		return nil, err
	}
	return c, nil
}

// Get return the current config
func (c *Config) Get() Config {
	c.RLock()
	defer c.RUnlock()

	return Config{
		Volume: c.Volume,
	}
}

// SetVolume modify configured volume
func (c *Config) SetVolume(n float32) error {
	c.Lock()
	c.Volume = n
	c.Unlock()
	return c.Save()
}

// Save saves config to file
func (c *Config) Save() error {
	c.RLock()
	defer c.RUnlock()

	b, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	return overwrite(c.filename, b, 0644)
}

// User ...
type User struct {
	ID       uint
	Username string
	Password string
	Role     string
}

// Media represent a media in the library
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

// NewMedia return a new Media
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
	return media, media.save()
}

// QueuedMedias return queueud media list
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

// ActiveMedias return active media list
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

// DeleteMedia removes mediafrom library
func DeleteMedia(id string) error {
	media, err := FindMedia(id)
	if err != nil {
		return nil
	}

	// Delete associations
	db.Exec("DELETE FROM list_media WHERE media_id = ?", media.ID)
	db.Delete(&media)

	// Remove all media files.
	files := []string{
		media.imageFile(),
		media.videoFile(),
		media.audioFile(),
		media.file(),
	}
	for _, f := range files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			continue
		}
		if err := os.Remove(f); err != nil {
			return err
		}
	}

	return db.Error
}

// DeleteList removes a playlist
func DeleteList(id string) error {
	list, err := findList(id)
	if err != nil {
		return err
	}
	// Delete associations
	db.Exec("DELETE FROM list_media WHERE list_id = ?", list.ID)
	db.Delete(&list)
	return db.Error
	//return os.Remove(list.file())
}

// FindMedia search media in library
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
	return nil, errMediaNotFound
}

// LoadMedia reads media file
func loadMedia(id string) (*Media, error) {
	var media Media
	db.First(&media, "ID = ?", id)
	return &media, db.Error
}

// ListMedias list medias in library
func ListMedias() ([]*Media, error) {
	/*files, err := ioutil.ReadDir(datadir)
	if err != nil {
		return nil, err
	}
	sort.Slice(files, func(i, j int) bool {
		return files[j].ModTime().Before(files[i].ModTime())
	})*/

	var mediasBDD []*Media
	var medias []*Media
	db.Order("modified desc").Find(&mediasBDD)
	for _, m := range mediasBDD {
		// must have an image file.
		if !m.hasImage() {
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

func (m Media) save() error {
	return db.Create(&m).Error
}

func (m Media) file() string {
	return mediaFile(m.ID)
}

func (m Media) imageFile() string {
	return filepath.Join(datadir, m.ID+".jpg")
}

func (m Media) videoFile() string {
	return filepath.Join(datadir, m.ID+".mp4")
}

func (m Media) audioFile() string {
	return filepath.Join(datadir, m.ID+".m4a")
}

func (m Media) hasImage() bool {
	_, err := os.Stat(m.imageFile())
	return err == nil
}

func (m Media) hasVideo() bool {
	_, err := os.Stat(m.videoFile())
	return err == nil
}

// HasAudio ...
func (m Media) HasAudio() bool {
	_, err := os.Stat(m.audioFile())
	return err == nil
}

// List represent a playlist
type List struct {
	ID    string `json:"id"`
	Title string `json:"title"`

	Medias []*Media `json:"medias" gorm:"many2many:list_media;AssociationForeignKey:ID;ForeignKey:ID"`

	Modified time.Time `json:"modified"`
	Created  time.Time `json:"created"`
}

func listFile(id string) string {
	if id == "" {
		panic("invalid list id")
	}
	return filepath.Join(datadir, id+".playlist")
}

func newList(title string) (*List, error) {
	id, err := randomNumber()
	if err != nil {
		return nil, err
	}
	list := &List{
		ID:       fmt.Sprintf("%d", id),
		Title:    title,
		Modified: time.Now(),
		Created:  time.Now(),
	}
	return list, list.save()
}

func (l *List) file() string {
	return listFile(l.ID)
}

func (l *List) save() error {
	l.Modified = time.Now()
	db.Where(List{ID: l.ID}).Assign(&l).FirstOrCreate(&l)
	return db.Error
}

// HasMedia ...
func (l *List) HasMedia(media *Media) bool {
	var medias []Media
	db.Model(&l).Related(&medias, "Medias")
	for _, m := range medias {
		if m.ID == media.ID {
			return true
		}
	}
	return false
}

// TotalLength ...
func (l *List) TotalLength() (total int64) {
	var medias []Media
	db.Model(&l).Related(&medias, "Medias")
	for _, m := range medias {
		total += m.Length
	}
	return total
}

// MediasCount ...
func (l *List) MediasCount() int {
	var medias []Media
	db.Model(&l).Related(&medias, "Medias")
	return len(medias)
}

func (l *List) shuffleMedia() error {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	var medias []*Media
	for _, i := range r.Perm(len(l.Medias)) {
		medias = append(medias, l.Medias[i])
	}
	l.Medias = medias
	return l.save()
}

func (l *List) addMedia(media *Media) error {
	l.Medias = append(l.Medias, media)
	return l.save()
}

func (l *List) removeMedia(media *Media) error {
	if !l.HasMedia(media) {
		return nil
	}
	db.Model(&l).Association("Medias").Delete(media)
	return db.Error
}

func findList(id string) (*List, error) {
	var list List
	db.First(&list, "ID = ?", id)
	return &list, db.Error
}

func listLists() ([]*List, error) {
	/*files, err := ioutil.ReadDir(datadir)
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
		l, err := findList(strings.TrimSuffix(f.Name(), ".playlist"))
		if err != nil {
			return nil, err
		}
		lists = append(lists, l)
	}
	sort.Slice(lists, func(i, j int) bool {
		return lists[j].Created.Before(lists[i].Created)
	})
	return lists, nil*/
	var lists []*List
	db.Find(&lists)
	return lists, db.Error
}

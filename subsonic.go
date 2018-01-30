package main

import (
	"encoding/xml"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"time"
)

// Credit: https://github.com/mdlayher/wavepipe/blob/master/subsonic/subsonic.go

const (
	// XMLName is the top-level name of a Subsonic XML response
	SubsonicXMLName = "subsonic-response"

	// XMLNS is the XML namespace of a Subsonic XML response
	SubsonicXMLNS = "http://subsonic.org/restapi"

	// Version is the emulated Subsonic API version
	SubsonicVersion = "1.8.0"
)

type SubsonicResponse struct {
	// Top-level container name
	XMLName xml.Name `xml:"subsonic-response"`

	// Attributes which are always present
	XMLNS   string `xml:"xmlns,attr"`
	Status  string `xml:"status,attr"`
	Version string `xml:"version,attr"`

	// Error, returned on failures
	SubError *SubsonicError

	// Nested data

	// getAlbum.view
	Album []SubsonicAlbum `xml:"album"`

	// getAlbumList2.view
	AlbumList2 *SubsonicAlbumList2

	// getIndexes.view
	Indexes *SubsonicIndexes

	// getLicense.view
	License *SubsonicLicense `xml:"license"`

	// getMusicDirectory.view
	MusicDirectory *SubsonicMusicDirectory

	// getMusicFolders.view
	MusicFolders *SubsonicMusicFolders

	// getPlaylists.view
	Playlists *SubsonicPlaylists

	// getPlaylist.view
	Playlist *SubsonicPlaylist

	// getLyrics.view
	Lyrics *SubsonicLyrics

	// getRandomSongs.view
	RandomSongs *SubsonicRandomSongs

	// getStarred.view
	Starred *SubsonicStarred `xml:"starred"`
}

// SubsonicError contains a Subsonic error, with status code and message
type SubsonicError struct {
	XMLName xml.Name `xml:"error,omitempty"`

	Code    int    `xml:"code,attr"`
	Message string `xml:"message,attr"`
}

// SubsonicArtist represents an emulated Subsonic artist
type SubsonicArtist struct {
	XMLName xml.Name `xml:"artist,omitempty"`

	// Subsonic fields
	Name string `xml:"name,attr"`
	ID   string `xml:"id,attr"`
}

// SubsonicAlbum represents an emulated Subsonic album
type SubsonicAlbum struct {
	// Subsonic fields
	ID        int    `xml:"id,attr"`
	Name      string `xml:"name,attr"`
	Artist    string `xml:"artist,attr"`
	ArtistID  int    `xml:"artistId,attr"`
	CoverArt  string `xml:"coverArt,attr"`
	SongCount int    `xml:"songCount,attr"`
	Duration  int    `xml:"duration,attr"`
	Created   string `xml:"created,attr"`

	// Nested data

	// getAlbum.view
	Songs []SubsonicSong `xml:"song"`
}

// SubsonicRandomSongs contains a random list of emulated Subsonic songs
type SubsonicRandomSongs struct {
	// Container name
	XMLName xml.Name `xml:"randomSongs,omitempty"`

	// Songs
	Songs []SubsonicSong `xml:"song"`
}

// SubsonicLyrics represents a Subsonic lyrics
type SubsonicLyrics struct {
	// Container name
	XMLName xml.Name `xml:"lyrics,omitempty"`

	Artist string `xml:"artist,attr,omitempty"`
	Title  string `xml:"title,attr,omitempty"`
}

// SubsonicLicense represents a Subsonic license
type SubsonicLicense struct {
	XMLName xml.Name `xml:"license,omitempty"`

	Valid bool   `xml:"valid,attr"`
	Email string `xml:"email,attr"`
	Key   string `xml:"key,attr"`
	Date  string `xml:"date,attr"`
}

// SubsonicSong represents an emulated Subsonic song
type SubsonicSong struct {
	ID          int    `xml:"id,attr"`
	Parent      int    `xml:"parent,attr"`
	Title       string `xml:"title,attr"`
	Album       string `xml:"album,attr"`
	Artist      string `xml:"artist,attr"`
	IsDir       bool   `xml:"isDir,attr"`
	CoverArt    string `xml:"coverArt,attr"`
	Created     string `xml:"created,attr"`
	Duration    int    `xml:"duration,attr"`
	BitRate     int    `xml:"bitRate,attr"`
	Track       int    `xml:"track,attr"`
	DiscNumber  int    `xml:"discNumber,attr"`
	Year        int    `xml:"year,attr"`
	Genre       string `xml:"genre,attr"`
	Size        int64  `xml:"size,attr"`
	Suffix      string `xml:"suffix,attr"`
	ContentType string `xml:"contentType,attr"`
	IsVideo     bool   `xml:"isVideo,attr"`
	Path        string `xml:"path,attr"`
	AlbumID     int    `xml:"albumId,attr"`
	ArtistID    int    `xml:"artistId,attr"`
	Type        string `xml:"type,attr"`
}

// SubsonicAlbumList2 contains a list of emulated Subsonic albums, by tags
type SubsonicAlbumList2 struct {
	// Container name
	XMLName xml.Name `xml:"albumList2,omitempty"`

	// Albums
	Albums []SubsonicAlbum `xml:"album"`
}

// SubsonicMusicFolder represents an emulated Subsonic music folder
type SubsonicMusicFolder struct {
	ID   int    `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

// SubsonicMusicDirectory contains a list of emulated Subsonic music folders
type SubsonicMusicDirectory struct {
	// Container name
	XMLName xml.Name `xml:"directory,omitempty"`

	// Attributes
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`

	Children []SubsonicChild `xml:"child"`
}

// SubsonicIndexes represents a Subsonic indexes container
type SubsonicIndexes struct {
	XMLName xml.Name `xml:"indexes,omitempty"`

	LastModified int64           `xml:"lastModified,attr"`
	Indexes      []SubsonicIndex `xml:"index"`
}

// SubsonicIndex represents an alphabetical Subsonic index
type SubsonicIndex struct {
	XMLName xml.Name `xml:"index"`

	Name string `xml:"name,attr"`

	Artists []SubsonicArtist `xml:"artist"`
}

// SubsonicChild is any item displayed to Subsonic when browsing using getMusicDirectory
type SubsonicChild struct {
	// Container name
	XMLName xml.Name `xml:"child,omitempty"`

	// Attributes
	ID       string `xml:"id,attr"`
	Title    string `xml:"title,attr"`
	Album    string `xml:"album,attr"`
	Artist   string `xml:"artist,attr"`
	IsDir    bool   `xml:"isDir,attr"`
	CoverArt int    `xml:"coverArt,attr"`
	Created  string `xml:"created,attr"`
}

// SubsonicPlaylists represents the Subsonic playlists container
type SubsonicPlaylists struct {
	XMLName xml.Name `xml:"playlists,omitempty"`

	Playlists []SubsonicPlaylist `xml:"playlist"`
}

type SubsonicPlaylist struct {
	XMLName xml.Name `xml:"playlist,omitempty"`

	ID        string    `xml:"id,attr"`
	Name      string    `xml:"name,attr"`
	Comment   string    `xml:"comment,attr"`
	Owner     string    `xml:"owner,attr"`
	Public    bool      `xml:"public,attr"`
	SongCount int       `xml:"songCount,attr"`
	Duration  int       `xml:"duration,attr"`
	CoverArt  string    `xml:"coverArt,attr"`
	Created   time.Time `xml:"created,attr"`

	Entry []SubsonicPlaylistEntry `xml:"entry"`
}

type SubsonicPlaylistEntry struct {
	XMLName     xml.Name  `xml:"entry,omitempty"`
	ID          string    `xml:"id,attr"`
	Parent      string    `xml:"parent,attr"`
	Title       string    `xml:"title,attr"`
	Album       string    `xml:"album,attr"`
	Artist      string    `xml:"artist,attr"`
	IsDir       string    `xml:"isDir,attr"`
	Duration    int       `xml:"duration,attr"`
	CoverArt    string    `xml:"coverArt,attr"`
	Created     time.Time `xml:"created,attr"`
	IsVideo     bool      `xml:"isVideo,attr"`
	Path        string    `xml:"path,attr"`
	BitRate     int       `xml:"bitRate,attr"`
	Suffix      string    `xml:"suffix,attr"`
	ContentType string    `xml:"contentType,attr"`
	Type        string    `xml:"type,attr"`
}

// SubsonicMusicFolders contains a list of emulated Subsonic music folders
type SubsonicMusicFolders struct {
	// Container name
	XMLName xml.Name `xml:"musicFolders,omitempty"`

	// Music folders
	MusicFolders []SubsonicMusicFolder `xml:"musicFolder"`
}

// SubsonicStarred represents a Subsonic license
type SubsonicStarred struct {
	XMLName xml.Name `xml:"starred,omitempty"`
}

func NewSubsonicResponse() *SubsonicResponse {
	return &SubsonicResponse{
		XMLNS:   SubsonicXMLNS,
		Version: SubsonicVersion,
		Status:  "ok",
	}
}

func subsonicPing(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := NewSubsonicResponse()
	XML(w, response)
}

func subsonicGetMusicFolders(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := NewSubsonicResponse()
	response.MusicFolders = &SubsonicMusicFolders{
		MusicFolders: []SubsonicMusicFolder{
			{
				ID:   1,
				Name: "Music",
			},
		},
	}
	XML(w, response)
}

func subsonicGetIndexes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := NewSubsonicResponse()

	indexes := &SubsonicIndexes{
		LastModified: time.Now().Unix(),
	}

	response.Indexes = indexes
	XML(w, response)
}

func subsonicGetPlaylists(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := NewSubsonicResponse()

	lists, err := ListLists()
	if err != nil {
		Error(w, err)
		return
	}

	var playlists []SubsonicPlaylist

	for _, list := range lists {
		playlists = append(playlists, SubsonicPlaylist{
			ID:        list.ID,
			Name:      list.Title,
			Comment:   list.Title,
			Owner:     "admin",
			Public:    false,
			SongCount: len(list.Medias),
			Duration:  int(list.TotalLength()),
		})
	}

	response.Playlists = &SubsonicPlaylists{Playlists: playlists}
	XML(w, response)
}

func subsonicGetPlaylist(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := NewSubsonicResponse()

	list, err := FindList(r.FormValue("id"))
	if err != nil {
		Error(w, err)
		return
	}

	var entries []SubsonicPlaylistEntry
	for _, media := range list.Medias {
		entries = append(entries, SubsonicPlaylistEntry{
			ID:          media.ID,
			Title:       media.Title,
			Duration:    int(media.Length),
			CoverArt:    media.ID,
			Path:        fmt.Sprintf("/soundscape/stream/%s/%s.mp3", list.ID, media.ID),
			ContentType: "audio/mp3",
			Suffix:      "mp3",
			//Path:        fmt.Sprintf("/soundscape/stream/%s/%s.m4a", list.ID, media.ID),
			//ContentType: "audio/mp4",
			//Suffix:      "m4a",
			Type: "music",
		})
	}

	response.Playlist = &SubsonicPlaylist{
		ID:        list.ID,
		Name:      list.Title,
		Comment:   list.Title,
		Owner:     "admin",
		Public:    false,
		SongCount: len(list.Medias),
		Duration:  int(list.TotalLength()),
		CoverArt:  list.ID,
		Entry:     entries,
	}
	XML(w, response)
}

func subsonicGetCoverArt(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	media, err := FindMedia(r.FormValue("id"))
	if err != nil {
		Error(w, err)
		return
	}
	size, err := strconv.Atoi(r.FormValue("size"))
	if err != nil {
		size = 640
	}

	img, err := imaging.Open(media.ImageFile())
	if err != nil {
		Error(w, err)
		return
	}

	img = imaging.Resize(img, size, 0, imaging.Lanczos)

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Vary", "Accept-Encoding")
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", 7*86400))
	if err := imaging.Encode(w, img, imaging.JPEG); err != nil {
		Error(w, err)
		return
	}
}

func subsonicGetLyrics(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := NewSubsonicResponse()
	response.Lyrics = &SubsonicLyrics{}
	XML(w, response)
}

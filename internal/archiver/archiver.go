package archiver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/rylio/ytdl"

	"go.uber.org/zap"
	// "github.com/streamlist/streamlist/internal/youtube"
)

var (
	httpUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36"
)

type job struct {
	id      string
	source  string
	context *context.Context
	cancel  *context.CancelFunc

	imagefile string
	videofile string
	audiofile string
}

// NewArchiver returns a new Archiver
func NewArchiver(datadir string, concurrency int, logger *zap.SugaredLogger) *Archiver {
	a := &Archiver{
		datadir:     datadir,
		concurrency: concurrency,
		active:      make(map[string]*job),
		failed:      make(map[string]error),
		logger:      logger,
	}
	go a.manager()
	return a
}

// Archiver object
type Archiver struct {
	mu          sync.RWMutex
	datadir     string
	concurrency int
	queue       []*job
	active      map[string]*job
	failed      map[string]error
	logger      *zap.SugaredLogger
	debug       bool
}

func (a *Archiver) setConcurrency(n int) {
	a.lock("concurrency")
	defer a.unlock("concurrency")
	a.concurrency = n
}

func (a *Archiver) getConcurrency() int {
	a.rlock("concurrency")
	defer a.runlock("concurrency")
	return a.concurrency
}

// QueuedJobs return job queue list
func (a *Archiver) QueuedJobs() []string {
	a.rlock("QueuedJobs")
	defer a.runlock("QueuedJobID")
	var ids []string
	for _, job := range a.queue {
		ids = append(ids, job.id)
	}
	sort.Strings(ids)
	return ids
}

// ActiveJobs return job active list
func (a *Archiver) ActiveJobs() []string {
	a.rlock("ActiveJobs")
	defer a.runlock("ActiveJobs")
	var ids []string
	for id := range a.active {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// InProgress return true if job id is in progress
func (a *Archiver) InProgress(id string) bool {
	for _, job := range a.QueuedJobs() {
		if job == id {
			return true
		}
	}
	for _, job := range a.ActiveJobs() {
		if job == id {
			return true
		}
	}
	return false
}

// Remove removes a job
func (a *Archiver) Remove(id string) {
	a.lock("Remove")
	defer a.unlock("Remove")
	// find job
	job, ok := a.active[id]
	if !ok {
		return
	}
	// cancel it
	cancel := *job.cancel
	cancel()
	// remove it
	delete(a.active, job.id)

	return
}

// Add adds a job
func (a *Archiver) Add(id string, source string) {
	a.lock("Add")
	defer a.unlock("Add")
	// Already running.
	if _, ok := a.active[id]; ok {
		return
	}
	// Already queued.
	if a.queued(id) {
		return
	}
	a.queue = append(a.queue, a.newJob(id, source))
}

func (a *Archiver) lock(loc string) {
	if a.debug {
		a.logger.Debugf("lock %q", loc)
	}
	a.mu.Lock()
}

func (a *Archiver) unlock(loc string) {
	if a.debug {
		a.logger.Debugf("unlock %q", loc)
	}
	a.mu.Unlock()
}

func (a *Archiver) rlock(loc string) {
	if a.debug {
		a.logger.Debugf("rlock %q", loc)
	}
	a.mu.RLock()
}

func (a *Archiver) runlock(loc string) {
	if a.debug {
		a.logger.Debugf("runlock %q", loc)
	}
	a.mu.RUnlock()
}

func (a *Archiver) newJob(id, source string) *job {
	ctx, cancel := context.WithCancel(context.Background())
	return &job{
		id:        id,
		source:    source,
		context:   &ctx,
		cancel:    &cancel,
		imagefile: filepath.Join(a.datadir, id+".jpg"),
		videofile: filepath.Join(a.datadir, id+".mp4"),
		audiofile: filepath.Join(a.datadir, id+".m4a"),
	}
}

func (a *Archiver) queued(id string) bool {
	for _, job := range a.queue {
		if job.id == id {
			return true
		}
	}
	return false
}

func (a *Archiver) manager() {
	for {
		a.lock("manager")

		if a.debug {
			a.logger.Debugf("queue: %d active: %d concurrency: %d", len(a.queue), len(a.active), a.concurrency)
		}

		if len(a.queue) > 0 && len(a.active) < a.concurrency {
			// Shift off queue.
			job := a.queue[0]
			a.queue = a.queue[1:]

			// Start archiving job.
			a.active[job.id] = job
			go a.archive(job)
		}
		a.unlock("manager")

		time.Sleep(2 * time.Second)
	}
}

func (a *Archiver) archive(job *job) {
	var failed error

	// Clean up on completion.
	defer func() {
		a.lock("archive complete")
		delete(a.active, job.id)
		if failed != nil {
			a.logger.Errorf("archive job %q failed: %s", job.id, failed)
			a.failed[job.id] = failed
		}
		a.unlock("archive complete")
	}()

	vinfo, err := ytdl.GetVideoInfoFromID(job.id)
	if err != nil {
		failed = err
		return
	}

	// image
	imgmax := fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", vinfo.ID)
	imgsd := fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", vinfo.ID)

	if maxerr := a.download(*job.context, imgmax, job.imagefile); maxerr != nil {
		if sderr := a.download(*job.context, imgsd, job.imagefile); sderr != nil {
			failed = fmt.Errorf("max: %s sd: %s", maxerr, sderr)
			return
		}
	}

	// video
	videourl, err := vinfo.GetDownloadURL(vinfo.Formats[0])
	if err != nil {
		failed = err
		return
	}

	defer os.Remove(job.videofile)

	if err := a.download(*job.context, videourl.String(), job.videofile); err != nil {
		failed = err
		return
	}

	// transcode to mp4/aac
	if err := a.transcode(*job.context, job.videofile, job.audiofile); err != nil {
		failed = err
		return
	}
}

func (a *Archiver) download(ctx context.Context, rawurl, filename string) error {
	// request file
	req, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", httpUserAgent)
	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("download %s failed: %s", rawurl, http.StatusText(res.StatusCode))
	}

	// write to file
	tmpname := filename + ".downloading"
	defer os.Remove(tmpname) // clean up on failure

	f, err := os.Create(tmpname)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, res.Body); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), filename)
}

func ffprobe(ctx context.Context, filename string) (*ffprobeInfo, error) {
	exe, err := exec.LookPath("ffprobe")
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, exe,
		"-i", filename,
		"-v", "quiet",
		"-print_format", "json",
		"-show_format", "-show_streams",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffprobe %q failed: %s\n%s", filename, err, string(output))
	}

	var ffinfo ffprobeInfo
	if err := json.Unmarshal(output, &ffinfo); err != nil {
		return nil, err
	}
	return &ffinfo, nil
}

func (a *Archiver) transcode(ctx context.Context, videoFile, audioFile string) error {
	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		return err
	}

	ffinfo, err := ffprobe(ctx, videoFile)
	if err != nil {
		return err
	}

	// Transcode to aac and delete the video.
	err = func() error {
		audioCodec := "aac"

		for i, stream := range ffinfo.Streams {
			typ := stream.CodecType
			name := stream.CodecName

			a.logger.Debugf("stream #%d %q codec is %q", i, typ, name)

			if typ == "audio" && name == "aac" {
				audioCodec = "copy"
			}
		}

		tmpname := videoFile + ".transcoding"
		defer os.Remove(tmpname)

		args := []string{
			"-y", "-i", videoFile,
			"-vn",
			"-c:a", audioCodec,
			"-strict", "experimental",
			"-movflags", "faststart",
			"-f", "mp4",
			tmpname,
		}
		a.logger.Debugf("transcoding with %s %s", ffmpeg, strings.Join(args, " "))

		output, err := exec.CommandContext(ctx, ffmpeg, args...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("transcoding %q to %q failed: %s\n%s", videoFile, tmpname, err, string(output))
		}
		return os.Rename(tmpname, audioFile)
	}()
	return err
}

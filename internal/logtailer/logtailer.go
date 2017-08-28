package logtailer

import (
	"github.com/armon/circbuf"
	"strings"
	"sync"
)

type Logtailer struct {
	sync.RWMutex

	tail *circbuf.Buffer
}

func NewLogtailer(size int64) (*Logtailer, error) {
	buf, err := circbuf.NewBuffer(size)
	if err != nil {
		return nil, err
	}
	return &Logtailer{tail: buf}, nil
}

func (l *Logtailer) Lines() []string {
	l.RLock()
	buf := l.tail.Bytes()
	l.RUnlock()

	s := string(buf)
	start := 0
	if nl := strings.Index(s, "\n"); nl != -1 {
		start = nl + len("\n")
	}
	return strings.Split(s[start:], "\n")
}

func (l *Logtailer) Write(buf []byte) (int, error) {
	l.Lock()
	n, err := l.tail.Write(buf)
	l.Unlock()
	return n, err
}

func (l *Logtailer) Sync() error {
	return nil
}

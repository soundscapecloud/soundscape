package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

type DiskInfo struct {
	free int64
	used int64
}

func (d *DiskInfo) Total() int64   { return d.free + d.used }
func (d *DiskInfo) TotalMB() int64 { return d.Total() / 1024 / 1024 }
func (d *DiskInfo) TotalGB() int64 { return d.TotalMB() / 1024 }

func (d *DiskInfo) Free() int64   { return d.free }
func (d *DiskInfo) FreeMB() int64 { return d.free / 1024 / 1024 }
func (d *DiskInfo) FreeGB() int64 { return d.FreeMB() / 1024 }

func (d *DiskInfo) Used() int64   { return d.used }
func (d *DiskInfo) UsedMB() int64 { return d.used / 1024 / 1024 }
func (d *DiskInfo) UsedGB() int64 { return d.UsedMB() / 1024 }

func (d *DiskInfo) UsedPercent() float64 {
	return (float64(d.used) / float64(d.Total())) * 100
}

func NewDiskInfo(path string) (*DiskInfo, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return nil, fmt.Errorf("diskinfo failed: %s", err)
	}
	free := stat.Bavail * uint64(stat.Bsize)
	used := (stat.Blocks * uint64(stat.Bsize)) - free
	return &DiskInfo{int64(free), int64(used)}, nil
}

func RandomNumber() (int, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return 0, err
	}
	return int(binary.LittleEndian.Uint32(b)), nil
}

func Overwrite(filename string, data []byte, perm os.FileMode) error {
	f, err := ioutil.TempFile(filepath.Dir(filename), filepath.Base(filename)+".tmp")
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Chmod(f.Name(), perm); err != nil {
		return err
	}
	return os.Rename(f.Name(), filename)
}

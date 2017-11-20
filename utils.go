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

type diskInfo struct {
	free int64
	used int64
}

func (d *diskInfo) Total() int64   { return d.free + d.used }
func (d *diskInfo) TotalMB() int64 { return d.Total() / 1024 / 1024 }
func (d *diskInfo) TotalGB() int64 { return d.TotalMB() / 1024 }

func (d *diskInfo) Free() int64   { return d.free }
func (d *diskInfo) FreeMB() int64 { return d.free / 1024 / 1024 }
func (d *diskInfo) FreeGB() int64 { return d.FreeMB() / 1024 }

func (d *diskInfo) Used() int64   { return d.used }
func (d *diskInfo) UsedMB() int64 { return d.used / 1024 / 1024 }
func (d *diskInfo) UsedGB() int64 { return d.UsedMB() / 1024 }

func (d *diskInfo) UsedPercent() float64 {
	return (float64(d.used) / float64(d.Total())) * 100
}

func newDiskInfo(path string) (*diskInfo, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return nil, fmt.Errorf("diskinfo failed: %s", err)
	}
	free := stat.Bavail * uint64(stat.Bsize)
	used := (stat.Blocks * uint64(stat.Bsize)) - free
	return &diskInfo{int64(free), int64(used)}, nil
}

func randomNumber() (int, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return 0, err
	}
	return int(binary.LittleEndian.Uint32(b)), nil
}

func overwrite(filename string, data []byte, perm os.FileMode) error {
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

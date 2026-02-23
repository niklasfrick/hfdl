package progress

import (
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	barWidth   = 40
	updateFreq = 100 * time.Millisecond
)

type ProgressBar struct {
	name       string
	total      int64
	downloaded int64
	startTime  time.Time
	lastUpdate time.Time
}

func NewProgressBar(name string, total int64) *ProgressBar {
	return &ProgressBar{
		name:       name,
		total:      total,
		startTime:  time.Now(),
		lastUpdate: time.Now(),
	}
}

func (p *ProgressBar) Start() {
	fmt.Printf("\r\033[2K")
	fmt.Printf("⏳ %s\n", p.name)
}

func (p *ProgressBar) Update() {
	now := time.Now()
	elapsed := now.Sub(p.startTime)
	speed := float64(p.downloaded) / elapsed.Seconds()

	var percent float64
	if p.total > 0 {
		percent = 100.0 * float64(p.downloaded) / float64(p.total)
	}

	var eta string
	if p.total > 0 && p.downloaded > 0 {
		remaining := p.total - p.downloaded
		if speed > 0 {
			etaSeconds := remaining / int64(speed)
			if etaSeconds < 60 {
				eta = fmt.Sprintf("%ds", etaSeconds)
			} else {
				eta = fmt.Sprintf("%dm%ds", etaSeconds/60, etaSeconds%60)
			}
		}
	}

	prog := int(percent / (100.0 / float64(barWidth)))
	bar := strings.Repeat("=", prog) + strings.Repeat(" ", barWidth-prog)

	var sizeStr, speedStr string
	if p.total > 0 {
		sizeStr = fmt.Sprintf("%.1f/%.1f MB", float64(p.downloaded)/1024/1024, float64(p.total)/1024/1024)
	} else {
		sizeStr = fmt.Sprintf("%.1f MB", float64(p.downloaded)/1024/1024)
	}
	if speed > 1024*1024 {
		speedStr = fmt.Sprintf("%.1f MB/s", speed/1024/1024)
	} else if speed > 1024 {
		speedStr = fmt.Sprintf("%.1f KB/s", speed/1024)
	} else {
		speedStr = fmt.Sprintf("%.0f B/s", speed)
	}

	etaStr := ""
	if eta != "" {
		etaStr = fmt.Sprintf(" ETA: %s", eta)
	}

	fmt.Printf("\r\033[2K[%s] %.1f%%  %s  %s%s", bar, percent, sizeStr, speedStr, etaStr)
}

func (p *ProgressBar) Increment(n int) {
	p.downloaded += int64(n)
	if time.Since(p.lastUpdate) >= updateFreq {
		p.Update()
		p.lastUpdate = time.Now()
	}
}

type ProgressReader struct {
	Reader io.Reader
	Bar    *ProgressBar
}

func (p *ProgressReader) Read(b []byte) (int, error) {
	n, err := p.Reader.Read(b)
	if n > 0 {
		p.Bar.Increment(n)
	}
	return n, err
}

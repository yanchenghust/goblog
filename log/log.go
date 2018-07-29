package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Provider interface {
	Write(ch *logChan)
	Stop()
	Start()
	Flush()
}

type SplitBy int16

const (
	SplitByHour SplitBy = 1
	SplitByDay  SplitBy = 2
)

type FileProvider struct {
	BaseProvider
	filename      string
	fp            *os.File
	buf           string
	bufLimit      int
	splitBy       SplitBy
	lastSplitName string
}

type ConsoleProvider struct {
	BaseProvider
}

func (p *ConsoleProvider) Start() {

}
func (p *ConsoleProvider) Stop() {

}
func (p *ConsoleProvider) Write(ch *logChan) {
	msg := p.format(ch)
	os.Stdout.WriteString(msg)
}
func (p *ConsoleProvider) Flush() {

}

type BaseProvider struct {
}

func NewFileProvider(filename string) *FileProvider {
	return &FileProvider{
		filename: filename,
		bufLimit: 1024,
		splitBy:  SplitByHour,
	}
}

func NewConsoleProvider() *ConsoleProvider {
	return &ConsoleProvider{}
}

func (p *FileProvider) split() {
	splitName := p.getSplitName()
	if p.lastSplitName == splitName {
		return
	}
	p.lastSplitName = splitName
	p.Flush()
	p.logFile()
}

func (p *FileProvider) getSplitName() string {
	now := time.Now()
	switch p.splitBy {
	case SplitByDay:
		return now.Format("2006-01-02")
	case SplitByHour:
		return now.Format("2006-01-02-03")
	}
	return ""
}

func (p *FileProvider) Write(ch *logChan) {
	if p.fp == nil {
		panic("fp nil")
	}
	p.split()
	msg := p.format(ch)
	p.buf += msg
	if len(msg) >= p.bufLimit {
		p.fp.WriteString(p.buf)
		p.buf = ""
	}
}

func (p *FileProvider) logFile() {
	logfile := p.filename + "." + p.lastSplitName + ".log"
	fp, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file error, err: %v, filename: %s", err, p.filename))
	}
	p.fp = fp
}

func (p *FileProvider) Flush() {
	if p.fp == nil {
		panic("fp nil")
	}
	_, err := p.fp.WriteString(p.buf)
	if err != nil {
		panic(fmt.Sprintf("WriteString error, err: %v", err))
	}
	p.buf = ""
	p.fp.Sync()
}

func (p *FileProvider) Start() {
	p.lastSplitName = p.getSplitName()
	p.logFile()
}

func (p *FileProvider) Stop() {
	if p.fp == nil {
		return
	}
	p.Flush()
	p.fp.Close()
}

func (p *BaseProvider) format(ch *logChan) string {
	now := time.Now().Format("Mon Jan 2 03:04:05.000 MST-07 2006")
	msg := fmt.Sprintf("[%s][%s]%s\n", formatLevel(ch.level), now, ch.msg)
	return msg
}

type Level int16

const (
	LvlDebug Level = 1
	LvlInfo  Level = 2
	LvlWarn  Level = 3
	LvlError Level = 4
)

var lvlStrs = map[Level]string{
	LvlDebug: "Debug",
	LvlInfo:  "Info",
	LvlWarn:  "Warn",
	LvlError: "Error",
}

type logChan struct {
	msg   string
	level Level
}
type Log struct {
	providers []Provider
	level     Level
	ch        chan *logChan
	stop      chan *sync.WaitGroup
	flush     chan *sync.WaitGroup
}

func NewLog(level Level, providers ...Provider) *Log {
	log := &Log{
		level: level,
		ch:    make(chan *logChan),
		stop:  make(chan *sync.WaitGroup),
		flush: make(chan *sync.WaitGroup),
	}
	for _, p := range providers {
		log.AppendProvider(p)
	}
	return log
}

func formatLevel(lvl Level) string {
	if s, found := lvlStrs[lvl]; found {
		return s
	}
	return "unknown"
}

func (l *Log) AppendProvider(provider Provider) {
	l.providers = append(l.providers, provider)
}

func (l *Log) Debugf(format string, args ...interface{}) {

	if LvlDebug < l.level {
		return
	}
	l.WriteLog(LvlDebug, format, args...)
}
func (l *Log) Infof(format string, args ...interface{}) {

	if LvlInfo < l.level {
		return
	}
	l.WriteLog(LvlInfo, format, args...)
}

func (l *Log) Errorf(format string, args ...interface{}) {

	if LvlError < l.level {
		return
	}
	l.WriteLog(LvlError, format, args...)
}

func (l *Log) Warnf(format string, args ...interface{}) {

	if LvlWarn < l.level {
		return
	}
	l.WriteLog(LvlWarn, format, args...)
}

func (l *Log) WriteLog(lvl Level, format string, args ...interface{}) {
	l.ch <- &logChan{
		msg:   fmt.Sprintf(format, args),
		level: lvl,
	}
}

func (l *Log) Stop() {
	l.Flush()
	var wg = sync.WaitGroup{}
	l.stop <- &wg
	wg.Wait()
}

func (l *Log) Flush() {
	var wg = sync.WaitGroup{}
	wg.Add(1)
	l.flush <- &wg
	wg.Wait()
}

func (l *Log) Start() {
	go func() {
		l.Loop()
	}()
}
func (l *Log) Loop() {
	for _, provider := range l.providers {
		provider.Start()
	}
	noStop := true
	for noStop {
		select {
		case ch := <-l.ch:
			for _, provider := range l.providers {
				provider.Write(ch)
			}
		case wg := <-l.stop:
			wg.Add(1)
			for _, provider := range l.providers {
				provider.Stop()
			}
			wg.Done()
			noStop = false
		case wg := <-l.flush:
			for _, provider := range l.providers {
				provider.Flush()
			}
			wg.Done()
		default:

		}
	}

}

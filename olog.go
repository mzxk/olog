package olog

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	//INFO 信息
	INFO = "[I]"
	//ERR 错误
	ERR = "[E]"
	//WARN 警告
	WARN = "[W]"
	//IMP 重要
	IMP = "[P]"
)

// var logs *Logg
var lg *logg
var once sync.Once

func E(v ...interface{}) {
	log.Println(ERR, v)
}
func P(v ...interface{}) {
	log.Println(IMP, v)
}
func I(v ...interface{}) {
	log.Println(INFO, v)
}
func W(v ...interface{}) {
	log.Println(WARN, v)
}
func init() {
	lg = &logg{}
	//lg.SaveFile()
	log.SetOutput(lg)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type logg struct {
	name  string
	write *os.File
}

func (t *logg) Write(w []byte) (n int, err error) {
	fmt.Print(string(w))
	once.Do(t.SaveFile)
	return t.write.Write(w)
}

func (t *logg) SaveFile() {
	t.getname()
	now := time.Now()
	filename := t.name + now.Format("2006-01-02") + ".log"
	fmt.Println("logFileName", filename)
	var logfile *os.File
	var err error
	logfile, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		//fmt.Println(err)
		logfile, err = os.Create(filename)
		if err != nil {
			fmt.Println("error save log", err)
		}
	}
	//logfile.Write([]byte(strings.Join(os.Args, "-")))
	lastlog := logfile
	t.write = logfile

	go func() {
		for {
			if time.Now().Day() != now.Day() {
				name := t.name + time.Now().Format("2006-01-02") + ".log"
				file, err := os.Create(name)
				if err != nil {
					fmt.Println("error save log", err)
					time.Sleep(1 * time.Second)
					continue
				}
				file.Write([]byte(strings.Join(os.Args, "-")))
				t.write = file
				time.Sleep(1 * time.Hour)
				lastlog.Close()
				lastlog = file
				now = time.Now()
			}
			time.Sleep(360 * time.Second)
		}
	}()
}
func (t *logg) getname() (s string) {
	for i, arg := range os.Args {
		if i == 0 {
			s = filepath.Base(arg)
			if runtime.GOOS != "windows" {
				u, e := user.Current()
				if e == nil {
					ss := strings.Split(arg, "/")
					s = u.HomeDir + "/log/" + ss[len(ss)-1]
				}
			}
		} else {
			s = s + "_" + arg
		}
	}
	if runtime.GOOS == "windows" {
		s = s + "_"
	} else {
		s = s + "/"
	}
	t.name = s
	os.MkdirAll(s, 0777)
	return

}

package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// File a file to waiting for processing
type File struct {
	ID       string
	Path     []string
	Date     time.Time
	FullPath string
}

var err = log.New(os.Stderr, "ERROR: ", 0)
var out = log.New(os.Stdout, "", 0)

func l(v ...interface{}) {
	out.Println(v...)
}

func exit(v ...interface{}) {
	err.Println(v...)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		exit("no directory is provided")
	}

	wds := os.Args[1:]
	all := []File{}
	for _, wd := range wds {
		if strings.HasSuffix(wd, string(os.PathSeparator)) {
			wd = wd[0 : len(wd)-len(string(os.PathSeparator))]
		}
		err := filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
			_, file := filepath.Split(path)
			// skip hidden directories
			if info.IsDir() && strings.HasPrefix(file, ".") {
				return filepath.SkipDir
			}
			// skip hidden files
			if !info.IsDir() && !strings.HasPrefix(file, ".") {
				rel, err := filepath.Rel(wd, path)
				if err != nil {
					exit(err.Error())
				}
				comps := strings.Split(rel, string(filepath.Separator))
				for i, v := range comps {
					comps[i] = fmt.Sprintf("%q", v)
				}
				comps = append([]string{fmt.Sprintf("%q", wd)}, comps...)
				all = append(all, File{
					ID:       fmt.Sprintf("%02x", sha1.Sum([]byte(rel))),
					FullPath: path,
					Date:     info.ModTime(),
					Path:     comps,
				})
			}
			return nil
		})
		if err != nil {
			exit(err.Error())
		}
	}

	pkg := os.Getenv("PKG")
	if len(pkg) == 0 {
		pkg = "main"
	}

	l(`// Generated by github.com/novakit/binfs`)
	l(``)
	l(`package ` + pkg)
	l(``)
	l(`import (`)
	l(`"time"`)
	l(`"github.com/novakit/binfs"`)
	l(`)`)
	l(``)
	l(`var (`)

	buf := make([]byte, 32)

	for _, f := range all {
		func() {
			// open file
			bs, err := os.Open(f.FullPath)
			if err != nil {
				exit(err.Error())
			}
			defer bs.Close()
			// write
			l(`binfs` + f.ID + ` = binfs.Chunk{`)
			l(`Path: []string{` + strings.Join(f.Path, ", ") + "},")
			l(`Date: time.Unix(` + fmt.Sprintf("%d", f.Date.Unix()) + `, 0),`)
			l(`Data: []byte{`)

			for {
				n, err := bs.Read(buf)
				if err != nil {
					if err != io.EOF {
						exit(err.Error())
					} else {
						break
					}
				} else if n > 0 {
					var line string
					for i := 0; i < n; i++ {
						line += fmt.Sprintf("0x%02x,", buf[i])
					}
					l(line)
				}
			}

			l(`},`)
			l(`}`)
		}()
	}

	l(`)`)
	l(`func init() {`)
	for _, v := range all {
		l(`binfs.Load(&binfs` + v.ID + `)`)
	}
	l(`}`)
}
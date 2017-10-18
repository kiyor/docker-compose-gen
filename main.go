/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 10-17-2017

* Last Modified : Wed 18 Oct 2017 01:30:59 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

var (
	flagPort       flagSliceString
	flagAddHost    flagSliceString
	composeVersion = flag.String("version", "3.3", "docker compose version")
)

func init() {
	flag.Var(&flagPort, "p", "docker port mount")
	flag.Var(&flagAddHost, "add-host", "add host to hosts file")
	flag.Parse()
}

func main() {
	d := data{
		Version:      *composeVersion,
		Name:         baseName(),
		Dir:          goPwd(),
		ContinerPort: getContinerPorts(flagPort),
		MountPort:    optimizeMountPort(flagPort),
		ExtraHosts:   flagAddHost,
	}
	err := write("Dockerfile", DOCKERFILE, d)
	if err != nil {
		log.Fatal(err)
	}
	err = write("docker-compose.yml", DOCKERCOMPOSE, d)
	if err != nil {
		log.Fatal(err)
	}
}

func getContinerPorts(ps flagSliceString) []string {
	var res []string
	for _, v := range ps {
		res = append(res, getContinerPort(v))
	}
	return res
}

func getContinerPort(p string) string {
	var rt string
	for _, v := range strings.Split(p, ":") {
		if _, err := strconv.ParseInt(v, 10, 64); err == nil {
			rt = v
		}
	}
	return rt
}

func optimizeMountPort(ps flagSliceString) []string {
	var res []string
	for _, v := range ps {
		if strings.Contains(v, ":") {
			res = append(res, v)
		} else {
			res = append(res, v+":"+v)
		}
	}
	return res
}

type data struct {
	Version      string
	Name         string
	Dir          string
	ContinerPort []string
	MountPort    []string
	ExtraHosts   []string
}

func write(file, tpl string, d interface{}, id ...int) error {
	filename := file
	if len(id) > 0 {
		filename = fmt.Sprintf("%s.%d", file, id[0])
	} else {
		id = append(id, 0)
	}
	if _, err := os.Stat(filename); err == nil {
		return write(file, tpl, d, id[0]+1)
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	t, err := template.New(filename).Parse(tpl)
	if err != nil {
		return err
	}
	err = t.Execute(f, d)
	if err != nil {
		return err
	}
	log.Println("write to file", filename)
	return nil
}

func goPwd() string {
	p := strings.Split(pwd(), "/src/")
	if len(p) > 0 {
		return fmt.Sprintf("/go/src/%s", p[1])
	}
	return ""
}

func baseName() string {
	return filepath.Base(pwd())
}

func pwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

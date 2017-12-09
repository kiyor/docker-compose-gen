/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 10-17-2017

* Last Modified : Sat 09 Dec 2017 01:32:43 AM UTC

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
	flagMount      flagSliceString
	flagCapAdd     flagSliceString
	flagEnv        flagSliceString
	composeVersion = flag.String("version", "3.4", "docker compose version")
	beego          = flag.Bool("beego", false, "beego app")
)

type data struct {
	Version      string
	Name         string
	Dir          string
	ContinerPort []string
	MountPort    []string
	MountDisk    []string
	ExtraHosts   []string
	CapAdd       []string
	Env          []string
	Image        string
	Command      string
}

func init() {
	flag.Var(&flagPort, "p", "docker port mount")
	flag.Var(&flagAddHost, "add-host", "add host to hosts file")
	flag.Var(&flagMount, "v", "docker volume mount")
	flag.Var(&flagCapAdd, "cap_add", "docker cap_add")
	flag.Var(&flagEnv, "e", "docker env")
	flag.Parse()
}

func main() {
	var image, command string
	if len(flag.Args()) > 0 {
		image = flag.Args()[0]
		if len(flag.Args()) > 1 {
			command = strings.Join(flag.Args()[1:], " ")
		}
	}
	d := data{
		Version:      *composeVersion,
		Name:         baseName(),
		Dir:          goPwd(),
		ContinerPort: getContinerPorts(flagPort),
		MountPort:    optimizeMountPort(flagPort),
		MountDisk:    getContinerMounts(flagMount),
		ExtraHosts:   flagAddHost,
		CapAdd:       flagCapAdd,
		Env:          flagEnv,
		Image:        image,
		Command:      command,
	}
	dockerfile := DOCKERFILE
	if *beego || isBeego() {
		dockerfile = DOCKERFILE_BEE
	}
	if len(image) == 0 {
		err := write("Dockerfile", dockerfile, d)
		if err != nil {
			log.Fatal(err)
		}
	}
	err := write("docker-compose.yml", DOCKERCOMPOSE, d)
	if err != nil {
		log.Fatal(err)
	}
}

func isBeego() bool {
	dirs := []string{"conf", "controllers", "models", "routers", "static", "views"}
	for _, name := range dirs {
		if _, err := os.Stat(name); err != nil {
			return false
		}
	}
	return true
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
func getContinerMounts(ps flagSliceString) []string {
	var res []string
	for _, v := range ps {
		res = append(res, getContinerMount(v))
	}
	return res
}
func getContinerMount(p string) string {
	if strings.Contains(p, ":") {
		return p
	}
	return p + ":" + p
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
	t, err := template.New(filename).Funcs(template.FuncMap{
		"split": split,
	}).Parse(tpl)
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
	if len(p) > 1 {
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

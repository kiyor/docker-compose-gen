/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : tmpl.go

* Purpose :

* Creation Date : 10-17-2017

* Last Modified : Wed 18 Oct 2017 01:29:56 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import ()

const (
	DOCKERFILE = `FROM golang
ADD . {{.Dir}} 
RUN cd {{.Dir}} && \
    go get && \
    go install
{{if .ContinerPort}}EXPOSE {{range .ContinerPort}}{{.}} {{end}}{{end}}
WORKDIR {{.Dir}}
ENTRYPOINT ["{{.Name}}"]
`
	DOCKERCOMPOSE = `version: '{{.Version}}'
services:
  {{.Name}}:
    container_name: {{.Name}}
    build:
      context: .
      dockerfile: Dockerfile{{if .MountPort}}
    ports:{{range .MountPort}}
      - {{.}}{{end}}{{end}}
    restart: always{{if .ExtraHosts}}
    extra_hosts:{{range .ExtraHosts}}
      - "{{.}}"{{end}}{{else}}
    extra_hosts:
      - "node:172.17.0.1"{{end}}
`
)

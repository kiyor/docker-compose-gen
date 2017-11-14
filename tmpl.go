/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : tmpl.go

* Purpose :

* Creation Date : 10-17-2017

* Last Modified : Mon 30 Oct 2017 10:33:33 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"strings"
)

const (
	DOCKERFILE = `FROM golang as builder
WORKDIR {{.Dir}}
COPY . .
RUN go get && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o {{.Name}} .

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY . .
COPY --from=builder {{.Dir}}/{{.Name}} .
{{if .ContinerPort}}EXPOSE {{range .ContinerPort}}{{.}} {{end}}{{end}}
{{if .MountDisk}}VOLUME [{{range $k, $v := .MountDisk}}"{{if $k}},{{end}}{{index ( split $v ":" ) 1}}"{{end}}]{{end}}
ENTRYPOINT ["/root/{{.Name}}"]
`
	DOCKERFILE_BEE = `FROM golang as builder
WORKDIR {{.Dir}}
COPY controllers controllers
COPY models models
COPY routers routers
COPY *.go ./
RUN go get && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o {{.Name}} .

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY conf conf
COPY static static
COPY views views
COPY --from=builder {{.Dir}}/{{.Name}} .
{{if .ContinerPort}}EXPOSE {{range .ContinerPort}}{{.}} {{end}}{{end}}
{{if .MountDisk}}VOLUME [{{range $k, $v := .MountDisk}}"{{if $k}},{{end}}{{index ( split $v ":" ) 1}}"{{end}}]{{end}}
ENTRYPOINT ["/root/{{.Name}}"]
`
	DOCKERCOMPOSE = `version: '{{.Version}}'
services:
  {{.Name}}:
    container_name: {{.Name}}{{if .Image}}
    image: {{.Image}}{{else}}
    build:
      context: .
      dockerfile: Dockerfile{{end}}{{if .Command}}
    command: {{.Command}}{{end}}{{if .MountPort}}
    ports:{{range .MountPort}}
      - "{{.}}"{{end}}{{end}}{{if .MountDisk}}
    volumes:{{range .MountDisk}}
      - "{{.}}"{{end}}{{end}}{{if .CapAdd}}
    cap_add:{{range .CapAdd}}
      - "{{.}}"{{end}}{{end}}
    restart: always{{if .ExtraHosts}}
    extra_hosts:{{range .ExtraHosts}}
      - "{{.}}"{{end}}{{else}}
    extra_hosts:
      - "node:172.17.0.1"{{end}}
`
)

func split(in, sep string) []string {
	return strings.Split(in, sep)
}

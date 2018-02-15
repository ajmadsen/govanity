package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/ajmadsen/govanity/server"
)

var (
	addr       string
	redirs     stringSlice
	serverName string
)

type stringSlice []string

func (ss *stringSlice) String() string {
	if ss == nil {
		return ""
	}
	return strings.Join(*ss, ",")
}

func (ss *stringSlice) Set(f string) error {
	*ss = append(*ss, f)
	return nil
}

func (ss *stringSlice) Get() interface{} {
	return *ss
}

func init() {
	flag.StringVar(&addr, "addr", ":http", "address to listen on")
	flag.Var(&redirs, "redir", "list of repo-name:canonical-path to redirect for")
	flag.StringVar(&serverName, "server-name", "", "the hostname of the server that will comprise the prefix of the repository")
}

func main() {
	flag.Parse()

	if len(redirs) == 0 {
		log.Fatalln("please add one or more redirections")
	}
	if serverName == "" {
		log.Fatalln("please supply a server name")
	}

	backend := make(server.TextBackend)
	for _, r := range redirs {
		if err := backend.Add(r); err != nil {
			log.Fatalf("could not add redirection %q: %v", r, err)
		}
	}

	s := &server.Server{
		Backend:    backend,
		ServerName: serverName,
	}

	log.Println("starting up")
	log.Fatalln(http.ListenAndServe(addr, s))
}

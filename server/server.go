package server

import (
	"log"
	"net/http"
	"path"
	"strings"

	"html/template"
)

var importPathTmpl = template.Must(template.New("importPath").Parse(`<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Moved</title>
		<meta name="go-import" content="{{.Prefix}} git {{.Canonical}}">
		<meta http-equiv="refresh" content="0;URL='{{.Canonical}}'">
	</head>
	<body>
		<p>Moved <a href="{{.Canonical}}">here</a>.</p>
		<script>window.location={{.Canonical}};</script>
	</body>
</html>
`))

type tmplArgs struct {
	Prefix    string
	Canonical string
}

type Server struct {
	Backend
	ServerName string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL

	rpath := path.Clean(url.Path)
	if strings.HasPrefix(rpath, "/") {
		rpath = rpath[1:]
	}

	pathParts := strings.SplitN(rpath, "/", 2)
	repo := pathParts[0]

	log.Printf("got request for repo %q", repo)

	isRepo, err := s.IsRepo(repo)
	if err != nil || !isRepo {
		http.NotFound(w, r)
		return
	}

	canonical := s.Canonical(repo)

	if url.Query().Get("go-get") != "1" {
		http.Redirect(w, r, canonical, http.StatusFound)
		return
	}

	prefix := strings.TrimPrefix(path.Join(s.ServerName, repo), "/")

	err = importPathTmpl.Execute(w, &tmplArgs{
		Canonical: canonical,
		Prefix:    prefix,
	})
	if err != nil {
		log.Printf("failed to execute template: %v", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

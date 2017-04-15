package views

import _solr "github.com/rtt/Go-Solr"

const (
	HOST = "localhost"
	PORT = 8983
	CORE = "newcore"
)

var solr, s_err = _solr.Init(HOST, PORT, CORE)

func init() {
	if s_err != nil {
		panic("Unable to init solr")
	}
}

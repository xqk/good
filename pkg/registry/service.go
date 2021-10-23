package registry

import "github.com/xqk/good/pkg/server"

type service struct {
	server.Server
}

func Service(srv server.Server) server.Server {
	return &service{Server: srv}
}

func (s *service) Serve() error {

	return s.Server.Serve()
}

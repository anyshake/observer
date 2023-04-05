package server

type ServerOptions struct {
	Listen string
	Port   string
	Cors   bool
	Gzip   int
}

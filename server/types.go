package server

type ServerOptions struct {
	Version string
	Host    string
	Port    int
	Cors    bool
	Gzip    int
}

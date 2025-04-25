package server

import (
	"quotobot/internal/server"
)

func main() {
	s := server.NewServer()
	s.Start()
}

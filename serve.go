package redis

import (
	"crypto/tls"
	"log"
	"time"

	"github.com/tidwall/redcon"
)

// Run runs the default redis server.
// Initializes the default redis if not already.
func Run(addr string) error {
	return Default().Run(addr)
}

// Run runs the redis server.
func (r *Redis) Run(addr string) error {
	go r.keyExpirer.Start(1 * time.Second)
	return redcon.ListenAndServe(
		addr,
		func(conn redcon.Conn, cmd redcon.Command) {
			addr := conn.RemoteAddr()
			c := r.clients[addr]

			if c == nil {
				log.Printf("Attempting to handle a command from a disconnected client with address [%s]\n", addr)
				return
			}

			r.handler(c, cmd)
		},
		func(conn redcon.Conn) bool {
			addr := conn.RemoteAddr()

			if r.clients[addr] != nil {
				log.Printf("Client from this address [%s] already exists\n", addr)
				return false
			}
			r.clients[addr] = r.NewClient(conn)
			return true

		},
		func(conn redcon.Conn, err error) {
			addr := conn.RemoteAddr()
			delete(r.clients, addr)
		},
	)
}

// Run runs the redis server with tls.
func (r *Redis) RunTLS(addr string, tls *tls.Config) error {
	go r.keyExpirer.Start(1 * time.Second)
	return redcon.ListenAndServeTLS(
		addr,
		func(conn redcon.Conn, cmd redcon.Command) {
			addr := conn.RemoteAddr()
			c := r.clients[addr]

			if c == nil {
				log.Printf("Attempting to handle a command from a disconnected client with address [%s]\n", addr)
				return
			}

			r.handler(c, cmd)
		},
		func(conn redcon.Conn) bool {
			addr := conn.RemoteAddr()

			if r.clients[addr] != nil {
				log.Printf("Client from this address [%s] already exists\n", addr)
				return false
			}
			r.clients[addr] = r.NewClient(conn)
			return true

		},
		func(conn redcon.Conn, err error) {
			addr := conn.RemoteAddr()
			delete(r.clients, addr)
		},
		tls,
	)
}

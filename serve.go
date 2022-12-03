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

			r.mu.RLock()

			c := r.clients[addr]

			r.mu.RUnlock()

			if c == nil {
				log.Printf("Attempting to handle a command from a disconnected client with address [%s]\n", addr)
				return
			}

			// TODO: There's a small race condition here where we unlocked the lock
			// above but then the connection is closed. That means that the connection is closed
			// but the handler will still use it. What should happen instead is that the
			// onClose function check if the client is still active then pause its deletion for later.
			// At the same time, this handler also checks if the client is pending for deletion.

			r.handler(c, cmd)
		},
		func(conn redcon.Conn) bool {
			addr := conn.RemoteAddr()

			if r.clients[addr] != nil {
				log.Printf("Client from this address [%s] already exists\n", addr)
				return false
			}

			r.mu.Lock()

			r.clients[addr] = r.NewClient(conn)

			r.mu.Unlock()

			return true

		},
		func(conn redcon.Conn, err error) {
			addr := conn.RemoteAddr()

			r.mu.Lock()

			delete(r.clients, addr)

			r.mu.Unlock()

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

package snippet

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
)

func reverse_proxy() {
	proxy := func(w http.ResponseWriter, req *http.Request) {
		log.Println("proxy", req.Method, req.RequestURI)
		if req.URL.Host != "" {
			if req.Method == http.MethodConnect {
				// tunnel
				conn, err := net.Dial("tcp", req.URL.Host)
				if err != nil {
					w.WriteHeader(502)
					fmt.Fprint(w, err)
					return
				}

				client, _, err := w.(http.Hijacker).Hijack()
				if err != nil {
					w.WriteHeader(502)
					fmt.Fprint(w, err)
					conn.Close()
					return
				}
				client.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))

				hr, hw := io.Pipe()
				go func() {
					io.Copy(os.Stdout, hr)
					hr.Close()
				}()
				go func() {
					// print response to stdout
					io.Copy(io.MultiWriter(client, hw), conn)
					client.Close()
					conn.Close()
					hw.Close()
				}()
				go func() {
					io.Copy(conn, client)
					client.Close()
					conn.Close()
				}()
				return
			}

			httputil.NewSingleHostReverseProxy(req.URL).ServeHTTP(w, req)
		}
	}
	http.ListenAndServe(":8021", http.HandlerFunc(proxy))
}

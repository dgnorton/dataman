package main

import (
        "fmt"
        "net"
        "time"
)

func main() {
        addr := new(net.TCPAddr)
        addr.Port = 2300
        srv, err := net.ListenTCP("tcp", addr)
        if err != nil {
                fmt.Println("error listening on port 23")
                fmt.Println(err)
                return
        }
        for {
                fmt.Println("waiting for connection...")
                conn , err := srv.AcceptTCP()
                if err != nil {
                        fmt.Println("error accepting connection")
                        continue
                }
                fmt.Println("connection accepted")
                //connected := make(chan bool)
                session(conn)
                //<-connected
        }
}

func session(conn *net.TCPConn) {
        fmt.Println("here");
        var buf[2048]byte
        code := 0
        for {
	        t := time.Now().Add(time.Millisecond*100)
                conn.SetReadDeadline(t)
                n, err := conn.Read(buf[:])
                e, ok := err.(net.Error)

                if err != nil && ok && !e.Timeout() {
                        fmt.Println(err)
                        break
                }

                if n > 0 {
                        process(conn, buf[:n])
                } else {
                        msg := fmt.Sprintf("%v", code)
                        code++
                        conn.Write([]byte(msg))
                }
        }
        fmt.Println("session ended")
}

func process(conn *net.TCPConn, buf []byte) {
        fmt.Println(string(buf))
}

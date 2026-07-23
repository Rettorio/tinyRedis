package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)


func main() {
    db := CreateAndSeed()
    ln,err := net.Listen("tcp", ":6739")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to start server : %v", err.Error())
    }
    fmt.Fprint(os.Stderr, "Server started 0.0.0.0:6739")

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed to accept connection %s", err.Error())
            continue
        }

        go handleConnection(conn,db)
    }
}


// handleConnection takes conn as input connection
// scan input text for Redis command then execute them directly to db
func handleConnection(conn net.Conn, db *rdb) {
    defer conn.Close()

    // live scanner
    scanner := bufio.NewScanner(conn)


    for scanner.Scan() {
        if scanner.Err() != nil {
            fmt.Fprintf(os.Stderr, "Error in input %v", scanner.Err().Error())
            break
        }
        input := scanner.Text()
        fmt.Fprintf(os.Stderr, "Got command: %s\n", input)

        //  split input by space
        parts := strings.Fields(input)
        if len(parts) == 0 {
            continue
        }

        command := strings.ToUpper(parts[0])
        //parts[0] is the command
        //parts[1] is key and parts[2] is the value
        switch command {
            case "SET":
                db.SET(parts[1], parts[2])
                conn.Write([]byte("+Ok\r\n"))
            case "GET":
                val := db.GET(parts[1])
                if val == "" {
                    val = "novalue"
                }
                conn.Write(fmt.Appendf(nil,"+%s\r\n", val))
                // conn.Write(fmt.Appendf("+%s\r\n", val))
            case "DEL":
                success := db.DEL(parts[1])
                if !success {
                    conn.Write(fmt.Appendf(nil, "+failed to delete %s\r\n", parts[1]))
                } else {
                    conn.Write(fmt.Appendf(nil, "+success to delete %s\r\n", parts[1]))
                }
            default :
                conn.Write([]byte("-ERR unknown command\r\n"))
        }
    }
}

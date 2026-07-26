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

    reader := bufio.NewReader(conn)
    for {
        cmd, err := ReadCommand(reader)
        if err != nil {
            break
        }

        switch strings.ToUpper(cmd.Method) {
        case "PING":
            WriteSimpleString(conn, "PONG")
        case "SET":
            if len(cmd.Args) < 2 { WriteError(conn,"Missing Value in COMMAND"); continue }
            db.SET(cmd.Args[0], cmd.Args[1])
            WriteSimpleString(conn, "OK")
        case "GET":
            if len(cmd.Args) < 1 { WriteError(conn,"Missing Key in Command"); continue }
            val := db.GET(cmd.Args[0])
            if val == "" {
                WriteNullBulkString(conn)
            } else {
                WriteBulkString(conn, val)  // $N\r\n<val>\r\n
            }
        case "COMMAND":
            WriteArrayLength(conn, 0)
        case "EXISTS":
            if len(cmd.Args) < 1 { WriteError(conn,"Missing Key in Command"); continue }
            exists := db.EXISTS(cmd.Args[0])
            WriteInteger(conn, map[bool]int{true: 1, false: 0}[exists])
        case "DEL":
            if len(cmd.Args) < 1 { WriteError(conn,"Missing Key in Command"); continue }
            ok := db.DEL(cmd.Args[0])
            WriteInteger(conn, map[bool]int{true: 1, false: 0}[ok])
        default:
            WriteError(conn, "unknown command '" + cmd.Method + "'")
        }
    }
}

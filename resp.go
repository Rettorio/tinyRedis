package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
)


type Command struct {
    Method      string
    Args        []string
}


func ReadCommand(reader *bufio.Reader) (*Command,error) {
    inpByte, err := reader.ReadByte()
    if err != nil {

        return nil,fmt.Errorf("Unkonwn error %v", err)
    }
    if inpByte != '*' {
        return nil,errors.New("Unknown Command")
    }

    slength, err := reader.ReadBytes('\n')
    if err != nil {
        return nil,fmt.Errorf("Malformed data : %v", err)
    }
    nlength, err := parseNumber(slength)
    if nlength <= -1 {
        return nil, errors.New("Malformed data")
    }
    if nlength == 0 {
        return nil, errors.New("Empty Command")
    }
    // Parse command first
    comm,err := parseBulkString(reader)
    if err != nil {
        return nil,err
    }
    var command Command
    if comm != "GET" && comm != "SET" && comm != "DEL" && comm != "PING" {
        return nil,errors.New("Unknown Command")
    }
    for i := 1; i < nlength; i++ {
        trg, err := parseBulkString(reader)
        if err != nil {
            return  nil,fmt.Errorf  ("Failed to parse command arguments : %v", err)
        }
        command.Args = append(command.Args, trg)
    }

    if comm == "SET" && len(command.Args) < 2 {
        return nil,errors.New("Missing Value in SET Command.")
    }
    if (comm == "GET" || comm == "DEL") && len(command.Args) == 0 {
        return nil,errors.New("Missing Key in Command")
    }
    command.Method = comm

    return &command,nil
}


// parseNumber parse integer from strint
// return error from strconv.Atoi()
func parseNumber(nums []byte) (int,error) {
    // eliminate \r\n
    l := len(nums) - 2
    num, err := strconv.Atoi(string(nums[0:l]));
    if err != nil {
        return 0,err
    }
    return num,nil
}

// parseBulkString takes bufio.Reader and
// parse string in resp BulkString format
func parseBulkString(reader *bufio.Reader) (string,error) {
    // Discard '$' bulk string notation
    reader.Discard(1)
    slength, err := reader.ReadBytes('\n')
    if err != nil {
        return "",fmt.Errorf("Malformed data : %v", err)
    }
    nlength, err := parseNumber(slength)
    if err != nil {
        return "",fmt.Errorf("Malformed data : %v", err)
    }
    if nlength <= 0 {
        return "",errors.New("Unknown input.")
    }

    s := make([]byte, nlength)
    _,err = io.ReadFull(reader, s)
    if err != nil {
        return "",fmt.Errorf("Something went wrong : %v", err)
    }
    reader.Discard(2)
    return string(s),nil
}


// WriteSimpleString takes net.Conn and string to write in
// resp format of simple string
// return nil in successful operation
func WriteSimpleString(conn net.Conn, content string) error {
    final := []byte{'+'}
    final = append(final, []byte(content)...)
    final = append(final, '\r', '\n')

    _,err := conn.Write(final)
    if err != nil {
        return err
    }

    return nil
}

// WriteBulkString takes net.Conn and string to write in
// resp format of bulk string
// return nil in succesful operation
func WriteBulkString(conn net.Conn, content string) error {
    final := []byte{'$'}
    final = append(final, []byte(strconv.Itoa(len(content)))...)
    final = append(final, '\r','\n')
    if content != "" {
        final = append(final, []byte(content)...)
    }
    final = append(final, '\r','\n')

    _,err := conn.Write(final)
    if err != nil {
        return err
    }

    return nil
}

// WriteNullBulkString write null string
// in resp bulkstring format
func WriteNullBulkString(conn net.Conn) error {
    _,err := conn.Write([]byte{'$','-','1','\r','\n'})
    if err != nil {
        return err
    }
    return nil
}

// WriteInteger write integer
// in resp format
// return nil in succesful operation
func WriteInteger(conn net.Conn, number int) error {
    final := []byte{':'}
    final = append(final, []byte(strconv.Itoa(number))...)

    _,err := conn.Write(append(final, '\r','\n'))
    if err != nil {
        return err
    }

    return nil
}

// WriteArrayLength writes RESP array header
func WriteArrayLength(conn net.Conn, length int) error {
	final := []byte{'*'}
	final = append(final, []byte(strconv.Itoa(length))...)
	_, err := conn.Write(append(final, '\r', '\n'))
	return err
}

// WriteError takes net.Conn and string of error
// in resp error format
// return nil if successful
func WriteError(conn net.Conn, message string) error {
    final := []byte{'-'}
    final = append(final, []byte(message)...)

    _,err := conn.Write(append(final, '\r','\n'))
    if err != nil {
        return err
    }

    return nil
}

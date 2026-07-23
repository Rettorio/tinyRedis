package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

// Session memory of reader
// The goal is to write resp reader that is aware
type ReaderSession struct {
    inputType      byte
    isComplete     bool
    inputLength    int
    content        any
    reader         *bufio.Reader
}

func NewSession(reader *bufio.Reader) (*ReaderSession,error) {
    inpByte, err := reader.ReadByte()
    if
        err != nil || (inpByte != '+' && inpByte != ':' && inpByte != '$' && inpByte != '-' && inpByte != '*') {

        return nil,errors.New("Malformed data : unknown type")
    }
    return &ReaderSession{inputType: inpByte, reader: reader},nil
}

func (session *ReaderSession) Read() error {
    // parse simple string
    if session.inputType == '+' {
        simpString, err := session.reader.ReadBytes('\n')
        session.content = string(simpString)
        session.isComplete = true
        if err != nil {
            return  err
        }
        session.reader.Discard(session.reader.Buffered())
        return  nil
    }
    // parse error
    if session.inputType == '-' {
        errMsg, err := session.reader.ReadBytes('\n')
        session.content = errors.New(string(errMsg))
        session.isComplete = true
        if err != nil {
            return  err
        }
        session.reader.Discard(session.reader.Buffered())
        return nil
    }
    // parse integer
    if session.inputType == ':' {
        numBytes, err := session.reader.ReadBytes('\n')
        if err != nil {
            return  err
        }
        num,err := parseNumber(numBytes)
        session.isComplete = true
        if err != nil {
            return err
        }
        session.content = num
        session.reader.Discard(session.reader.Buffered())
        return nil
    }
    // parse bulkstring
    if session.inputType == '$' {
        // if session.inputLength == 0 {
        slength, err := session.reader.ReadBytes('\n')
        if err != nil {
            return fmt.Errorf("Malformed data : %v", err)
        }
        // minlength of []byte is 3
        if len(slength) < 3 {
            return fmt.Errorf("Malformed data")
        }
        nlength, err := parseNumber(slength)
        if err != nil {
            return fmt.Errorf("Malformed data : %v", err)
        }
        if nlength == -1 {
            session.content = ""
            session.isComplete = true
            session.reader.Discard(session.reader.Buffered())
            return  nil
        }
        session.inputLength = nlength
        // }

        s := make([]byte, session.inputLength)
        _,err = io.ReadFull(session.reader, s)
        if err != nil {
            return err
        }
        session.content = string(s)
        session.isComplete = true
        session.reader.Discard(session.reader.Buffered())
        return nil
    }
    // parse array
    if session.inputType == '*' {
        slength, err := session.reader.ReadBytes('\n')
        if err != nil {
            session.isComplete = true
            session.reader.Discard(session.reader.Buffered())
            return fmt.Errorf("Malformed data : %v", err)
        }
        nlength, err := parseNumber(slength)
        var content []string
        // empty array
        if nlength <= 0 {
            session.content = content
            session.isComplete = true
            session.reader.Discard(session.reader.Buffered())
            return nil
        }
        // temp item length
        ntemp := 0
        for i := 1; i <= (nlength*2);i++ {
            if i % 2 == 1 {
                nbyte,err := session.reader.ReadBytes('\n')
                if err != nil {
                    session.isComplete = true
                    session.reader.Discard(session.reader.Buffered())
                    return fmt.Errorf("Malformed data : %v", err)
                }
                // {'$','4','\r','\n'} skip type notation
                ntemp,err = parseNumber(nbyte[1:])
                if err != nil {
                    session.isComplete = true
                    session.reader.Discard(session.reader.Buffered())
                    return fmt.Errorf("Malformed data : %v", err)
                }
                continue
            }
            if i % 2 == 0 {
                line := make([]byte, ntemp)
                _,err := io.ReadFull(session.reader, line)
                if err != nil {
                    session.isComplete = true
                    session.content = content
                    session.reader.Discard(session.reader.Buffered())
                    return errors.New("Malformed data : %v")
                }
                content = append(content, string(line))
                // Discard \r\n
                session.reader.Discard(2)
            }
        }
        // safely assume content is match with arr length
        session.content = content
        session.isComplete = true
        session.reader.Discard(session.reader.Buffered())
        return nil
    }
    // fallback
    return nil
}


func parseInputLength(reader *bufio.Reader) (int,error) {
    slength, err := reader.ReadBytes('\n')
    if err != nil {
        return 0,fmt.Errorf("unable to parse input length : %v", err)
    }
    // minlength of []byte is 3
    if len(slength) < 3 {
        return 0,fmt.Errorf("unable to parse input length")
    }
    nlength, err := parseNumber(slength)
    if err != nil {
        return 0,fmt.Errorf("unable to parse input length : %v", err)
    }
    if nlength == -1 {
        return 0,errors.New("unable to parse input length")
    }
    return  nlength,nil
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

// parseBulkString takes string bytes
// and length of string then return the trimmed string
// if the actual string is less than given length return halformed string and an error.
func parseBulkString(words []byte, length int) (string,error) {
    s := make([]byte, length)
    for i := 0;i < length;{
        if words[i] == '\r' || words[i] == '\n' {
            continue
        }
        s = append(s, words[i])
        i++
    }
    if len(s) != length {
        return string(s),errors.New("given string length is less than given length")
    }
    return string(s),nil
}

func parseString(words []byte) string {
    if len(words) == 0 {
        return ""
    }
    return string(words)
}

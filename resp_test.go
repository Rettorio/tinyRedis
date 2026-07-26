package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"
)

type writeRecorder struct {
	bytes.Buffer
}

func (w *writeRecorder) Close() error                       { return nil }
func (w *writeRecorder) Read(b []byte) (int, error)         { return 0, io.EOF }
func (w *writeRecorder) LocalAddr() net.Addr                { return nil }
func (w *writeRecorder) RemoteAddr() net.Addr               { return nil }
func (w *writeRecorder) SetDeadline(t time.Time) error      { return nil }
func (w *writeRecorder) SetReadDeadline(t time.Time) error  { return nil }
func (w *writeRecorder) SetWriteDeadline(t time.Time) error { return nil }

func TestReadCommand(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Command
		wantErr string
	}{
		{
			name:  "ping",
			input: "*1\r\n$4\r\nPING\r\n",
			want:  &Command{Method: "PING"},
		},
		{
			name:  "get key",
			input: "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n",
			want:  &Command{Method: "GET", Args: []string{"key"}},
		},
		{
			name:  "set key value",
			input: "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n",
			want:  &Command{Method: "SET", Args: []string{"key", "value"}},
		},
		{
			name:  "del key with space",
			input: "*2\r\n$3\r\nDEL\r\n$4\r\nHP 1\r\n",
			want:  &Command{Method: "DEL", Args: []string{"HP 1"}},
		},
		{
			name:  "set key with spaces",
			input: "*3\r\n$3\r\nSET\r\n$4\r\nHP 1\r\n$40\r\nHarry Potter and the Philosopher's Stone\r\n",
			want:  &Command{Method: "SET", Args: []string{"HP 1", "Harry Potter and the Philosopher's Stone"}},
		},
		{
			name:    "empty array",
			input:   "*0\r\n",
			wantErr: "Empty Command",
		},
		{
			name:    "not array simple string",
			input:   "+OK\r\n",
			wantErr: "Unknown Command",
		},
		{
			name:    "not array bulk string",
			input:   "$5\r\nhello\r\n",
			wantErr: "Unknown Command",
		},
		{
			name:    "inline text",
			input:   "PING\r\n",
			wantErr: "Unknown Command",
		},
		{
			name:    "unknown command",
			input:   "*1\r\n$4\r\nFOOB\r\n",
			wantErr: "Unknown Command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			got, err := ReadCommand(reader)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("error mismatch: got %q, want %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got  %+v\nwant %+v", got, tt.want)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	t.Run("simple string", func(t *testing.T) {
		var buf writeRecorder
		err := WriteSimpleString(&buf, "OK")
		if err != nil {
			t.Fatal(err)
		}
		want := "+OK\r\n"
		if buf.String() != want {
			t.Fatalf("got %q, want %q", buf.String(), want)
		}
	})

	t.Run("bulk string", func(t *testing.T) {
		var buf writeRecorder
		err := WriteBulkString(&buf, "hello")
		if err != nil {
			t.Fatal(err)
		}
		want := "$5\r\nhello\r\n"
		if buf.String() != want {
			t.Fatalf("got %q, want %q", buf.String(), want)
		}
	})

	t.Run("bulk string empty", func(t *testing.T) {
		var buf writeRecorder
		err := WriteBulkString(&buf, "")
		if err != nil {
			t.Fatal(err)
		}
		want := "$0\r\n\r\n"
		if buf.String() != want {
			t.Fatalf("got %q, want %q", buf.String(), want)
		}
	})

	t.Run("null bulk string", func(t *testing.T) {
		var buf writeRecorder
		err := WriteNullBulkString(&buf)
		if err != nil {
			t.Fatal(err)
		}
		want := "$-1\r\n"
		if buf.String() != want {
			t.Fatalf("got %q, want %q", buf.String(), want)
		}
	})

	t.Run("integer", func(t *testing.T) {
		var buf writeRecorder
		err := WriteInteger(&buf, 1)
		if err != nil {
			t.Fatal(err)
		}
		want := ":1\r\n"
		if buf.String() != want {
			t.Fatalf("got %q, want %q", buf.String(), want)
		}
	})

	t.Run("integer zero", func(t *testing.T) {
		var buf writeRecorder
		err := WriteInteger(&buf, 0)
		if err != nil {
			t.Fatal(err)
		}
		want := ":0\r\n"
		if buf.String() != want {
			t.Fatalf("got %q, want %q", buf.String(), want)
		}
	})

	t.Run("error", func(t *testing.T) {
		var buf writeRecorder
		err := WriteError(&buf, "ERR unknown command 'FOOB'")
		if err != nil {
			t.Fatal(err)
		}
		want := "-ERR unknown command 'FOOB'\r\n"
		if buf.String() != want {
			t.Fatalf("got %q, want %q", buf.String(), want)
		}
	})
}

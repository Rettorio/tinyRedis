package main

import (
	"bufio"
	"errors"
	"strings"
	"testing"
)

func TestRespRead(t *testing.T) {
	testCase := map[string]struct {
		label          string
		input          string
		want           any
		wantSessionErr string
		wantReadErr    string
	}{
		// Simple String (+)
		"simple_string_ok": {
			label: "simple OK",
			input: "+OK\r\n",
			want:  "OK\r\n",
		},
		"simple_string_hello": {
			label: "simple hello",
			input: "+hello\r\n",
			want:  "hello\r\n",
		},
		"simple_string_empty": {
			label: "simple empty",
			input: "+\r\n",
			want:  "\r\n",
		},

		// Error (-)
		"error_unknown_command": {
			label: "ERR unknown command",
			input: "-ERR unknown command\r\n",
			want:  errors.New("ERR unknown command\r\n"),
		},
		"error_wrongtype": {
			label: "WRONGTYPE",
			input: "-WRONGTYPE Operation against a key holding the wrong kind of value\r\n",
			want:  errors.New("WRONGTYPE Operation against a key holding the wrong kind of value\r\n"),
		},

		// Integer (:)
		"integer_one": {
			label: ":1",
			input: ":1\r\n",
			want:  1,
		},
		"integer_zero": {
			label: ":0",
			input: ":0\r\n",
			want:  0,
		},
		"integer_negative": {
			label: ":-1",
			input: ":-1\r\n",
			want:  -1,
		},
		"integer_large": {
			label: ":1000",
			input: ":1000\r\n",
			want:  1000,
		},

		// Bulk String ($)
		"bulk_string_hello": {
			label: "$5 hello",
			input: "$5\r\nhello\r\n",
			want:  "hello",
		},
		"bulk_string_spaces": {
			label: "$11 hello world",
			input: "$11\r\nhello world\r\n",
			want:  "hello world",
		},
		"bulk_string_empty": {
			label: "$0 empty",
			input: "$0\r\n\r\n",
			want:  "",
		},
		"bulk_string_null": {
			label: "$-1 null",
			input: "$-1\r\n",
			want:  "",
		},

		// Array (*)
		"array_empty": {
			label: "*0 empty",
			input: "*0\r\n",
			want:  []string{},
		},
		"array_one_ping": {
			label: "*1 PING",
			input: "*1\r\n$4\r\nPING\r\n",
			want:  []string{"PING"},
		},
		"array_two_get_key": {
			label: "*2 GET key",
			input: "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n",
			want:  []string{"GET", "key"},
		},
		"array_three_set_key_value": {
			label: "*3 SET key value",
			input: "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n",
			want:  []string{"SET", "key", "value"},
		},

		// Error / Edge Cases
		"unknown_type_ampersand": {
			label:          "unknown type &",
			input:          "&5\r\nhello\r\n",
			wantSessionErr: "Malformed data : unknown type",
		},
		"unknown_type_plain_text": {
			label:          "plain text PING",
			input:          "PING\r\n",
			wantSessionErr: "Malformed data : unknown type",
		},
		"malformed_bulk_length_no_digits": {
			label:       "$ no length",
			input:       "$\r\n",
			wantReadErr: "Malformed data",
		},
	}

	for name, tc := range testCase {
		t.Run(name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tc.input))
			session, err := NewSession(reader)

			if err != nil {
				if tc.wantSessionErr == "" {
					t.Fatalf("NewSession unexpected error: %v", err)
				}
				if err.Error() != tc.wantSessionErr {
					t.Fatalf("NewSession error mismatch: got %q, want %q", err.Error(), tc.wantSessionErr)
				}
				return
			}
			if tc.wantSessionErr != "" {
				t.Fatalf("expected NewSession error %q, got nil", tc.wantSessionErr)
			}

			err = session.Read()
			if err != nil {
				if tc.wantReadErr == "" {
					t.Fatalf("Read unexpected error: %v", err)
				}
				if err.Error() != tc.wantReadErr {
					t.Fatalf("Read error mismatch: got %q, want %q", err.Error(), tc.wantReadErr)
				}
				return
			}
			if tc.wantReadErr != "" {
				t.Fatalf("expected Read error %q, got nil", tc.wantReadErr)
			}

			if !readEqual(session.content, tc.want) {
				t.Fatalf("content mismatch:\ngot  %#v\nwant %#v", session.content, tc.want)
			}
		})
	}
}

func readEqual(a, b any) bool {
	switch x := a.(type) {
	case string:
		y, ok := b.(string)
		return ok && x == y
	case int:
		y, ok := b.(int)
		return ok && x == y
	case []string:
		y, ok := b.([]string)
		if !ok || len(x) != len(y) {
			return ok && false
		}
		for i := range x {
			if x[i] != y[i] {
				return false
			}
		}
		return true
	case error:
		y, ok := b.(error)
		return ok && x.Error() == y.Error()
	default:
		return false
	}
}

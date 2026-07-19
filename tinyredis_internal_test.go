package main

import "testing"




func TestStorage(t *testing.T) {
    testCase := map[string]struct {
        testName    string
        want        any
        run         func(*rdb) any
    } {
        "insert item": {
            testName: "new hp movie", want: map[string]string{
                "HP 1": "Harry Potter and the Philosopher's Stone",
                "HP 2": "Harry Potter and the Chamber of Secrets",
                "HP 3": "Harry Potter and the Prisoner Of Azkabat",
                "HP 4": "Harry Potter and the Goblet of Fire",
                "HP 5": "Harry Potter and the Order of the Phoenix",
                "HP 6": "Harry Potter and the Half-Blood Prince",
                "HP 7": "Harry Potter and the Deathly Hallows - Part 1",
                "HP 8": "Harry Potter and the Deathly Hallows - Part 2",
                "HP 9": "Harry Potter and their gibberish",
            }, run: func(db *rdb) any {
                db.SET("HP 9", "Harry Potter and their gibberish")
                return db.rmap
            },
        },
        "update item": {
            testName: "update new hp movie", want: map[string]string{
                "HP 1": "Harry Potter and the Philosopher's Stone",
                "HP 2": "Harry Potter and the Chamber of Secrets",
                "HP 3": "Harry Potter and the Prisoner Of Azkabat",
                "HP 4": "Harry Potter and the Goblet of Fire",
                "HP 5": "Harry Potter and the Order of the Phoenix",
                "HP 6": "Harry Potter and the Half-Blood Prince",
                "HP 7": "Harry Potter and the Deathly Hallows - Part 1",
                "HP 8": "Harry Potter and the Deathly Hallows - Part 2",
                "HP 9": "Harry Potter and fantastic beast",
            }, run: func(db *rdb) any {
                db.SET("HP 9", "Harry Potter and fantastic beast")
                return db.rmap
            },
        },
        "delete existing item": {
            testName: "remove new hp movie", want: true, run: func(db *rdb) any {
                return db.DEL("HP 9")
            },
        },
        "delete not existed item": {
            testName: "remove new hp movie", want: false, run: func(db *rdb) any {
                return db.DEL("HP 11")
            },
        },
        "get hp movie 1": {
            testName: "get harry potter first movie name", want: "Harry Potter and the Philosopher's Stone", run: func(db *rdb) any {
                return db.GET("HP 1")
            },
        },
    }

    rdb := CreateAndSeed()
    for name,tc := range testCase {
        t.Run(name,func(t *testing.T) {
            got := tc.run(rdb)

            if !Equal(t, got,tc.want) {
                t.Fatalf("different result: got %v expected %v", got,tc.want)
            }
        })
    }
}

func Equal(t *testing.T, a, b any) bool {
    t.Helper()

	switch x := a.(type) {
	case string:
		y, ok := b.(string)
		return ok && x == y

	case map[string]string:
		y, ok := b.(map[string]string)
		if !ok {
			return false
		}

		if len(x) != len(y) {
			return false
		}

		for k, v := range x {
			if y[k] != v {
				return false
			}
		}

		return true

	case bool:
	    y,ok := b.(bool)
		return ok && x == y
	default:
		return false
	}
}

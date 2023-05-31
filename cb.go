package jacksapi

import (
	"bufio"
	"fmt"
	"net/http"
)

func Output(resp *http.Response, err error) {
	With(func(r rune, done bool, err error) {
		if done {
			fmt.Println()
		} else if err != nil {
			panic(err)
		} else {
			fmt.Print(string(r))
		}
	})(resp, err)
}

func With(f func(r rune, done bool, err error)) Callback {
	return func(resp *http.Response, err error) {
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		reader := bufio.NewReader(resp.Body)

		for {
			r, _, err := reader.ReadRune()
			if err != nil {
				if err.Error() == "EOF" {
					f(-1, true, nil)
					return
				}
				f(-1, false, err)
			}

			f(r, false, nil)
		}
	}
}

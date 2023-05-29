package cb

import (
	"io"
	"net/http"
	"os"
)

func CopyToStdio(resp *http.Response, err error) {
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}

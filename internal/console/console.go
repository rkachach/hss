package console

import (
	"fmt"
	"net/http"
)

func ConsoleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the management console!")
}

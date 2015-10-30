package pigae

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"appengine"
)

// Bad stuff started 2015-10...
// Had to take action.

var spammerSet = map[string]bool{}

func init() {
	b, err := ioutil.ReadFile("spammers.csv")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	full := string(b)
	spammers := strings.Split(full, ",")
	spammerSet = make(map[string]bool, len(spammers))
	for _, spammer := range spammers {
		spammerSet[spammer] = true
	}
	for _, ip := range os.Args[1:] {
		fmt.Println(ip, ":", spammerSet[ip])
	}
	// How to log this at startup ...?
	// fmt.Println("Registered list of", len(spammerSet), "spam IPs.")
}

func isSpam(w http.ResponseWriter, r *http.Request) bool {
	ip := r.RemoteAddr
	if spammerSet[ip] {
		c := appengine.NewContext(r)
		c.Infof("Detected spammer %v", ip)
		// Let's return a nice 200 ... nothing to see here
		fmt.Fprintln(w, "<html><body>This site is under construction</body></html>")
		return true
	}
	return false
}

package cert

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func HelloHandler(world string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, false)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		log.Printf("%s", dump)

		who := world
		if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
			if len(r.TLS.PeerCertificates[0].EmailAddresses) > 0 {
				who = r.TLS.PeerCertificates[0].EmailAddresses[0]
			}
		}
		fmt.Fprintf(w, "hello %v", who)
	}
}

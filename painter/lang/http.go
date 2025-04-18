package lang

import (
	"io"
	"log"
	"net/http"
	"strings"
	"fmt"

	"github.com/maxnetyaga/software-architecture-lab3/painter"
)

func HttpHandler(loop *painter.Loop, p *Parser) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var in io.Reader = r.Body
		if r.Method == http.MethodGet {
			in = strings.NewReader(r.URL.Query().Get("cmd"))
		}

		cmds, err := p.Parse(in)
		if err != nil {
			log.Printf("Bad script: %s", err)
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf("Error parsing commands: %s", err)))
			return
		}

		for _, cmd := range cmds {
			loop.Post(cmd)
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Commands received and posted to the event loop"))
	})
}

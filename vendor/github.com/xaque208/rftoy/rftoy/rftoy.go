package rftoy

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type RFToy struct {
	Address string
}

type state bool

func (t *RFToy) call(sid int, turn string) {
	addr := fmt.Sprintf("%s/cc", t.Address)

	req, err := http.NewRequest("GET", addr, nil)
	q := req.URL.Query()
	q.Add("sid", strconv.Itoa(sid))
	q.Add("turn", turn)

	if err != nil {
		log.Error(err)
	}

	client := &http.Client{}

	req.URL.RawQuery = q.Encode()
	client.Do(req)
}

func (t *RFToy) On(sid int) {
	t.call(sid, "on")
}

func (t *RFToy) Off(sid int) {
	t.call(sid, "off")
}

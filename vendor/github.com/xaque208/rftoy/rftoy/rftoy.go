package rftoy

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/common/log"
)

type RFToy struct {
	Address string
}

func (t *RFToy) call(sid int, turn string) error {
	addr := fmt.Sprintf("%s/cc", t.Address)

	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("sid", strconv.Itoa(sid))
	q.Add("turn", turn)

	client := &http.Client{}

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Response: %+v", resp)
		return err
	}

	return nil
}

func (t *RFToy) On(sid int) error {
	err := t.call(sid, "on")
	return err
}

func (t *RFToy) Off(sid int) error {
	err := t.call(sid, "off")
	return err
}

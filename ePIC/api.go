package ePIC

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/exp/maps"
	"github.com/jak3kaj/ePICmon/log"
)

type Host struct {
	Host    string
	Port    int
	timeout int
	counter int
	Error   error
}

type POSTPerpetualTune struct {
	Params   *POSTTuneParams `json:"param"`
	Password string          `json:"password"`
}

type POSTTuneParams struct {
	Algo   string `json:"algo"`
	Target int    `json:"target"`
}

func makeHost(host string) *Host {
	return &Host{
		Host:    host,
		Port:    4028,
		timeout: 5,
	}
}

func GetLog(host string) (log.Logs, error) {
	h := makeHost(host)
	ls := h.getLog()
	if ls == nil || h.Error != nil {
		return nil, h.Error
	}
	return *ls, nil
}

func (h *Host) getLog() *log.Logs {
	endpoint := "log"
	respData := h.getIt(endpoint)
	if respData == nil {
		return nil
	}

	// Unmarshall will throw an error by design
	if ls, err := log.UnmarshalJSON(*respData); err != nil {
		h.Error = fmt.Errorf("Failed to Unmarshall JSON from %s endpoint. Response Body: %s\n", endpoint, err)
	} else {
		return &ls
	}
	return nil
}

func GetSummary(host string) *Summary {
	h := makeHost(host)
	return h.getSummary()
}

func (h *Host) getSummary() *Summary {
	sum := &Summary{}
	h.counter += 1

	respData := h.getIt("summary")
	if respData == nil {
		return nil
	}

	if err := json.Unmarshal(*respData, sum); err != nil {
		h.Error = fmt.Errorf("Failed to Unmarshall JSON from response Body: %s\n", err)
		return nil
	}

	if sum == nil || sum.Result == nil {
		return sum
	}

	if h.counter <= h.timeout {
		return h.getSummary()
	}

	return nil
}

func (h *Host) getIt(endpoint string) *[]byte {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/%s", h.Host, h.Port, endpoint))

	if err != nil {
		h.Error = fmt.Errorf("Failed to Get data: %s\n", err)
		return nil 
	}

	if resp.StatusCode < 200 && resp.StatusCode >= 400 {
		h.Error = fmt.Errorf("Failed to Get data: %s\n", resp.Status)
		return nil
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		h.Error = fmt.Errorf("Failed to read response Body: %s\n", err)
		return nil
	}
	resp.Body.Close()

	return &respData
}

func ResetThrottle(host string, s *Summary) bool {
	h := makeHost(host)
	p := &POSTTuneParams{
		Algo:   maps.Keys(s.PerpetualTune.Algorithm)[0],
		Target: s.PerpetualTune.Algorithm[maps.Keys(s.PerpetualTune.Algorithm)[0]].Tgt,
	}
	pt := &POSTPerpetualTune{
		Params:   p,
		Password: "Kr0d@D0rk!",
	}
	return h.resetThrottle(pt)
}

func (h *Host) resetThrottle(pt *POSTPerpetualTune) bool {
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(pt)
	_, err := http.Post(fmt.Sprintf("http://%s:%d/%s", h.Host, h.Port, "perpetualtune/algo"), "application/json", payloadBuf)
	if err != nil {
		return false
	} else {
		return true
	}
}

func GetBoards(host string) *[3]log.Board {
	var ls log.Logs
	var err error
	if ls, err = GetLog(host); err != nil {
		fmt.Printf("GetLog Error: %#v\n", err)
	}

	//log.GetBoard(
	logBytes := []byte{}
	for _, l := range ls {
		logBytes = append(logBytes, *l.Bytes()...)
		logBytes = append(logBytes, []byte("\n")...)
	}

	var boards []*log.Board
	if err := log.FindBoards(&logBytes, &boards); err != nil {
		fmt.Printf("Error running log.FindBoards: %s\n", err)
	}
	var a [3]log.Board
	for _, b := range boards {
		a[b.HB] = *b
	}
	return &a
}

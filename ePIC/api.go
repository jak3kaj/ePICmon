package ePIC

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jak3kaj/ePICmon/log"
	"golang.org/x/exp/maps"
)

const Password = "Kr0d@D0rk!"

type Host struct {
	Host    string
	Port    int
	timeout int
	counter int
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
	ls, err := h.getLog()
	if err != nil {
		return nil, err
	} else if ls == nil {
		return nil, fmt.Errorf("getLog returned nil ls\n")
	}
	return *ls, nil
}

func (h *Host) getLog() (*log.Logs, error) {
	endpoint := "log"
	respData, err := h.getIt(endpoint)
	if err != nil {
		return nil, err
	} else if respData == nil {
		return nil, fmt.Errorf("getIt returned nil ls\n")
	}

	// Unmarshall will throw an error by design
	if ls, err := log.UnmarshalJSON(*respData); err == nil {
		return &ls, nil
	}
	err = fmt.Errorf("Failed to Unmarshall JSON from %s endpoint. Response Body: %s\n", endpoint, err)
	return nil, err
}

func GetSummary(host string) (*Summary, error) {
	h := makeHost(host)
	return h.getSummary()
}

func (h *Host) getSummary() (*Summary, error) {
	sum := &Summary{}
	h.counter += 1

	respData, err := h.getIt("summary")
	if err != nil {
		return nil, err
	} else if respData == nil {
		return nil, fmt.Errorf("Getting Summary returned nil\n")
	}

	if err := json.Unmarshal(*respData, sum); err != nil {
		return nil, fmt.Errorf("Failed to Unmarshall JSON from response Body: %s\n", err)
	}

	if sum == nil || sum.Result == nil {
		return sum, nil
	}

	if h.counter <= h.timeout {
		return h.getSummary()
	}

	return nil, nil
}

func (h *Host) getIt(endpoint string) (*[]byte, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/%s", h.Host, h.Port, endpoint))

	if err != nil {
		return nil, fmt.Errorf("Failed to Get data: %s\n", err)
	}

	if resp.StatusCode < 200 && resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Failed to Get data: %s\n", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to read response Body: %s\n", err)
	}

	return &respData, nil
}

func ResetThrottle(host string, s *Summary) bool {
	h := makeHost(host)
	p := &POSTTuneParams{
		Algo:   maps.Keys(s.PerpetualTune.Algorithm)[0],
		Target: s.PerpetualTune.Algorithm[maps.Keys(s.PerpetualTune.Algorithm)[0]].Tgt,
	}
	pt := &POSTPerpetualTune{
		Params:   p,
		Password: Password,
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

func GetBoards(host string) (*[3]log.Board, error) {
	var ls log.Logs
	var err error = nil
	if ls, err = GetLog(host); err != nil {
		// Don't report an error
		return nil, nil
	}

	//log.GetBoard(
	logBytes := []byte{}
	for _, l := range ls {
		logBytes = append(logBytes, *l.Bytes()...)
		logBytes = append(logBytes, []byte("\n")...)
	}

	var boards []*log.Board
	if err = log.FindBoards(&logBytes, &boards); err != nil {
		err = fmt.Errorf("Error running log.FindBoards: %s\n", err)
	}
	var a [3]log.Board
	for _, b := range boards {
		a[b.HB] = *b
	}
	return &a, err
}

package log

import (
	"bytes"
	"encoding/json"
	"io"
	"fmt"
	"regexp"
	"time"
	"gopkg.in/yaml.v2"

	"github.com/acarl005/stripansi"
)

type TS struct {
	S int `json:"secs_since_epoch"`
	NSs int `json:"nanos_since_epoch"`
}

type LogSection struct {
	TS *TS
	Log *string
}

type LogSectionTypes interface {
	TS | string
}

func (ts TS) String() string {
	t := time.Unix(int64(ts.S), int64(ts.NSs)).UTC()
	return fmt.Sprintf("%s", t)
	//return fmt.Sprintf("%d.%d", ts.S, ts.NSs)
}

func (l LogSection) String() string {
	var r string
	r += fmt.Sprintf("%s\n", l.TS)
	if l.Log != nil {
		r += fmt.Sprintf("%s\n", *l.Log)
	}
	return r
}

func (l LogSection) Bytes() *[]byte {
	var r []byte
	if l.Log != nil {
		r = []byte(*l.Log)
	}
	return &r
}

type TSList [][]TS
type StringList [][]string
type Logs []*LogSection

func UnmarshalJSON(data []byte) (Logs, error) {
	// First pass pul out TSs
	var logParseTS TSList
	var ls Logs
	// Note json.Unmarshal will throw errors because there are mixed types in the json output
	json.Unmarshal(data, &logParseTS)
	lnil := TS{}
	for _, a := range logParseTS{
		// Make a nil log.TS object
		for _, b := range a {
			// Compare b to a nil l.TS object
			if b != lnil {
				l := LogSection{}
				bb := b
				l.TS = &bb
				ls = append(ls, &l)
			}
		}
	}

	// Second pass pull out log strings
	var logParseString StringList
	// Note json.Unmarshal will throw errors because there are mixed types in the json output
	json.Unmarshal(data, &logParseString)
	for i, a := range logParseString{
		for _, b := range a {
			if b != "" {
				ls[i].Log = &b
			}
		}
	}
	return ls, nil
}

type Board struct {
	HB int `yaml:"HB"` //   0
	SerialNum string `yaml:"Board Serial No"` //   NGSBYPDBCJHAA0BKC
	ChipDie string `yaml:"Chip Die"` //          ED
	ChipMarking string `yaml:"Chip Marking"` //      S1GX23BF1L
	ChipBin int `yaml:"Chip Bin"` //          4
	FTVersion string `yaml:"FT Version"` //        F1V22B3C1
	PCBVersion int `yaml:"PCB Version"` //       220
	BOMVersion int `yaml:"BOM Version"` //       0
	ASICSensorType string  `yaml:"ASIC Sensor Type"` //  Some(0)
	ASICSensorAddr0 string `yaml:"ASIC Sensor Addr0"` // Some(0)
	ASICSensorAddr1 string `yaml:"ASIC Sensor Addr1"` // Some(0)
	ASICSensorAddr2 string `yaml:"ASIC Sensor Addr2"` // Some(0)
	ASICSensorAddr3 string `yaml:"ASIC Sensor Addr3"` // Some(0)
	PICSensorType string `yaml:"PIC Sensor Type"` //   Some(0)
	PICSensorAddr string `yaml:"PIC Sensor Addr"` //   Some(0)
	ChipTech string `yaml:"Chip Tech"` //         AL
	BoardName string `yaml:"Board Name"` //        Some("BHB56801")
	FactoryJob string `yaml:"Factory Job"` //       Some("NGSB20230801001")
	DefaultVolt int `yaml:"Default Volt"` //      12900
	DefaultClk int `yaml:"Default Clk"` //       485
	NonceRate string `yaml:"Nonce Rate"` //        Some(9950)
	PCBTempIn string `yaml:"PCB Temp In"` //       Some(0)
	PCBTempOut string `yaml:"PCB Temp Out"` //      Some(0)
	TestVersion string `yaml:"Test Version"` //      Some(0)
	TestStandard string `yaml:"Test Standard"` //     Some(1)
	PT2Result string `yaml:"PT2 Result"` //        Some(1)
	PT2Count string `yaml:"PT2 Count"` //         Some(2)
}

func FindBoards(data *[]byte, boards *[]*Board) error {
	// Hopefully stripping all ansi colors from the entire log output and converting from string to []byte is efficient enough.
	*data = []byte(stripansi.Strip(string(*data)))

	// Get the HB id and add a document separator before the HB row.
	// [2024-01-19 19:10:32][bm_miner::controller][INFO] - is board id 0 detected?: true
	// [2024-01-19 19:10:32][bm_miner::controller][INFO] - HB: 0
	re, err  := regexp.Compile("(?m)^" + regexp.QuoteMeta("[") +
		"[[:digit:]]{4}-[[:digit:]]{2}-[[:digit:]]{2} [[:digit:]]{2}:[[:digit:]]{2}:[[:digit:]]{2}" +
		regexp.QuoteMeta("]") + regexp.QuoteMeta("[") + "bm_miner::controller" +
		regexp.QuoteMeta("][") + "INFO" + regexp.QuoteMeta("]") +  " - HB:")
	if err != nil {
		return fmt.Errorf("Regexp Compile failed: %s\n", err)
	}
	*data  = re.ReplaceAll(*data, []byte("| ---\n| HB:"))

	re, err  = regexp.Compile("(?m)^" + regexp.QuoteMeta("|") + " (---|.+:.+)$")
	if err != nil {
		return fmt.Errorf("Regexp Compile failed: %s\n", err)
	}
	find := re.FindAllSubmatchIndex(*data, -1)

	foundYAML := []byte{}
	for _, in := range find {
		foundYAML = append(foundYAML, (*data)[in[2]:in[3]]...)
		foundYAML = append(foundYAML, []byte("\n")...)
	}

	// Uncomment to see the YAML output
	// fmt.Printf("%s\n", foundYAML)

	decoder := yaml.NewDecoder(bytes.NewBuffer(foundYAML))
	for {
	    var d Board
	    if err := decoder.Decode(&d); err != nil {
	        if err == io.EOF {
	            break
	        }
	        return fmt.Errorf("Document decode failed: %w", err)
	    }
		*boards = append(*boards, &d)
	}
	return nil
}

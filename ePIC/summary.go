package ePIC

type Status struct {
	OperState      string            `json:"Operating State"`
	LastCmd        string            `json:"Last Command"`
	LastCmdResults map[string]string `json:"Last Command Result"`
	Lasterr        string            `json:"Last Error"`
}

type Stratum struct {
	ConfigID         int     `json:"Config Id"`
	CurrentPool      string  `json:"Current Pool"`
	CurrentUser      string  `json:"Current User"`
	PoolConnect      bool    `json:"IsPoolConnected"`
	AvgLatency       float64 `json:"Average Latency"`
	WorkerUniqID     bool    `json:"Worker Unique Id"`
	WorkerUniqIDType string  `json:"Worker Unique Id Variant"`
}

type Session struct {
	StartupTS        int        `json:"Startup Timestamp"`
	StartupTime      string     `json:"Startup String"`
	Uptime           int        `json:"Uptime"`
	LastWorkTS       int        `json:"Last Work Timestamp"`
	WorkRX           int        `json:"WorkReceived"`
	ActiveHashboards int        `json:"Active HBs"`
	AvgMHs           float64    `json:"Average MHs"`
	LastAvgMHs       LastAvgMHs `json:"LastAverageMHs"`
	Accepted         int        `json:"Accepted"`
	Rejeted          int        `json:"Rejected"`
	Submitted        int        `json:"Submitted"`
	LastAcceptedTS   int        `json:"Last Accepted Share Timestamp"`
	Difficulty       float64    `json:"Difficulty"`
}

type LastAvgMHs struct {
	Hashrate float64 `json:"Hashrate"`
	TS       int     `json:"Timestamp"`
}

type Hashboard struct {
	ID       int       `json:"Index"`
	In_v     float64   `json:"Input Voltage"`
	Out_v    float64   `json:"Output Voltage"`
	In_a     float64   `json:"Input Current"`
	Out_a    float64   `json:"Output Current"`
	In_w     float64   `json:"Input Power"`
	Out_w    float64   `json:"Output Power"`
	Temp     float64   `json:"Temperature"`
	ClkList  []float64 `json:"Core Clock"`
	Hashrate []float64 `json:"Hashrate"`
	ClkAvg   float64   `json:"Core Clock Avg"`
}

type Fans struct {
	Speed int            `json:"Fans Speed"`
	Mode  map[string]int `json:"Fan Mode"`
}

type Misc struct {
	Locate       bool    `json:"Locate Miner State"`
	ShutdownTemp float64 `json:"Shutdown Temp"`
}

type StratumConfig struct {
	Pool  string `json:"pool"`
	Login string `json:"login"`
	Pass  string `json:"password"`
}

type PerpetualTune struct {
	Running   bool                  `json:"Running"`
	Algorithm map[string]TuneParams `json:"Algorithm"`
}

type TuneParams struct {
	Optimized   bool   `json:"Optimized"`
	Tgt         int    `json:"Target"`
	ThrottleTgt int    `json:"Throttle Target"`
	Unit        string `json:"Unit"`
}

type PsuStats struct {
	In_v  float64 `json:"Input Voltage"`
	Out_v float64 `json:"Output Voltage"`
	In_a  float64 `json:"Input Current"`
	Out_a float64 `json:"Output Current"`
	In_w  float64 `json:"Input Power"`
	Out_w float64 `json:"Output Power"`
	Tgt_v int     `json:"Target Voltage"`
}

type HwConfig struct {
	TgtClk []Clk `json:"Boards Target Clock"`
}

type Clk struct {
	ID   int     `json:"Index"`
	Data float64 `json:"Data"`
}

type Summary struct {
	Status        Status            `json:"Status"`
	Hostname      string            `json:"Hostname"`
	PresetInfo    map[string]int    `json:"PresetInfo"`
	Software      string            `json:"Software"`
	Mining        map[string]string `json:"Mining"`
	Stratum       Stratum           `json:"Stratum"`
	Session       Session           `json:"Session"`
	Hashboards    []Hashboard       `json:"HBs"`
	Fans          Fans              `json:"Fans"`
	FanRPMs       map[string]int    `json:"Fans Rpm"`
	Misc          Misc              `json:"Misc"`
	StratumConfig []StratumConfig   `json:"StratumConfigs"`
	PerpetualTune PerpetualTune     `json:"PerpetualTune"`
	PsuStats      PsuStats          `json:"Power Supply Stats"`
	HwConfig      HwConfig          `json:"HwConfig"`
	Result        *bool             `json:"result"`
	Error         string            `json:"error"`
}

type Board struct {
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

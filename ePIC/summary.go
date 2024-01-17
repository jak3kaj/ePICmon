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
	AvgLatency       float32 `json:"Average Latency"`
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
	AvgMHs           float32    `json:"Average MHs"`
	LastAvgMHs       LastAvgMHs `json:"LastAverageMHs"`
	Accepted         int        `json:"Accepted"`
	Rejeted          int        `json:"Rejected"`
	Submitted        int        `json:"Submitted"`
	LastAcceptedTS   int        `json:"Last Accepted Share Timestamp"`
	Difficulty       float32    `json:"Difficulty"`
}

type LastAvgMHs struct {
	Hashrate float32 `json:"Hashrate"`
	TS       int     `json:"Timestamp"`
}

type Hashboard struct {
	ID       int       `json:"Index"`
	In_v     float32   `json:"Input Voltage"`
	Out_v    float32   `json:"Output Voltage"`
	In_a     float32   `json:"Input Current"`
	Out_a    float32   `json:"Output Current"`
	In_w     float32   `json:"Input Power"`
	Out_w    float32   `json:"Output Power"`
	Temp     float32   `json:"Temperature"`
	ClkList  []float32 `json:"Core Clock"`
	Hashrate []float32 `json:"Hashrate"`
	ClkAvg   float32   `json:"Core Clock Avg"`
}

type Fans struct {
	Speed int            `json:"Fans Speed"`
	Mode  map[string]int `json:"Fan Mode"`
}

type Misc struct {
	Locate       bool    `json:"Locate Miner State"`
	ShutdownTemp float32 `json:"Shutdown Temp"`
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
	In_v  float32 `json:"Input Voltage"`
	Out_v float32 `json:"Output Voltage"`
	In_a  float32 `json:"Input Current"`
	Out_a float32 `json:"Output Current"`
	In_w  float32 `json:"Input Power"`
	Out_w float32 `json:"Output Power"`
	Tgt_v int     `json:"Target Voltage"`
}

type HwConfig struct {
	TgtClk []Clk `json:"Boards Target Clock"`
}

type Clk struct {
	ID   int     `json:"Index"`
	Data float32 `json:"Data"`
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

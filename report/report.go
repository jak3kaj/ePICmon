package report

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/jak3kaj/ePICmon/Ocean"
	"github.com/jak3kaj/ePICmon/ePIC"
)

func Report(s *ePIC.Summary, o *Ocean.UserTable) string {
	//Mining
	rpt := s.Status.OperState
	rpt += fmt.Sprintf(" - Uptime: %02d:%02d:%02d", uptime(s.Session.Uptime)...)

	//Poll Connected
	if s.Stratum.PoolConnect {
		rpt += fmt.Sprintf(" - Pool: Connected")
	} else {
		rpt += fmt.Sprintf(" - Pool: Not Connected")
	}
	rpt += "\n"

	rpt += fmt.Sprintf("%.0fW %.1fV", s.PsuStats.In_w, s.PsuStats.Out_v)

	if len(s.Fans.Mode) > 0 {
		rpt += fmt.Sprintf(" %s", maps.Keys(s.Fans.Mode)[0])
	}

	rpt += fmt.Sprintf(" Fan: %d%%\n", s.Fans.Speed)

	var sumMH float32
	var sumC float32

	for _, Hashboard := range s.Hashboards {
		rpt += fmt.Sprintf("\tBoard %d", Hashboard.ID)
		if len(Hashboard.Hashrate) > 0 {
			rpt += fmt.Sprintf(" %.2fTHs", Hashboard.Hashrate[0]/1000000)
			sumMH += Hashboard.Hashrate[0]
			sumC += Hashboard.Temp
		}
		rpt += fmt.Sprintf(" %.1fC %.2fMHz\n", Hashboard.Temp, Hashboard.ClkAvg)
	}
	rpt += fmt.Sprintf("Avg: %.2fTHs %.1fC Efficiency: %.1fJ/TH\n", sumMH/1000000, sumC/3, s.PsuStats.In_w/(sumMH/1000000))

	if s.PerpetualTune.Running {
		algorithm := maps.Keys(s.PerpetualTune.Algorithm)[0]
		rpt += fmt.Sprintf("Tgt: %d.00THs", s.PerpetualTune.Algorithm[algorithm].Tgt)
		rpt += fmt.Sprintf(" %s", algorithm)
		if s.PerpetualTune.Algorithm[algorithm].Optimized {
			rpt += fmt.Sprintf(" - Optimized\n")
		} else {
			rpt += fmt.Sprintf(" - Not Optimized\n")
		}
	} else {
		rpt += "\n"
	}
	rpt += fmt.Sprintf("Pool: %sTHs Earnings: %s\n", strings.TrimSuffix(o.Hashrate3hr, " Th/s"), o.Earnings)

	return rpt
}

func Status(s *ePIC.Summary) string {
	//Mining
	rpt := fmt.Sprintf("Status: %s\n", s.Status.OperState)
	rpt += fmt.Sprintf("Uptime: %03d:%02d:%02d\n", uptime(s.Session.Uptime)...)

	//Poll Connected
	if s.Stratum.PoolConnect {
		rpt += fmt.Sprintf("Pool:   Connected\n")
	} else {
		rpt += fmt.Sprintf("Pool:   Not Connected\n")
	}
	return rpt

}

func Performance(s *ePIC.Summary, o *Ocean.UserTable) string {
	var sumMH float32
	var rpt string

	if s.PerpetualTune.Running {
		algorithm := maps.Keys(s.PerpetualTune.Algorithm)[0]
		rpt += fmt.Sprintf("Tgt:  %d.00THs", s.PerpetualTune.Algorithm[algorithm].Tgt)
		rpt += fmt.Sprintf(" %s:", algorithm)
		if s.PerpetualTune.Algorithm[algorithm].Optimized {
			rpt += fmt.Sprintf(" Optimized\n")
		} else {
			rpt += fmt.Sprintf(" Not Optimized\n")
		}
	}

	for _, Hashboard := range s.Hashboards {
		if len(Hashboard.Hashrate) > 0 {
			sumMH += Hashboard.Hashrate[0]
		}
	}

	rpt += fmt.Sprintf("Avg:  %.2fTHs\n", sumMH/1000000)

	rpt += fmt.Sprintf("Pool: %s0THs Earnings: %s\n", strings.TrimSuffix(o.Hashrate3hr, " Th/s"), o.Earnings)

	return rpt
}

func Psu(s *ePIC.Summary) string {
	var sumC float32
	var sumMH float32

	for _, Hashboard := range s.Hashboards {
		if len(Hashboard.Hashrate) > 0 {
			sumC += Hashboard.Temp
			sumMH += Hashboard.Hashrate[0]
		}
	}

	rpt := fmt.Sprintf("Efficiency: %2.1fJ/TH\n", s.PsuStats.In_w/(sumMH/1000000))
	rpt += fmt.Sprintf("%2.1fV       %4.0fW\n", s.PsuStats.Out_v, s.PsuStats.In_w)

	rpt += fmt.Sprintf("Fan: %d%%   %.1fC", s.Fans.Speed, sumC/3)
	if len(s.Fans.Mode) > 0 {
		rpt += fmt.Sprintf("    Mode: %s", maps.Keys(s.Fans.Mode)[0])
	}
	rpt += "\n"


	return rpt
}

func Board(s *ePIC.Summary) string {
	var rpt string
	for _, Hashboard := range s.Hashboards {
		rpt += fmt.Sprintf("Board %d", Hashboard.ID)
		if len(Hashboard.Hashrate) > 0 {
			rpt += fmt.Sprintf(" %.2fTHs", Hashboard.Hashrate[0]/1000000)
		}
		rpt += fmt.Sprintf(" %.1fC %.2fMHz\n", Hashboard.Temp, Hashboard.ClkAvg)
	}

	return rpt
}

func uptime(d int) []any {
	h := d / 3600
	m := d % 3600 / 60
	s := d % 3600 % 60
	return []any{h, m, s}
}

package report

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/jak3kaj/ePICmon/Ocean"
	"github.com/jak3kaj/ePICmon/ePIC"
	"github.com/jak3kaj/ePICmon/log"
	"github.com/jak3kaj/ePICmon/power"
)

const Sqrt3 float64 = 1.7320508075688772935274463415058723669428052538103806280558069794519330169088000370811461867572485756756261414154 //https://oeis.org/A002194

func Report(s *ePIC.Summary, o *Ocean.UserTable) string {
	var rpt string
	if s != nil {
		//Mining
		rpt += s.Status.OperState
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

		var sumMH float64
		var sumC float64

		for _, Hashboard := range s.Hashboards {
			rpt += fmt.Sprintf("\tBoard %d", Hashboard.ID)
			if len(Hashboard.Hashrate) > 0 {
				rpt += fmt.Sprintf(" %.2fTHs", Hashboard.Hashrate[0]/1000000)
				sumMH += Hashboard.Hashrate[0]
				sumC += Hashboard.Temp
			}
			rpt += fmt.Sprintf(" %.1fC %.2fV %.2fMHz\n", Hashboard.Temp, Hashboard.In_v, Hashboard.ClkAvg)
		}
		rpt += fmt.Sprintf("Avg: %.2fTHs %.1fC Efficiency: %.1fJ/TH\n", sumMH/1000000, sumC/3, s.PsuStats.In_w/(sumMH/1000000))

		if s.PerpetualTune.Running {
			algorithm := maps.Keys(s.PerpetualTune.Algorithm)[0]
			rpt += fmt.Sprintf("%s Tgt: %d.00THs", algorithm, s.PerpetualTune.Algorithm[algorithm].Tgt)
			if s.PerpetualTune.Algorithm[algorithm].Optimized {
				rpt += fmt.Sprintf(" Optimized\n")
			} else {
				rpt += fmt.Sprintf(" Not Optimized\n")
			}
		} else {
			rpt += "\n"
		}
		rpt += fmt.Sprintf("Pool: %sTHs Earnings: %s\n", strings.TrimSuffix(o.Hashrate3hr, " Th/s"), o.Earnings)
	}
	return rpt
}

func Status(s *ePIC.Summary) string {
	var rpt string
	if s != nil {
		//Mining
		rpt += fmt.Sprintf("Status: %s\n", s.Status.OperState)
		rpt += fmt.Sprintf("Uptime: %03d:%02d:%02d\n", uptime(s.Session.Uptime)...)

		//Poll Connected
		if s.Stratum.PoolConnect {
			rpt += fmt.Sprintf("Pool:   Connected\n")
		} else {
			rpt += fmt.Sprintf("Pool:   Not Connected\n")
		}
	}
	return rpt
}

func Performance(s *ePIC.Summary, o *Ocean.UserTable) string {
	var rpt string
	var sumMH float64

	if s != nil {
		if s.PerpetualTune.Running {
			algorithm := maps.Keys(s.PerpetualTune.Algorithm)[0]
			rpt += fmt.Sprintf("Tgt:  %d.00THs", s.PerpetualTune.Algorithm[algorithm].Tgt)
			if s.PerpetualTune.Algorithm[algorithm].ThrottleTgt > 0 {
				rpt += fmt.Sprintf(" [\"Th\"]Throttled: %d.00THs[\"\"]", s.PerpetualTune.Algorithm[algorithm].ThrottleTgt)
			}
			rpt += fmt.Sprintf(" %s:", algorithm)
			if s.PerpetualTune.Algorithm[algorithm].Optimized {
				rpt += fmt.Sprintf(" Optimized")
			} else {
				rpt += fmt.Sprintf(" Not Optimized")
			}
			rpt += "\n"
		}

		for _, Hashboard := range s.Hashboards {
			if len(Hashboard.Hashrate) > 0 {
				sumMH += Hashboard.Hashrate[0]
			}
		}

		rpt += fmt.Sprintf("Avg:  %.2fTHs Shutdown: %.2fC\n", sumMH/1000000, s.Misc.ShutdownTemp)

		if o != nil {
			rpt += fmt.Sprintf("Pool: %s0THs Earnings: %s\n", strings.TrimSuffix(o.Hashrate3hr, " Th/s"), o.Earnings)
		}
	}
	return rpt
}

func Psu(s *ePIC.Summary, Vc float64) string {
	var rpt string
	var sumC float64
	var sumMH float64

	if s != nil {
		for _, Hashboard := range s.Hashboards {
			if len(Hashboard.Hashrate) > 0 {
				sumC += Hashboard.Temp
				sumMH += Hashboard.Hashrate[0]
			}
		}

		rpt += fmt.Sprintf("Efficiency: %2.2fJ/TH\n", s.PsuStats.In_w/(sumMH/1000000))
		rpt += fmt.Sprintf("%2.2fV      %4.0fW     %1.1fA\n", s.PsuStats.Out_v, s.PsuStats.In_w, s.PsuStats.In_w/Vc)

		rpt += fmt.Sprintf("Fan: %d%%   %.1fC", s.Fans.Speed, sumC/3)
		/*
			if len(s.Fans.Mode) > 0 {
				rpt += fmt.Sprintf("     %s", maps.Keys(s.Fans.Mode)[0])
			}
		*/
		software, _ := strings.CutPrefix(s.Software, "PowerPlay-BM v")
		rpt += fmt.Sprintf("     %s\n", software)
	}

	return rpt
}

func Board(s *ePIC.Summary, b *[3]log.Board) string {
	var rpt string
	nilBoard := log.Board{}
	if s != nil {
		for _, Hashboard := range s.Hashboards {
			rpt += fmt.Sprintf("Board %d", Hashboard.ID)
			if len(Hashboard.Hashrate) > 0 {
				rpt += fmt.Sprintf(" %.2fTHs", Hashboard.Hashrate[0]/1000000)
			}
			rpt += fmt.Sprintf(" %.1fC %.2fV %.2fMHz ", Hashboard.Temp, Hashboard.In_v, Hashboard.ClkAvg)
			if b != nil && b[Hashboard.ID] != nilBoard {
				rpt += b[Hashboard.ID].SerialNum
			}
			rpt += "\n"
		}
	}

	return rpt
}

func Power(s *power.Power) []string {
	var rpt []string
	for _, p := range s.Panels {
		for _, c := range p.Circuits {
			var r string
			var totW float64
			for _, l := range c.Legs {
				a := l.W / p.V
				totW += l.W
				r += fmt.Sprintf("%s: %8.2fW %3.0fV %5.2fA %5.2fA\n", l.ID, l.W, p.V, a, a/Sqrt3)
			}
			r += fmt.Sprintf("     %8.2fW %3.0fV %5.2fA %5.2fA\n", totW, p.V, totW/p.V, totW/p.V/Sqrt3)
			rpt = append(rpt, r)
		}
	}
	return rpt
}

func TotalPower(s *power.Power) string {
	var totW float64
	var v float64
	for _, p := range s.Panels {
		for _, c := range p.Circuits {
			for _, l := range c.Legs {
				v = p.V
				totW += l.W
			}
		}
	}
	return fmt.Sprintf("%8.2fW %3.0fV %5.2fA %5.2fA\n", totW, v, totW/v, totW/v/Sqrt3)
}

func uptime(d int) []any {
	h := d / 3600
	m := d % 3600 / 60
	s := d % 3600 % 60
	return []any{h, m, s}
}

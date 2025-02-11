package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/maps"

	"github.com/jak3kaj/ePICmon/Ocean"
	"github.com/jak3kaj/ePICmon/ePIC"
	"github.com/jak3kaj/ePICmon/log"
	"github.com/jak3kaj/ePICmon/power"
	"github.com/jak3kaj/ePICmon/report"

	"github.com/rivo/tview"
)

type textView map[string]*tview.TextView
type textUpdate map[string]textView

type Model struct {
	hosts          []string
	power          *power.Power
	hostPower      map[string][]*power.Leg
	btcAddr        string
	siteData       map[string]*ePIC.Summary
	oceanData      map[string]*Ocean.UserTable
	textUpdate     textUpdate
	boardData      map[string]*[3]log.Board
	siteDataError  map[string]error
	boardDataError map[string]error
	app            *tview.Application
	mutex          *sync.RWMutex
}

func initModel() *Model {
	m := new(Model)
	m.hostPower = make(map[string][]*power.Leg)
	m.btcAddr = "bc1qluhcxmzf8up8m8625gtl74458jemt8jcgrp3u3"
	m.hosts = []string{"miner001", "miner002", "miner003", "miner004", "miner005", "miner006", "miner007", "miner008", "miner009", "miner010", "miner011", "miner012"}
	// m.hosts = []string{"192.168.1.11", "192.168.1.12", "192.168.1.13", "192.168.1.14", "192.168.1.15", "192.168.1.16", "192.168.1.17", "192.168.1.18", "192.168.1.19", "192.168.1.20", "192.168.1.21","192.168.1.22"}

	m.siteData = make(map[string]*ePIC.Summary)
	m.oceanData = make(map[string]*Ocean.UserTable)
	m.textUpdate = make(textUpdate)

	// Define power connections
	m.power = power.Init()

	// Circuit 1
	m.hostPower["192.168.1.19"] = []*power.Leg{m.power.Panels[0].Circuits[0].Legs[0], m.power.Panels[0].Circuits[0].Legs[1]}
	m.hostPower["192.168.1.21"] = []*power.Leg{m.power.Panels[0].Circuits[0].Legs[0], m.power.Panels[0].Circuits[0].Legs[2]}
	m.hostPower["192.168.1.22"] = []*power.Leg{m.power.Panels[0].Circuits[0].Legs[1], m.power.Panels[0].Circuits[0].Legs[2]}

	// Circuit 2
	m.hostPower["192.168.1.13"] = []*power.Leg{m.power.Panels[0].Circuits[1].Legs[0], m.power.Panels[0].Circuits[1].Legs[1]}
	m.hostPower["192.168.1.14"] = []*power.Leg{m.power.Panels[0].Circuits[1].Legs[0], m.power.Panels[0].Circuits[1].Legs[2]}
	m.hostPower["192.168.1.15"] = []*power.Leg{m.power.Panels[0].Circuits[1].Legs[1], m.power.Panels[0].Circuits[1].Legs[2]}

	// Circuit 3
	m.hostPower["192.168.1.17"] = []*power.Leg{m.power.Panels[0].Circuits[2].Legs[0], m.power.Panels[0].Circuits[2].Legs[2]}
	m.hostPower["192.168.1.18"] = []*power.Leg{m.power.Panels[0].Circuits[2].Legs[1], m.power.Panels[0].Circuits[2].Legs[0]}
	m.hostPower["192.168.1.11"] = []*power.Leg{m.power.Panels[0].Circuits[2].Legs[2], m.power.Panels[0].Circuits[2].Legs[1]}

	// Circuit 4
	m.hostPower["192.168.1.20"] = []*power.Leg{m.power.Panels[0].Circuits[3].Legs[0], m.power.Panels[0].Circuits[3].Legs[2]}
	m.hostPower["192.168.1.16"] = []*power.Leg{m.power.Panels[0].Circuits[3].Legs[0], m.power.Panels[0].Circuits[3].Legs[1]}
	m.hostPower["192.168.1.12"] = []*power.Leg{m.power.Panels[0].Circuits[3].Legs[1], m.power.Panels[0].Circuits[3].Legs[2]}

	m.hostPower["miner001"] = m.hostPower["192.168.1.11"]
	m.hostPower["miner002"] = m.hostPower["192.168.1.12"]
	m.hostPower["miner003"] = m.hostPower["192.168.1.13"]
	m.hostPower["miner004"] = m.hostPower["192.168.1.14"]
	m.hostPower["miner005"] = m.hostPower["192.168.1.15"]
	m.hostPower["miner006"] = m.hostPower["192.168.1.16"]
	m.hostPower["miner007"] = m.hostPower["192.168.1.17"]
	m.hostPower["miner008"] = m.hostPower["192.168.1.18"]
	m.hostPower["miner009"] = m.hostPower["192.168.1.19"]
	m.hostPower["miner010"] = m.hostPower["192.168.1.20"]
	m.hostPower["miner011"] = m.hostPower["192.168.1.21"]
	m.hostPower["miner012"] = m.hostPower["192.168.1.22"]

	m.mutex = &sync.RWMutex{}
	m.boardData = make(map[string]*[3]log.Board)
	for _, host := range m.hosts {
		m.boardData[host] = nil
	}

	m.siteDataError = make(map[string]error)
	m.boardDataError = make(map[string]error)

	return m
}

func main() {
	m := initModel()

	m.getData()

	newPrimitive := func(label string) *tview.TextView {
		return tview.NewTextView().
			SetTextAlign(tview.AlignLeft).
			SetLabel(label)

	}

	//grid := tview.NewGrid().SetBorders(true).SetColumns(-1, -2, -1, -2)
	grid := tview.NewGrid().SetBorders(true)

	var h int
	var i int
	//i += 1

	if m.textUpdate["Power"] == nil {
		m.textUpdate["Power"] = make(textView)
	}
	m.textUpdate["Power"]["Total"] = newPrimitive("")
    // AddItem(p Primitive, row, column, rowSpan, colSpan, minGridHeight, minGridWidth int, focus bool)
	grid.AddItem(m.textUpdate["Power"]["Total"], i, h, 2, 1, 2, 1, true)
	for pi, p := range m.power.Panels {
		for ci, _ := range p.Circuits {
			h += 1
			cName := fmt.Sprintf("Circuit %d ", pi+ci+1)
			m.textUpdate["Power"][cName] = newPrimitive("")
			grid.AddItem(m.textUpdate["Power"][cName], i, h, 2, 1, 2, 1, false)
		}
	}
    // Increment row count
    i += 1

	for _, host := range m.hosts {
        i += 1
		if m.textUpdate[host] == nil {
			m.textUpdate[host] = make(textView)
		}

		m.textUpdate[host]["status"] = newPrimitive(host + " ")
		grid.AddItem(m.textUpdate[host]["status"], i, 0, 1, 1, 1, 1, false)

		m.textUpdate[host]["host"] = newPrimitive("")
		grid.AddItem(m.textUpdate[host]["host"], i, 1, 1, 1, 1, 2, false)

		m.textUpdate[host]["psu"] = newPrimitive("")
		grid.AddItem(m.textUpdate[host]["psu"], i, 2, 1, 1, 1, 1, false)

		m.textUpdate[host]["board"] = newPrimitive("")
		grid.AddItem(m.textUpdate[host]["board"], i, 3, 1, 2, 1, 2, false)
	}

	m.app = tview.NewApplication()

	go m.refreshData(5)
	go m.getLogData(600)

	if err := m.app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}

}

func (m Model) clearThrottle(t int) {
	for {
		m.getData()
		for _, host := range m.hosts {
			if m.siteData[host].PerpetualTune.Running {
				algorithm := maps.Keys(m.siteData[host].PerpetualTune.Algorithm)[0]
				if m.siteData[host].PerpetualTune.Algorithm[algorithm].ThrottleTgt > 0 {
					m.app.QueueUpdateDraw(func() {
						m.mutex.RLock()
						m.textUpdate[host]["status"].Clear()
						m.textUpdate[host]["host"].
							SetText(fmt.Sprintf("Resetting Throttle..."))
						m.textUpdate[host]["psu"].Clear()
						m.textUpdate[host]["board"].Clear()
						m.mutex.RUnlock()
					})
					var status string
					if ok := ePIC.ResetThrottle(host, m.siteData[host]); !ok {
						status = "Resetting Throttle Failed"
					} else {
						status = "Resetting Throttle Failed"
					}
					m.app.QueueUpdateDraw(func() {
						m.mutex.RLock()
						m.textUpdate[host]["status"].Clear()
						m.textUpdate[host]["host"].
							SetText(fmt.Sprintf(status))
						m.textUpdate[host]["psu"].Clear()
						m.textUpdate[host]["board"].Clear()
						m.mutex.RUnlock()
					})
				}
			}

		}
	}
}

func (m Model) refreshData(t int) {
	for {
		//var t map[string]float64
		//t = make(map[string]float64)
		m.getData()
		for _, p := range m.power.Panels {
			for _, c := range p.Circuits {
				for _, leg := range c.Legs {
					leg.ClearLoad()
				}
			}
		}
		for _, host := range m.hosts {
			if m.siteData[host] == nil {
				continue
			}

			for _, leg := range m.hostPower[host] {
				leg.AddLoad(m.siteData[host].PsuStats.In_w / 2)
			}

			m.app.QueueUpdateDraw(func() {
				m.mutex.RLock()
				if m.siteDataError[host] == nil {
					m.textUpdate[host]["status"].
						SetText(report.Status(m.siteData[host]))
					m.textUpdate[host]["host"].
						SetText(report.Performance(m.siteData[host], m.oceanData[host]))
					m.textUpdate[host]["psu"].
						SetText(report.Psu(m.siteData[host], m.power.Panels[0].V))
					if m.boardDataError[host] == nil {
						m.textUpdate[host]["board"].
							SetText(report.Board(m.siteData[host], m.boardData[host]))
					} else {
						m.textUpdate[host]["board"].SetText(fmt.Sprint(m.boardDataError[host]))
					}
				} else {
					m.textUpdate[host]["status"].SetText(fmt.Sprint(m.siteDataError[host]))
					m.textUpdate[host]["host"].Clear()
					m.textUpdate[host]["psu"].Clear()
					m.textUpdate[host]["board"].Clear()
				}
				m.mutex.RUnlock()
			})
		}
		m.app.QueueUpdateDraw(func() {
			m.mutex.RLock()
			rpt := report.Power(m.power)
			for pi, p := range m.power.Panels {
				for ci, _ := range p.Circuits {
					cName := fmt.Sprintf("Circuit %d ", pi+ci+1)
					m.textUpdate["Power"][cName].
						SetText("[yellow]" + cName + "[-]\n" + rpt[pi+ci]).
						SetDynamicColors(true)
				}
			}

			m.textUpdate["Power"]["Total"].
				SetText("[yellow]Total[-]\n" + report.TotalPower(m.power)).
				SetDynamicColors(true)
			m.mutex.RUnlock()
		})
		time.Sleep(time.Duration(t) * time.Second)
	}
}

func (m Model) getData() {

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		oceanTable := Ocean.DumpTable(m.btcAddr)
		if oceanTable != nil {
			for _, row := range *oceanTable {
				row := row
				m.oceanData["miner"+row.Nickname] = &row
			}
		}

	}()

	for _, host := range m.hosts {
		wg.Add(1)
		host := host
		go func() {
			defer wg.Done()
			m.mutex.Lock()
			m.siteData[host], m.siteDataError[host] = ePIC.GetSummary(host)
			m.mutex.Unlock()
		}()

	}

	wg.Wait()
}

func (m Model) getLogData(t int) {
	for {
		for _, host := range m.hosts {
			m.boardData[host], m.boardDataError[host] = ePIC.GetBoards(host)
		}
		/*
			host := host
			go func() {
				m.mutex.Lock()
				if b := ePIC.GetBoards(host); b != nil {
					m.boardData[host] = b
				}
				m.mutex.Unlock()
			}()
		*/

		time.Sleep(time.Duration(t) * time.Second)

	}
}

/*
	for _, host := range m.hosts {
		if m.boardData[host] != nil {
			for i, b := range m.boardData[host] {
				fmt.Printf("%s Board %d: %+v\n", host, i, b)
			}
		}
	}
*/

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
	m.hosts = []string{"miner001", "miner002", "miner005", "miner006", "miner007", "miner008"}
	//m.hosts = []string{"miner005", "miner006", "miner007"}

	m.siteData = make(map[string]*ePIC.Summary)
	m.oceanData = make(map[string]*Ocean.UserTable)
	m.textUpdate = make(textUpdate)

	// Define power connections
	m.power = power.Init()
	m.hostPower["miner005"] = []*power.Leg{m.power.Legs[0], m.power.Legs[1]}
	m.hostPower["miner006"] = []*power.Leg{m.power.Legs[1], m.power.Legs[2]}
	m.hostPower["miner007"] = []*power.Leg{m.power.Legs[0], m.power.Legs[2]}
	m.hostPower["miner008"] = []*power.Leg{m.power.Legs[0], m.power.Legs[1]}
	m.hostPower["miner001"] = []*power.Leg{m.power.Legs[1], m.power.Legs[2]}
	m.hostPower["miner002"] = []*power.Leg{m.power.Legs[0], m.power.Legs[2]}

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
			SetLabel(label).
			SetRegions(true).
			Highlight("Th")

	}

	grid := tview.NewGrid().SetBorders(true).SetColumns(-1, -2, -1, -2)

	var i int
	for _, host := range m.hosts {
		i += 1
		if m.textUpdate[host] == nil {
			m.textUpdate[host] = make(textView)
		}

		m.textUpdate[host]["status"] = newPrimitive(host + " ")
		grid.AddItem(m.textUpdate[host]["status"], i, 0, 1, 1, 0, 0, true)

		m.textUpdate[host]["host"] = newPrimitive("")
		grid.AddItem(m.textUpdate[host]["host"], i, 1, 1, 1, 0, 0, true)

		m.textUpdate[host]["psu"] = newPrimitive("")
		grid.AddItem(m.textUpdate[host]["psu"], i, 2, 1, 1, 0, 0, true)

		m.textUpdate[host]["board"] = newPrimitive("")
		grid.AddItem(m.textUpdate[host]["board"], i, 3, 1, 1, 0, 0, true)
	}
	// Line after the last host
	i += 1
	//if m.textUpdate["Total"] == nil {
	//	m.textUpdate["Total"] = make(textView)
	//}
	//m.textUpdate["Total"]["THs"] = newPrimitive("Hashrate ")
	//grid.AddItem(m.textUpdate["Total"]["THs"], i, 2, 1, 1, 0, 0, true)

	if m.textUpdate["Power"] == nil {
		m.textUpdate["Power"] = make(textView)
	}
	m.textUpdate["Power"]["Watts"] = newPrimitive("")
	m.textUpdate["Power"]["Watts"] = newPrimitive("Circuit ")
	grid.AddItem(m.textUpdate["Power"]["Watts"], i, 1, 1, 1, 0, 0, true)

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
		for _, leg := range m.power.Legs {
			leg.ClearLoad()
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
						SetText(report.Psu(m.siteData[host]))
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
			m.textUpdate["Power"]["Watts"].
				SetText(report.Power(m.power))
			//m.textUpdate["Total"]["THs"].
			//	SetText(report.Total(t))
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

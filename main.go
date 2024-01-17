package main

import (
	"sync"
	"time"

	"github.com/jak3kaj/ePICmon/Ocean"
	"github.com/jak3kaj/ePICmon/ePIC"
	"github.com/jak3kaj/ePICmon/report"

	"github.com/rivo/tview"
)

type textView map[string]*tview.TextView
type textUpdate map[string]textView

type model struct {
	hosts      []string
	btcAddr    string
	siteData   map[string]*ePIC.Summary
	oceanData  map[string]*Ocean.UserTable
	textUpdate textUpdate
	app        *tview.Application
	mutex      *sync.RWMutex
}

func main() {
	m := new(model)
	m.hosts = []string{"miner001", "miner002", "miner005", "miner006", "miner007", "miner008"}
	m.btcAddr = "bc1qluhcxmzf8up8m8625gtl74458jemt8jcgrp3u3"
	m.siteData = make(map[string]*ePIC.Summary)
	m.oceanData = make(map[string]*Ocean.UserTable)
	m.textUpdate = make(textUpdate)
	m.mutex = &sync.RWMutex{}

	m.getData()

	newPrimitive := func(label string) *tview.TextView {
		return tview.NewTextView().
			SetTextAlign(tview.AlignLeft).
			SetLabel(label)
	}

	grid := tview.NewGrid().SetBorders(true).SetColumns(-1, -2, -2, -2)

	for i, host := range m.hosts {
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
	m.app = tview.NewApplication()
	go m.refreshData(5)
	if err := m.app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}

}

func (m model) refreshData(t int) {
	for {
		m.getData()
		for _, host := range m.hosts {
			m.app.QueueUpdateDraw(func() {
				m.mutex.RLock()
				m.textUpdate[host]["status"].
					SetText(report.Status(m.siteData[host]))
				m.textUpdate[host]["host"].
					SetText(report.Performance(m.siteData[host], m.oceanData[host]))
				m.textUpdate[host]["psu"].
					SetText(report.Psu(m.siteData[host]))
				m.textUpdate[host]["board"].
					SetText(report.Board(m.siteData[host]))
				m.mutex.RUnlock()
			})
		}
		//m.app.Sync()
		time.Sleep(time.Duration(t) * time.Second)
	}
}

func (m model) getData() {

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
			m.siteData[host] = ePIC.GetSummary(host)
			m.mutex.Unlock()
		}()

	}

	wg.Wait()}

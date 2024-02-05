package power

type Leg struct {
	ID string
	W  float64
}

type Circuit struct {
	Legs []*Leg
}

type Panel struct {
	Circuits []*Circuit
	V        float64
}

type Power struct {
	Panels []*Panel
}

func Init() *Power {
	p := Power{}
	p0 := new(Panel)
	p0.V = 212.0
	c0 := new(Circuit)
	c1 := new(Circuit)
	c0.Legs = []*Leg{&Leg{ID: "1-2"}, &Leg{ID: "2-3"}, &Leg{ID: "3-1"}}
	c1.Legs = []*Leg{&Leg{ID: "1-2"}, &Leg{ID: "2-3"}, &Leg{ID: "3-1"}}
	p0.Circuits = []*Circuit{c0, c1}
	p.Panels = []*Panel{p0}
	return &p
}

func (p *Leg) AddLoad(load float64) {
	p.W += load
}

func (p *Leg) ClearLoad() {
	p.W = 0
}

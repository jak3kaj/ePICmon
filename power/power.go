package power

type Leg struct {
	ID string
	W  float64
}

type Power struct {
	Legs []*Leg
}

func Init() *Power {
	p := Power{}
	p.Legs = []*Leg{&Leg{ID: "1-2"}, &Leg{ID: "2-3"}, &Leg{ID: "3-1"}}
	return &p
}

func (p *Leg) AddLoad(load float64) {
	p.W += load
}

func (p *Leg) ClearLoad() {
	p.W = 0
}

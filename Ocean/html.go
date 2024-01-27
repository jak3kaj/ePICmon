package Ocean

import (
	"context"
	//"net/url"
	"github.com/nfx/go-htmltable"
)

type UserTable struct {
	Nickname    string `header:"Nickname"`
	Status      string `header:"Status"`
	LastShareTS string `header:"Last Share"`
	Hashrate60s string `header:"Hashrate (60s)"`
	Hashrate3hr string `header:"Hashrate (3hr)"`
	Earnings    string `header:"Earnings"`
}

func DumpTable(user string) *[]UserTable {
	oceanUrl := "https://ocean.xyz/stats/" + user

	htmltable.Logger = func(_ context.Context, msg string, fields ...any) {
	}
	out, err := htmltable.NewSliceFromURL[UserTable](oceanUrl)
	if err != nil {
		return nil
	}
	return &out
}

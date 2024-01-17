package Ocean

import (
	"context"
	"fmt"
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
	/*
		oceanUrl := "https://ocean.xyz/template/workers/rows"
		query := url.Values{
			"user": {user},
		}
		oceanUrl += query.Encode()
	*/
	oceanUrl := "https://ocean.xyz/stats/" + user

	htmltable.Logger = func(_ context.Context, msg string, fields ...any) {
	}
	out, err := htmltable.NewSliceFromURL[UserTable](oceanUrl)
	if err != nil {
		fmt.Printf("Failed to parse table: %s\n", err)
		return nil
	}
	//fmt.Println(out)
	return &out
	/*
		resp, err := http.Get(oceanUrl)
		if err != nil {
			fmt.Printf("Failed to Get Ocean table: %s\n", err)
			return nil
		}
		if resp.StatusCode < 200 && resp.StatusCode >= 400 {
			fmt.Printf("Failed to Get data: %s\n", resp.Status)
			return nil
		}

		respData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed to read response Body: %s\n", err)
			return nil
		}
		resp.Body.Close()
	*/

}

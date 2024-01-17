package ePIC

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetSummary(host string) *Summary {
	return getSummary(host, 0)
}
func getSummary(host string, i int) *Summary {
	port := 4028
	i += 1
	sum := getData(host, port, "summary")

	if sum.Result == nil {
		return sum
	}
	//fmt.Printf("GetSummary - host: %s i:%d Result: %t Error: %s\n", host, i, *sum.Result, sum.Error)

	if i <= 5 {
		return getSummary(host, i)
	}

	return nil
}

func getData(host string, port int, endpoint string) *Summary {
	sum := Summary{}

	resp, err := http.Get(fmt.Sprintf("http://%s:%d/%s", host, port, endpoint))
	if err != nil {
		fmt.Printf("Failed to Get data: %s\n", err)
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

	if err := json.Unmarshal(respData, &sum); err != nil {
		fmt.Printf("Failed to Unmarshall JSON from response Body: %s\n", err)
		return nil
	}

	return &sum
}

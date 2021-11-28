package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"time"
)

func JianhangSearch() {
	uri, _ := url.Parse("http://127.0.0.1:8888")
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(uri),
		},
	}
	url := "http://www1.ccb.com/tran/WCCMainPlatV5"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Host", "www1.ccb.com")
	req.Header.Set("Proxy-Connection", "keep-alive")
	//req.Header.Set("Referer", "http://www1.ccb.com/cn/home/map/branchSearch.html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	query := req.URL.Query()
	query.Add("CCB_IBSVersion", "V5")
	query.Add("SERVLET_NAME", "WCCMainPlatV5")
	query.Add("isAjaxRequest", "true")
	query.Add("TXCODE", "NZX010")
	query.Add("ADiv_Cd", "310000")
	query.Add("Kywd_List_Cntnt", "")
	query.Add("Enqr_MtdCd", "4")
	query.Add("PAGE", "1")
	query.Add("Cur_StCd", "4")
	cookie := http.Cookie{Name: "tranCCBIBS1", Value: "PTv%2C31lqRX9yr2CuQXbxiCKdMUZRICtdFkS1XCCeW0sdMBvpI3KK2LOnzkeV4HKdlEhx6Boe5E0pzKgen0odqGUdiEltGQBdVg", Expires: time.Now().Add(365 * 24 * time.Hour)}
	req.URL.RawQuery = query.Encode()
	req.AddCookie(&cookie)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	//doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		fmt.Println("decode error", err)
	}
	bank_list := data["OUTLET_DTL_LIST"].([]interface{})
	for _, bank := range bank_list {
		bank_data := bank.(map[string]interface{})
		fmt.Println(bank_data["CCBIns_Nm"])
		fmt.Println(bank_data["Dtl_Adr"])
		fmt.Println(bank_data["Fix_TelNo"])
	}
}

func main() {
	JianhangSearch()
}

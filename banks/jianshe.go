package banks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Set query format both get cookie and use cookie.
func AddQuery(query *url.Values, cityNum string, getCookie bool) {
	if !getCookie {
		query.Set("TXCODE", "NZX010")
		query.Set("ADiv_Cd", cityNum)
		query.Set("Kywd_List_Cntnt", "")
		query.Set("Enqr_MtdCd", "1")
		query.Set("PAGE", "1")
		query.Set("Cur_StCd", "4")
	} else {
		query.Set("CCB_IBSVersion", "V5")
		query.Set("SERVLET_NAME", "WCCMainPlatV5")
		query.Set("isAjaxRequest", "true")
		query.Set("TXCODE", "100119")
	}
}

func CityConvert(cityName string) string {
	cityMap := make(map[string]string)
	cityMap["上海"] = "310000"
	cityMap["北京"] = "110000"
	cityMap["天津"] = "120000"
	cityMap["石家庄"] = "130100"
	cityMap["太原"] = "140100"
	cityMap["呼和浩特"] = "150100"
	cityMap["沈阳"] = "210100"
	cityMap["长春"] = "220100"
	cityMap["哈尔滨"] = "230100"
	cityMap["南京"] = "320100"
	cityMap["杭州"] = "330100"
	cityMap["合肥"] = "340100"
	cityMap["福州"] = "350100"
	cityMap["南昌"] = "360100"
	cityMap["济南"] = "370100"
	cityMap["郑州"] = "410100"
	cityMap["武汉"] = "420100"
	cityMap["长沙"] = "430100"
	cityMap["广州"] = "440100"
	cityMap["南宁"] = "450100"
	cityMap["海口"] = "460100"
	cityMap["重庆"] = "500000"
	cityMap["成都"] = "510100"
	cityMap["贵阳"] = "520100"
	cityMap["昆明"] = "530100"
	cityMap["拉萨"] = "540100"
	cityMap["西安"] = "610100"
	cityMap["兰州"] = "620100"
	cityMap["西宁"] = "630100"
	cityMap["银川"] = "640100"
	cityMap["乌鲁木齐"] = "650100"
	return cityMap[cityName]
}

func JianhangSearch(cityName string) {
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
	req.Header.Set("Referer", "http://www1.ccb.com/cn/home/map/branchSearch.html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	query := req.URL.Query()
	cityNum := CityConvert(cityName)
	AddQuery(&query, cityNum, true)
	req.URL.RawQuery = query.Encode()
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Body)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("goquery error", err)
	}
	cookie := http.Cookie{Name: "tranCCBIBS1", Value: doc.Text(), Expires: time.Now().Add(365 * 24 * time.Hour)}
	AddQuery(&query, cityNum, false)
	req.URL.RawQuery = query.Encode()
	req.AddCookie(&cookie)
	res, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	if err != nil {
		log.Fatal(err)
	}
	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		fmt.Println("decode error", err)
	}
	// get bank list on page 1.
	bank_list := data["OUTLET_DTL_LIST"].([]interface{})
	for _, bank := range bank_list {
		bank_data := bank.(map[string]interface{})
		fmt.Println(bank_data["CCBIns_Nm"])
		fmt.Println(bank_data["Dtl_Adr"])
		fmt.Println(bank_data["Fix_TelNo"])
	}
	// get total page number.
	total_page, _ := strconv.Atoi(data["TOTAL_PAGE"].(string))
	for i := 2; i <= total_page; i++ {
		query.Set("PAGE", strconv.Itoa(i))
		fmt.Print(query.Encode())
		req.URL.RawQuery = query.Encode()
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		var data map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			fmt.Println("decode error", err)
		}
		bank_list = data["OUTLET_DTL_LIST"].([]interface{})
		sli := bank_list[:len(bank_list)-1]
		for _, bank := range sli {
			bank_data := bank.(map[string]interface{})
			fmt.Println(bank_data["CCBIns_Nm"])
			fmt.Println(bank_data["Dtl_Adr"])
			fmt.Println(bank_data["Fix_TelNo"])
		}
	}
}

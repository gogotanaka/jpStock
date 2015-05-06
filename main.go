package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Basis struct {
	title              string
	market             string
	industry           string
	price_limit        string
	capitalization     string
	shares_outstanding string
	dividend_yield     string
	minimum_purchase   string
	share_unit         string
	yearly_high        string
	yearly_low         string
}

type Margin struct {
	margin_buying    string
	d_margin_buying  string
	margin_rate      string
	margin_selling   string
	d_margin_selling string
}

type Index struct {
	price          string
	previousprice  string
	opening        string
	high           string
	low            string
	turnover       string
	trading_volume string
	dps            string
	per            string
	pbr            string
	eps            string
	bps            string
}

func GetPage(url string) {
	var b *Basis = &Basis{}
	var i *Index = &Index{}
	var m *Margin = &Margin{}

	doc, _ := goquery.NewDocument(url)

	textMap := func(q string, vars ...*string) {
		for i, v := range vars {
			*v = doc.Find(q).Eq(i).Text()
		}
	}

	b.title = doc.Find("table.stocksTable th.symbol h1").Text()
	b.market = doc.Find("div#ddMarketSelect span.stockMainTabName").Text()
	b.industry = doc.Find("div.stocksDtl dd.category a").Text()
	i.price = doc.Find("table.stocksTable td.stoksPrice").Last().Text()

	textMap("div.innerDate dd strong", &i.previousprice, &i.opening, &i.high, &i.low, &i.turnover, &i.trading_volume, &b.price_limit)
	textMap("div.ymuiDotLine div.yjMS dd.ymuiEditLink strong", &m.margin_buying, &m.d_margin_buying, &m.margin_rate, &m.margin_selling, &m.d_margin_selling)
	textMap("div#main div.main2colR div.chartFinance div.lineFi dl dd strong", &b.capitalization, &b.shares_outstanding, &b.dividend_yield, &i.dps, &i.per, &i.pbr, &i.eps, &i.bps, &b.minimum_purchase, &b.share_unit, &b.yearly_high, &b.yearly_low)

	fmt.Println(*b)
	fmt.Println(*i)
	fmt.Println(*m)
}

func main() {
	dat, err := ioutil.ReadFile("./codes")
	check(err)

	var limit = make(chan int, 10)
	for _, i := range strings.Split(string(dat), "\n") {
		go func(i string) {
			limit <- 1
			var url string = fmt.Sprintf("http://stocks.finance.yahoo.co.jp/stocks/detail/?code=%s", i)
			GetPage(url)
			log.Println("Finish: ", url)

			<-limit
		}(i)
	}
	select {}
}

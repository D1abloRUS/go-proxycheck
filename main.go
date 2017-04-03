package main

import (
	"flag"
	"fmt"
	"go-proxycheck/config"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang"
	"github.com/parnurzeal/gorequest"
)

type getWithproxy struct {
	proxy     string
	url       string
	fileout   string
	newstring string
	info      bool
}

func (g *getWithproxy) getproxy() {
	split := strings.Split(g.proxy, ":")
	ip := fmt.Sprintf("%s%%", split[0])
	//запрос к бд

	//

	if existStr == false {
		httpProxy := fmt.Sprintf("http://%s", g.proxy)
		request := gorequest.New().Proxy(httpProxy).Timeout(2 * time.Second)
		timeStart := time.Now()
		_, _, err := request.Get(g.url).End()
		if err != nil {
			//можно поискать удаление строки из файла, но хз
			fmt.Println("BAD: ", g.proxy)
		} else {
			fmt.Println("GOOD: ", g.proxy)
			//скорее всего убрать условие т.к все будет в бд
			if g.info == true {
				country := ipToCountry(ip)
				respone := time.Since(timeStart)
				g.newstring = fmt.Sprintf("%s;%s;%s\n", g.proxy, country, respone)
			} else {
				//должно стать не актуально
				g.newstring = fmt.Sprintf("%s\n", g.proxy)
			}
			//как и это, хотя тут скорее всего будет запрос к бд, а выше что нито типа db.prepare

		}
	}
}

func ipToCountry(ip string) string {
	db, err := geoip2.Open("/usr/share/GeoIP/GeoLite2-Country.mmdb")
	if err != nil {
		fmt.Printf("Could not open GeoIP database\n")
		os.Exit(1)
	}
	defer db.Close()
	country, _ := db.Country(net.ParseIP(ip))
	return country.Country.IsoCode
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	var (
		url    = flag.String("url", "https://m.vk.com", "")
		fileIn = flag.String("in", "proxylist.txt", "full path to proxy file")
		//удалить при работе с бд
		fileOut = flag.String("out", "goodlist.txt", "full path to output file")
		info    = flag.Bool("info", false, "info about proxy: Country, Respone")
		treds   = flag.Int("treds", 50, "number of treds")
	)

	db, err := config.NewDB("postgres://proxy:proxy@localhost/proxy")
	if err != nil {
		log.Panic(err)
	}

	env := &config.Env{DB: db}

	flag.Parse()

	//так же не понадобится
	content, _ := ioutil.ReadFile(*fileIn)
	proxys := strings.Split(string(content), "\n")

	workers := *treds

	wg := new(sync.WaitGroup)
	in := make(chan string, 2*workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for proxy := range in {
				gp := getWithproxy{
					proxy:   proxy,
					url:     *url,
					fileout: *fileOut,
					info:    *info,
				}
				gp.getproxy()
			}
		}()
	}

	for _, proxy := range proxys {
		if proxy != "" {
			in <- proxy
		}
	}
	close(in)
	wg.Wait()
}
package main

import (
	"flag"
	"log"

	"github.com/JimYJ/scanproxy/config"

	"github.com/JimYJ/scanproxy/scanproxy"
)

var (
	ipStep        = 10
	area          string
	ipstep        int
	mode          bool
	maxConcurrent int
)

func main() {
	mysql := config.MySQL()
	defer mysql.Close()
	if mode {
		log.Println("now work in fast mode,ip scan step:", ipstep, "area:", area, "maxConcurrent:", maxConcurrent)
		scanproxy.SetQueueMaxConcurrent(maxConcurrent)
		scanproxy.InternetFastScan(area, ipstep)
	} else {
		log.Println("now work in normal mode,ip scan step:", ipstep, "area:", area)
		scanproxy.InternetAllScan(area, ipstep)
	}
}

func init() {
	flag.BoolVar(&mode, "f", false, "scan mode, fast is only scan common port, default is scan all port(shorthand)")
	flag.IntVar(&ipstep, "i", 10, "set scan how many IP segment in same times, it will affect memory footprint(shorthand)")
	flag.IntVar(&maxConcurrent, "m", 200, "maximum concurrency number(shorthand)")
	flag.StringVar(&area, "a", "CN", "country codes, see ISO 3166-1(shorthand)")
	flag.Parse()
}

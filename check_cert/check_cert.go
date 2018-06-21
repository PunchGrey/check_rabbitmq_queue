package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	var (
		host       = flag.String("host", "gptl.ru:443", "The host with port for checked certificate, host:port")
		warnYears  = flag.Int("wy", 0, "Warn if the certificate will expire within this many years.")
		errYears   = flag.Int("ey", 0, "Error if the certificate will expire within this many years.")
		warnMonths = flag.Int("wm", 0, "Warn if the certificate will expire within this many months.")
		errMonths  = flag.Int("ym", 0, "Error if the certificate will expire within this many months.")
		warnDays   = flag.Int("wd", 25, "Warn if the certificate will expire within this many days.")
		errDays    = flag.Int("ed", 15, "Error if the certificate will expire within this many days.")
	)
	var infoStr string
	var infoCert []string

	flag.Parse()
	conn, err := tls.Dial("tcp", *host, nil)
	defer conn.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error whith connect %v", err)
		os.Exit(3)
	}

	for _, chain := range conn.ConnectionState().VerifiedChains {
		for certNum, cert := range chain {
			if time.Now().AddDate(*errYears, *errMonths, *errDays).After(cert.NotAfter) {
				fmt.Println("validity:", cert.NotAfter, certNum)
				os.Exit(2)
			}
			if time.Now().AddDate(*warnYears, *warnMonths, *warnDays).After(cert.NotAfter) {
				fmt.Println("validity:", cert.NotAfter, certNum)
				os.Exit(1)
			}
			infoStr = fmt.Sprintf("validity: %v %d", cert.NotAfter, certNum)
			infoCert = append(infoCert, infoStr)
		}
	}
	fmt.Println(infoCert)
	os.Exit(0)
}

package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
	"time"
)

func main() {
	response, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)

	}
	currentTime := time.Now().Add(response.ClockOffset)
	fmt.Println("Текущее время:", currentTime.Format("2006-01-02 15:04:05 MST"))

}

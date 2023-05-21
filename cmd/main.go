package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fachebot/a3spow"
)

func main() {
	cfg := a3spow.MustReadConfig("config.yaml")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	waitChan := make(chan bool)
	outAddressChan := make(chan a3spow.OutAddress, 1024)
	addresses := make([]a3spow.OutAddress, 0)
	go func() {
		for address := range outAddressChan {
			addresses = append(addresses, address)
		}
		close(waitChan)
	}()

	startTime := time.Now()
	a3spow.StartMining(ctx, cfg.Owner, &cfg.Filter, cfg.Number, outAddressChan)

	close(outAddressChan)
	<-waitChan

	fmt.Println()
	data := a3spow.RenderTable(addresses)
	fmt.Println(string(data))
	fmt.Println("Total elapsed time:", time.Since(startTime))

	filename := fmt.Sprintf("address-%d.txt", time.Now().Unix())
	if err := os.WriteFile(filename, data, 0666); err == nil {
		fmt.Printf(`The address has been saved to the file "%s"`, filename)
	}
}

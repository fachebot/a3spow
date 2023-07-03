package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fachebot/a3spow"
)

func waitForQuitSignals(callback func(os.Signal)) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		sig := <-sigs
		callback(sig)
	}()
}

func main() {
	cfg := a3spow.MustReadConfig("config.yaml")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	waitChan := make(chan bool)
	outAddressChan := make(chan a3spow.OutAddress, 128)
	addresses := make([]a3spow.OutAddress, 0)
	go func() {
		for address := range outAddressChan {
			addresses = append(addresses, address)
		}
		close(waitChan)
	}()

	startTime := time.Now()
	exitHandle := func(sig os.Signal) {
		close(outAddressChan)
		<-waitChan

		fmt.Println()
		fmt.Println("Total elapsed time:", time.Since(startTime))

		if len(addresses) > 0 {
			data := a3spow.RenderTable(addresses)
			fmt.Println(string(data))

			filename := fmt.Sprintf("address-%d.txt", time.Now().Unix())
			if err := os.WriteFile(filename, data, 0666); err == nil {
				fmt.Printf(`The address has been saved to the file "%s"`, filename)
			}
		}

		os.Exit(1)
	}
	waitForQuitSignals(exitHandle)
	if !cfg.HD {
		a3spow.StartMining(ctx, cfg.Owner, &cfg.Filter, cfg.Number, outAddressChan)
	} else {
		a3spow.StartHdMining(ctx, cfg.Owner, &cfg.Filter, cfg.Number, outAddressChan)
	}

	if ctx.Err() == nil {
		exitHandle(syscall.Signal(0x0))
	}
}

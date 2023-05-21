package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fachebot/a3spow"
)

func main() {
	filter := &a3spow.LongRepeatedFilter{MinLength: 8}
	owner := common.HexToAddress("0x999999273b1f52e3243f526dd54c974b46cd4f05")

	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan a3spow.OutAddress, 1024)
	addresses := make([]a3spow.OutAddress, 0)
	go func() {
		for address := range ch {
			fmt.Println(address.Address)
			addresses = append(addresses, address)
			cancel()
		}
	}()

	startTime := time.Now()
	a3spow.StartMining(ctx, owner, filter, 100000000, ch)
	fmt.Println(time.Since(startTime))
}

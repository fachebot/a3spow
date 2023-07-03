package a3spow

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

func GetHdAddress(owner common.Address) (string, common.Address) {
	mnemonic, err := hdwallet.NewMnemonic(128)
	if err != nil {
		panic(err)
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		panic(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		panic(err)
	}

	return mnemonic, account.Address
}

func StartHdMining(ctx context.Context, owner common.Address, filter Filter, count int32, outChan chan<- OutAddress) {
	var total uint64
	var validAddressCount int32
	var waitGroup sync.WaitGroup
	threads := int(float64(runtime.NumCPU()) * 1.5)
	for i := 0; i < threads; i++ {
		waitGroup.Add(1)

		go func() {
			for {
				if ctx.Err() != nil {
					break
				}

				salt, address := GetHdAddress(owner)
				if filter.Filter(address.Hex()) {
					outChan <- OutAddress{
						Address: address,
						Salt:    salt,
					}
					atomic.AddInt32(&validAddressCount, 1)
				}

				atomic.AddUint64(&total, 1)

				if atomic.LoadInt32(&validAddressCount) >= count {
					break
				}
			}

			waitGroup.Done()
		}()
	}

	var lastTotal uint64
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			currTotal := atomic.LoadUint64(&total)
			fmt.Printf("[%s] %d addr/s | Generate: %d addr | Minted: %d addr \n",
				time.Now().Format(time.RFC3339), currTotal-lastTotal, currTotal, atomic.LoadInt32(&validAddressCount))
			lastTotal = currTotal
		}
	}()

	waitGroup.Wait()
	ticker.Stop()
}

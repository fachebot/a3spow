package a3spow

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type OutAddress struct {
	Salt    string
	Address common.Address
}

func RandomSalt() [32]byte {
	var salt [32]byte
	for i := 0; i < len(salt)/8; i++ {
		binary.BigEndian.PutUint64(salt[i*8:], rand.Uint64())
	}

	return salt
}

func GetAddress(owner common.Address, salt [32]byte) common.Address {
	mutantSalt := solsha3.SoliditySHA3(
		[]string{"address", "bytes32"},
		[]interface{}{owner, salt},
	)

	var buf [32]byte
	copy(buf[:], mutantSalt)
	return crypto.CreateAddress2(FactoryAddress, buf, WalletByteHash)
}

func StartMining(ctx context.Context, owner common.Address, filter Filter, count int32, outChan chan<- OutAddress) {
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

				salt := RandomSalt()
				address := GetAddress(owner, salt)
				if filter.Filter(address.Hex()) {
					outChan <- OutAddress{
						Address: address,
						Salt:    hex.EncodeToString(salt[:]),
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

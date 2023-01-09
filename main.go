package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/sync/errgroup"
)

var (
	RPC = "ws://127.0.0.1:8545"
	TICK = 2
)

type txPoolStatus struct {
	BaseFee string `json:"baseFee"`
	Pending string `json:"pending"`
	Queued  string `json:"queued"`
}

type TxPoolStatusFormated struct {
	BaseFee 	uint `json:"baseFee"`
	Pending 	uint `json:"pending"`
	Queued  	uint `json:"queued"`
}


func txpoolStatus() (txPoolStatus txPoolStatus, err error) {
	url := RPC
	rpcClient, err := rpc.DialContext(context.Background(), url)
	if err != nil {
		log.Fatalln(err)
	}	
	err = rpcClient.Call(&txPoolStatus, "txpool_statu")
	if err != nil {
		log.Fatalln(err)
		return txPoolStatus, err
	}
	return txPoolStatus, nil
}

// Gather the "txpool_status" response
func PoolStatus(TxPoolStatusChan chan TxPoolStatusFormated) {
	ctx, cancel 	:= context.WithCancel(context.Background())
	g, ctx 			:= errgroup.WithContext(ctx)
	ticker 			:= time.NewTicker(time.Duration(TICK) * time.Second)
	defer ticker.Stop()

	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				txpool, err := txpoolStatus()
				if err != nil {
					return err
				}
				// faster than strconv.ParseUint
				baseFee, _ 	:= strToUint(txpool.BaseFee, 0, 32)
				pending, _ 	:= strToUint(txpool.Pending, 0, 32)
				queued, _ 	:= strToUint(txpool.Queued, 0, 32)
				if pending > 0 {
					TxPoolStatusChan <- TxPoolStatusFormated{
						BaseFee: baseFee,
						Pending: pending,
						Queued: queued,
					}
				}
			}
		}
	})

	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
	cancel()
}

func strToUint(s string, base int, bitSize int) (uint, error) {
	if base == 0 {
		if len(s) > 2 && s[0:2] == "0x" {
			base = 16
			s = s[2:]
		} else {
			base = 10
		}
	}
	if base == 16 {
		i, err := strconv.ParseUint(s, base, bitSize)
		return uint(i), err
	} else if base == 32 {
		n := 0
		values := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, -1, -1, -1, -1, -1, -1, -1, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
		for i := 0; i < len(s); i++ {
			c := s[i]
			if c >= '0' && c <= '9' {
				n = n*32 + values[c-'0']
			} else if c >= 'A' && c <= 'Z' {
				n = n*32 + values[c-'A'+10]
			} else {
				return 0, fmt.Errorf("invalid character: %c", c)
			}
		}
		return uint(n), nil
	}
	return 0, fmt.Errorf("unsupported base: %d", base)
}

func main() {
	errChan 			:= make(chan error)
	TxPoolStatusChan 	:= make(chan TxPoolStatusFormated)

	go PoolStatus(TxPoolStatusChan)

	for {
		select {
		case errs := <-errChan:
			log.Println("[ERROR]", errs)
		case txpool 	:= <-TxPoolStatusChan:
			fmt.Printf("[MONIT] baseFee %d\n", txpool.BaseFee)
			fmt.Printf("[MONIT] pending %d\n", txpool.Pending)
			fmt.Printf("[MONIT] queued %d\n\n", txpool.Queued) 
		}
	}
}
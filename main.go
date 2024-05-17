package main

import (
	"fmt"
)

type DsNode struct {
	dsId uint32
}

func (ds *DsNode) DsGetSegments() (sids []uint64, err error) {
	// 暂时删除全部
	// segment ID生成算法: DS ID 左移 32位
	var (
		dsCapacity uint64 = 257
	)
	if dsCapacity < 256 {
		sids = append(sids, uint64((*ds).dsId)<<32)
	} else if dsCapacity%256 != 0 {
		sids = append(sids, uint64((*ds).dsId)<<32+dsCapacity/256)
	}
	for sid, i := uint64((*ds).dsId)<<32, uint64(0); i < dsCapacity/256; i++ {
		sids = append(sids, sid+i)
	}
	return sids, err
}

func main() {
	// stripeWidth := 3
	// stripeSize := 524288
	// stripeSize = stripeSize << 9
	// segSize := 256 * 1024 * 1024
	// newCap := 256 + 20
	// segStripeNum := ((newCap << 20) % (segSize * stripeWidth))
	// fmt.Printf("segStripeNum: %d\n", segStripeNum)
	// segStripeNum = ((newCap << 20) % (segSize * stripeWidth)) / stripeSize
	// fmt.Printf("segStripeNum: %d\n", segStripeNum)
	dataLen := 128
	sids := make([]struct{}, 256+63)
	iterations := len(sids) / dataLen
	if len(sids)%dataLen != 0 {
		iterations += 1
	}
	for i := 0; i < iterations; i++ {
		start := i * dataLen
		end := start + dataLen
		if end > len(sids) {
			end = len(sids)
		}
		fmt.Printf("ITERATIONS: %d, process len:%d\n", i, end-start)
		for j := start; j < end; j++ {
			fmt.Printf("i:%v, start:%v, end:%v, j:%v\n", i, start, end, j)
		}
	}
}

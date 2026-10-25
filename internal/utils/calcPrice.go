package utils

import (
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

func CalcTotalPrice(startTime time.Time, endTime time.Time) (int64, error) {
	var totalPrice int64
	singlePrice, err := strconv.Atoi(os.Getenv("price"))
	if err != nil {
		log.Printf("在计算总价时出现错误:%s", err.Error())
		return 0, err
	}
	// 每半小时结一次单价
	duration := math.Ceil(endTime.Sub(startTime).Seconds())
	log.Printf("开始于:%s,结束于:%s,经过了:%f", startTime, endTime, duration)
	d := math.Ceil(duration/60/60) * 2
	totalPrice = int64(d) * int64(singlePrice)
	return totalPrice, nil
}

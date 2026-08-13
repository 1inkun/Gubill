package utils

func CalculateTotalPrice(duration int64, singlePrice int64) int {
	// 先处理duration，单个计价阶梯最大计费时长为24小时，半小时为一个计费单位，超过24小时的部分按阶梯另外计费
	var thisCycleDuration int64
	var nextCycleDuration int64
	var totalPrice int
	if duration <= 24*60*60 {
		thisCycleDuration = duration
		nextCycleDuration = 0
	} else {
		thisCycleDuration = 24 * 60 * 60
		nextCycleDuration = duration - thisCycleDuration
	}
	token := thisCycleDuration / (30 * 60)
	if thisCycleDuration%(30*60) != 0 {
		token++
	}
	// 本周期总价按阶梯计费，12小时为分界，12小时内最大收费额价格为 单价 x 10，24小时内最大收费价格为 单价 x 20
	if token <= 24 {
		totalPrice = int(token * singlePrice)
		if totalPrice > int(10*singlePrice) {
			totalPrice = int(10 * singlePrice)
		}
	} else {
		token -= 24
		totalPrice = int(token*singlePrice) + int(10*singlePrice)
		if totalPrice > int(20*singlePrice) {
			totalPrice = int(20 * singlePrice)
		}
	}
	// 如果下周期价格不为0,加上下周期的价格
	if nextCycleDuration != 0 {
		totalPrice += CalculateTotalPrice(nextCycleDuration, singlePrice)
	}
	return totalPrice
}

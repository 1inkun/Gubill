const CalcTotalValue = function(duration:number,singlePrice:number):number{
    // 先处理duration，单个计价阶梯最大计费时长为24小时，半小时为一个计费单位，超过24小时的部分按阶梯另外计费
    let thisCycleDuration:number
    let nextCycleDuration:number
    let totalPrice:number = 0
    if (duration <= 24*60*60) {
		thisCycleDuration = duration
		nextCycleDuration = 0
	} else {
		thisCycleDuration = 24 * 60 * 60
		nextCycleDuration = duration - thisCycleDuration
	}
	let token = Math.floor(thisCycleDuration / (30 * 60)) 
	if (thisCycleDuration%(30*60) != 0) {
		token++
	}
    // console.log(token)
    // 本周期总价按阶梯计费，12小时为分界，12小时内最大收费额价格为 单价 x 10，24小时内最大收费价格为 单价 x 20
	if (token <= 24) {
		totalPrice = (token * singlePrice)
		if (totalPrice > (10*singlePrice) ){
			totalPrice = (10 * singlePrice)
		}
	} else {
		token -= 24
		totalPrice = (token*singlePrice) + (10*singlePrice)
		if (totalPrice > (20*singlePrice)){
			totalPrice = (20 * singlePrice)
		}
	}
	// 如果下周期价格不为0,加上下周期的价格
	if (nextCycleDuration != 0) {
		totalPrice += CalcTotalValue (nextCycleDuration, singlePrice)
	}
    return totalPrice
}

export default CalcTotalValue
export function ShowErrMsg(msg:string) {
    uni.showToast({
        title: msg,
        icon: 'error',
        mask: true
    })
}

export function ShowSuccessMsg(msg:string){
    uni.showToast({
        title: msg,
        icon: 'success',
        mask: true
    })
}


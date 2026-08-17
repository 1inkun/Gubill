import { jwtDecode, JwtPayload } from "jwt-decode";

const CheckLoginStatus = function(tokenStr: string):boolean{
    const now = Math.ceil(Date.now() / 1000) 
    if (tokenStr == ""){
        return false
    }
    const decode = jwtDecode<JwtPayload>(tokenStr)
    if(decode.exp == undefined){
        return false
    }
    if(decode.exp < now){
        return false
    }
    return true
}
export default CheckLoginStatus
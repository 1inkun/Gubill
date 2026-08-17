type LoginData = {
    username: string,
    password: string
}

type LoginRes = {
    token: string
}

type RegisterData = {
    username: string
    nickname: string
    password: string
    email: string
}

type RegisterRes = {
    userId: string
}
export type {LoginData,LoginRes,RegisterData,RegisterRes}
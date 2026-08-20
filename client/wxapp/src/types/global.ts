type ResponseData<T> = {
	code: number;
	status: string;
	msg: string;
	data: T | Array<T>;
};
type PaginData<T> = {
	results: Array<T>,
	page: number,
	pageSize: number
}
type Pagin = {
	page: number,
	pageSize: number
}
type UserInfo = {
	userId: string;
	username: string;
	nickname: string;
	role: string;
};
type CustomJWTClaims = {
	username: string,
	userId: string,
	nickname: string,
	role: string,
	exp: number
}
export type { ResponseData, UserInfo, CustomJWTClaims, Pagin, PaginData };

type LoginData = {
	username: string;
	password: string;
};

type LoginRes = {
	token: string;
};

type RegisterData = {
	username: string;
	nickname: string;
	password: string;
	email: string;
};

type RegisterRes = {
	userId: string;
};

type SignData = {
	signId: string;
	userId: string;
	start_at: number;
	end_at: number;
	status: number;
	value: number;
};

type GenSignRes = {
    signId: string
}

type FinishSignRes = {
    payId: string
}
export type { LoginData, LoginRes, RegisterData, RegisterRes, SignData, GenSignRes, FinishSignRes};

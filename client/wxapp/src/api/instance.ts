import { un, UnInstance } from "@uni-helper/uni-network";

const NewInstance = function (baseUrl: string, tokenStr: string): UnInstance {
	const instance = un.create({
		baseUrl: baseUrl,
		timeout: 1000,
		headers: { Authorization: tokenStr },
	});
	return instance;
};

export default NewInstance;

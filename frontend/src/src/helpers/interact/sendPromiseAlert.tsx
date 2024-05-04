import toast from "react-hot-toast";

export const sendPromiseAlert = async <T,>(
	promiseFn: Promise<T>,
	loading: string,
	success: string,
	error: string,
	hideError = true
) => {
	if (!hideError) {
		return await toast.promise(promiseFn, {
			loading,
			success,
			error
		});
	}

	try {
		return await toast.promise(promiseFn, {
			loading,
			success,
			error
		});
	} catch {
		/* empty */
	}
};

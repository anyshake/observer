import toast from 'react-hot-toast';

export const sendPromiseAlert = async <T,>(
    promiseFn: Promise<T>,
    loading: string,
    success: string | ((data: T) => string),
    error: string | ((error: unknown) => string),
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

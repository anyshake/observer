import toast from 'react-hot-toast';

import { isRequestCanceled } from '../request/isRequestCanceled';

export const sendPromiseAlert = async <T,>(
    promiseFn: Promise<T>,
    loading: string,
    success: string | ((data: T) => string),
    error: string | ((error: unknown) => string),
    hideError = true
) => {
    const id = toast.loading(loading);
    try {
        const data = await promiseFn;
        toast.success(typeof success === 'function' ? success(data) : success, { id });
        return data;
    } catch (reason) {
        if (isRequestCanceled(reason)) {
            toast.dismiss(id);
        } else {
            toast.error(typeof error === 'function' ? error(reason) : error, { id });
        }
        if (!hideError) {
            throw reason;
        }
    }
};

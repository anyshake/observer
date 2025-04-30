import toast from 'react-hot-toast';

export const sendUserAlert = (message: string, error = false, duration = 2000) => {
    if (error) {
        toast.error(message, { duration });
    } else {
        toast.success(message, { duration });
    }
};

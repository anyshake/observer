import toast from "react-hot-toast";

export const sendUserAlert = (
    message: string,
    error: boolean = false,
    duration: number = 2000
) => {
    if (error) {
        toast.error(message, { duration });
    } else {
        toast.success(message, { duration });
    }
};

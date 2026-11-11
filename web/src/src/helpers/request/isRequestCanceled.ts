import axios from 'axios';

export const isRequestCanceled = (error: unknown): boolean => {
    if (axios.isCancel(error) || (error instanceof Error && error.name === 'AbortError')) {
        return true;
    }
    return (
        error !== null &&
        typeof error === 'object' &&
        'networkError' in error &&
        isRequestCanceled(error.networkError)
    );
};

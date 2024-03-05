import ErrorIcon from "../assets/icons/circle-exclamation-solid.svg";

const NotFound = () => {
    return (
        <div className="w-full min-h-screen flex flex-col items-center justify-center">
            <div className="animate-bounce">
                <img className="py-2 w-20 h-20" src={ErrorIcon} alt="" />
            </div>

            <h2 className="py-2 text-2xl font-bold text-gray-600">
                404 Not Found
            </h2>
        </div>
    );
};

export default NotFound;

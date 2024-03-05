interface ButtonProps {
    readonly label: string;
    readonly disabled?: boolean;
    readonly className?: string;
    readonly onClick?: () => void;
}

export const Button = (props: ButtonProps) => {
    const { className, label, onClick } = props;

    return (
        <button
            className={`w-full text-white font-medium text-sm shadow-lg rounded-lg py-2 ${
                className ?? ""
            }`}
            onClick={onClick}
        >
            {label}
        </button>
    );
};

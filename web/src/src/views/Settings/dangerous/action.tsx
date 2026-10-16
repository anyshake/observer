import type { MouseEventHandler } from 'react';

interface IDangerousAction {
    readonly title: string;
    readonly description: string;
    readonly buttonText: string;
    readonly onClick: MouseEventHandler<HTMLButtonElement>;
    readonly value?: string | number;
    readonly disabled?: boolean;
    readonly activeLabel?: string;
}

export const DangerousAction = ({
    title,
    description,
    buttonText,
    onClick,
    value,
    disabled,
    activeLabel
}: IDangerousAction) => (
    <div
        className={`flex flex-col justify-between gap-4 border-b border-gray-300 p-4 transition-colors sm:flex-row sm:items-center ${
            activeLabel ? 'bg-amber-50' : ''
        }`}
    >
        <div className="min-w-0 space-y-2">
            <div className="flex flex-wrap items-center gap-2">
                <h3 className="text-lg font-semibold text-gray-800">{title}</h3>
                {activeLabel && (
                    <span className="rounded-full bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-800">
                        {activeLabel}
                    </span>
                )}
            </div>
            <p className="text-sm text-gray-500">{description}</p>
        </div>
        <button
            type="button"
            value={value}
            className="btn w-full shrink-0 rounded-lg bg-red-500 px-4 py-2 font-medium text-white shadow-lg transition-all hover:bg-red-700 disabled:cursor-not-allowed disabled:bg-gray-400 sm:w-auto"
            disabled={disabled}
            onClick={onClick}
        >
            {buttonText}
        </button>
    </div>
);

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
    <div className={`app-dangerous-action ${activeLabel ? 'app-dangerous-action-active' : ''}`}>
        <div className="app-action-description">
            <div className="app-toolbar">
                <h3 className="app-subheading">{title}</h3>
                {activeLabel && <span className="app-action-status">{activeLabel}</span>}
            </div>
            <p className="app-description">{description}</p>
        </div>
        <button
            type="button"
            value={value}
            className="btn app-dangerous-action-button"
            disabled={disabled}
            onClick={onClick}
        >
            {buttonText}
        </button>
    </div>
);

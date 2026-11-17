import type { ReactNode } from 'react';

export type PageHeaderStatus = 'ok' | 'warning' | 'error';

export interface IPageHeaderNotice {
    readonly status: PageHeaderStatus;
    readonly message: ReactNode;
}

interface IPageHeader {
    readonly title: ReactNode;
    readonly meta?: ReactNode;
    readonly status: PageHeaderStatus;
    readonly statusText: ReactNode;
    readonly detail?: ReactNode;
    // Extra messages (e.g. upgrade hints, errors) shown below the header row
    readonly notices?: readonly IPageHeaderNotice[];
}

const statusClass: Record<PageHeaderStatus, string> = {
    ok: 'status-success',
    warning: 'status-warning',
    error: 'status-error'
};

// Neutral page heading card: title and meta on the left, live status on the right.
export const PageHeader = ({ title, meta, status, statusText, detail, notices }: IPageHeader) => (
    <section className="card app-panel app-page-header">
        <div className="app-page-header-row">
            <div className="app-page-header-main">
                <h1 className="app-page-header-title">{title}</h1>
                {meta && <p className="app-page-header-meta">{meta}</p>}
            </div>
            <div className="app-page-header-side" role="status">
                <span className="app-page-header-state">
                    <span
                        className={`status ${statusClass[status]} ${status === 'ok' ? '' : 'app-bounce'}`}
                        aria-hidden
                    />
                    {statusText}
                </span>
                {detail && <span className="app-page-header-meta">{detail}</span>}
            </div>
        </div>
        {notices && notices.length > 0 && (
            <ul className="app-page-header-notices">
                {notices.map(({ status: level, message }, index) => (
                    <li key={index} className="app-page-header-notice" role="alert">
                        <span className={`status ${statusClass[level]}`} aria-hidden />
                        <span>{message}</span>
                    </li>
                ))}
            </ul>
        )}
    </section>
);

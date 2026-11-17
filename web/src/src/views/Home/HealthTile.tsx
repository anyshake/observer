import { healthStatusClass, type IHealth } from './healthStatus';

interface IHealthTile extends IHealth {
    readonly title: string;
}

export const HealthTile = ({ title, value, detail, level }: IHealthTile) => (
    <div className="card app-panel app-health-tile">
        <div className="app-health-head">
            <span className="app-health-title">{title}</span>
            <span className={`status ${healthStatusClass[level]}`} aria-hidden />
        </div>
        <div className="app-health-value">{value}</div>
        <p className="app-health-detail">{detail}</p>
    </div>
);

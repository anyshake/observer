import Icon from '@mdi/react';
import { ReactNode } from 'react';

interface IChartPanel {
    readonly title: string;
    readonly icon: string;
    readonly extra?: ReactNode;
    readonly children: ReactNode;
}

export const ChartPanel = ({ title, icon, extra, children }: IChartPanel) => (
    <div className="app-chart-panel">
        <div className="app-chart-panel-header">
            <Icon className="app-chart-panel-icon" path={icon} size={0.65} />
            <span className="app-chart-panel-title">{title}</span>
            {extra && <div className="app-chart-panel-extra">{extra}</div>}
        </div>
        <div className="app-chart-panel-body">{children}</div>
    </div>
);

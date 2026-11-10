import Icon from '@mdi/react';
import { ReactNode } from 'react';

interface IChartPanel {
    readonly title: string;
    readonly icon: string;
    readonly children: ReactNode;
}

export const ChartPanel = ({ title, icon, children }: IChartPanel) => (
    <div className="app-chart-panel">
        <div className="app-chart-panel-header">
            <Icon className="app-chart-panel-icon" path={icon} size={0.65} />
            <span className="app-chart-panel-title">{title}</span>
        </div>
        <div className="app-chart-panel-body">{children}</div>
    </div>
);

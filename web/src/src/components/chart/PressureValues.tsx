import { formatChartNumber } from './formatChartNumber';

export const PressureValues = ({ value, unit }: { value?: number; unit: string }) => (
    <dl className="app-ground-motion">
        <div className="app-ground-motion-item">
            <dt className="app-ground-motion-label">P</dt>
            <dd className="app-ground-motion-value">
                {formatChartNumber(value)}
                <span className="app-ground-motion-unit">{unit}</span>
            </dd>
        </div>
    </dl>
);

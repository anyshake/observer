import type { TooltipComponentOption } from 'echarts';

import { getTimeString } from '../../helpers/utils/getTimeString';
import { formatChartNumber } from './formatChartNumber';

export const getPressureTooltip = (unit?: string): TooltipComponentOption => ({
    show: Boolean(unit),
    trigger: 'axis',
    axisPointer: { type: 'none' },
    backgroundColor: 'transparent',
    borderWidth: 0,
    padding: 0,
    extraCssText: 'box-shadow: none;',
    position: ([x, y], _params, _dom, _rect, { viewSize, contentSize }) => [
        x + 14 + 168 > viewSize[0] ? x - 14 - contentSize[0] : x + 14,
        y > viewSize[1] / 2 ? y - 14 - contentSize[1] : y + 14
    ],
    formatter: (params) => {
        const { value } = Array.isArray(params) ? params[0] : params;
        if (!Array.isArray(value) || typeof value[0] !== 'number' || typeof value[1] !== 'number') {
            return '';
        }

        const card = document.createElement('div');
        card.className = 'app-ground-motion-overlay';
        card.style.position = 'relative';
        const time = document.createElement('div');
        time.className = 'app-ground-motion-overlay-time';
        time.textContent = getTimeString(value[0]);
        const values = document.createElement('dl');
        values.className = 'app-ground-motion';
        const row = document.createElement('div');
        row.className = 'app-ground-motion-item';
        const label = document.createElement('dt');
        label.className = 'app-ground-motion-label';
        label.textContent = 'P';
        const reading = document.createElement('dd');
        reading.className = 'app-ground-motion-value';
        reading.textContent = formatChartNumber(value[1]);
        const suffix = document.createElement('span');
        suffix.className = 'app-ground-motion-unit';
        suffix.textContent = unit ?? '';
        reading.append(suffix);
        row.append(label, reading);
        values.append(row);
        card.append(time, values);
        return card;
    }
});

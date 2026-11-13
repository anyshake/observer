import { ReactNode, useCallback, useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { getTimeString } from '../../helpers/utils/getTimeString';
import { findGroundMotion, hasGroundMotion, IGroundMotionRecord } from './groundMotion';
import { GroundMotionValues } from './GroundMotionValues';
import type { HoverTimeHandler } from './useHoverTime';

interface IGroundMotionHover {
    // Returns the sorted per-second records at hover time; read lazily so live data
    // can grow without re-rendering the chart.
    readonly getRecords: () => readonly IGroundMotionRecord[] | undefined;
    readonly children: (onHoverTime: HoverTimeHandler) => ReactNode;
}

// Overlay distance from the pointer and the width reserved before flipping sides.
const OVERLAY_OFFSET = 14;
const OVERLAY_WIDTH = 168;

// Wraps a chart and shows the PGA/PGV/PGD of the second under the pointer.
export const GroundMotionHover = ({ getRecords, children }: IGroundMotionHover) => {
    const { t } = useTranslation();
    const containerRef = useRef<HTMLDivElement>(null);
    const frameRef = useRef(0);
    const [hover, setHover] = useState<{
        x: number;
        y: number;
        record: IGroundMotionRecord;
    } | null>(null);

    useEffect(() => () => cancelAnimationFrame(frameRef.current), []);

    const handleHover = useCallback<HoverTimeHandler>(
        (position) => {
            cancelAnimationFrame(frameRef.current);
            frameRef.current = requestAnimationFrame(() => {
                const record = position && findGroundMotion(getRecords(), position.time);
                setHover(position && record ? { x: position.x, y: position.y, record } : null);
            });
        },
        [getRecords]
    );

    const width = containerRef.current?.clientWidth ?? 0;
    const height = containerRef.current?.clientHeight ?? 0;
    const flipX = hover && hover.x + OVERLAY_OFFSET + OVERLAY_WIDTH > width;
    const flipY = hover && hover.y > height / 2;

    return (
        <div ref={containerRef} className="app-chart-hover-area">
            {children(handleHover)}
            {hover && (
                <div
                    className="app-ground-motion-overlay"
                    style={{
                        ...(flipX
                            ? { right: width - hover.x + OVERLAY_OFFSET }
                            : { left: hover.x + OVERLAY_OFFSET }),
                        ...(flipY
                            ? { bottom: height - hover.y + OVERLAY_OFFSET }
                            : { top: hover.y + OVERLAY_OFFSET })
                    }}
                >
                    <div className="app-ground-motion-overlay-time">
                        {getTimeString(hover.record.start)}
                    </div>
                    {hasGroundMotion(hover.record) ? (
                        <GroundMotionValues motion={hover.record} />
                    ) : (
                        <div className="app-ground-motion-overlay-empty">
                            {t('components.GroundMotion.unavailable')}
                        </div>
                    )}
                </div>
            )}
        </div>
    );
};

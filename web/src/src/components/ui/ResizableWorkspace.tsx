import { mdiDragHorizontalVariant } from '@mdi/js';
import Icon from '@mdi/react';
import {
    KeyboardEvent,
    PointerEvent,
    ReactNode,
    RefObject,
    useCallback,
    useRef,
    useState
} from 'react';

import { WorkspaceConstraints } from '../../config/constraints';
import { useLayoutStore } from '../../stores/layout';

interface IResizableWorkspace {
    readonly id: string;
    readonly className: string;
    readonly containerRef?: RefObject<HTMLDivElement | null>;
    readonly children: ReactNode;
}

interface IDragState {
    readonly startY: number;
    readonly startHeight: number;
    moved: boolean;
}

const { minHeight, maxHeight, keyboardStep } = WorkspaceConstraints;

const KEYBOARD_STEPS: Record<string, number> = {
    ArrowUp: -keyboardStep,
    ArrowDown: keyboardStep
};

const clampHeight = (height: number) =>
    Math.round(Math.min(maxHeight, Math.max(minHeight, height)));

export const ResizableWorkspace = ({
    id,
    className,
    containerRef,
    children
}: IResizableWorkspace) => {
    const { workspaceHeights, setWorkspaceHeight } = useLayoutStore();
    const fallbackRef = useRef<HTMLDivElement>(null);
    const workspaceRef = containerRef ?? fallbackRef;
    const dragRef = useRef<IDragState | null>(null);
    const [dragHeight, setDragHeight] = useState<number | null>(null);
    const height = dragHeight ?? workspaceHeights[id];

    const handlePointerDown = useCallback(
        (event: PointerEvent<HTMLDivElement>) => {
            const workspace = workspaceRef.current;
            if (!workspace || event.button !== 0) {
                return;
            }
            event.preventDefault();
            event.currentTarget.setPointerCapture(event.pointerId);
            dragRef.current = {
                startY: event.clientY,
                startHeight: workspace.offsetHeight,
                moved: false
            };
        },
        [workspaceRef]
    );

    const handlePointerMove = useCallback((event: PointerEvent<HTMLDivElement>) => {
        const drag = dragRef.current;
        if (!drag) {
            return;
        }
        drag.moved = drag.moved || event.clientY !== drag.startY;
        if (drag.moved) {
            setDragHeight(clampHeight(drag.startHeight + event.clientY - drag.startY));
        }
    }, []);

    const handlePointerUp = useCallback(
        (event: PointerEvent<HTMLDivElement>) => {
            const drag = dragRef.current;
            dragRef.current = null;
            setDragHeight(null);
            if (drag?.moved && event.type === 'pointerup') {
                setWorkspaceHeight(id, clampHeight(drag.startHeight + event.clientY - drag.startY));
            }
        },
        [id, setWorkspaceHeight]
    );

    const handleKeyDown = useCallback(
        (event: KeyboardEvent<HTMLDivElement>) => {
            const workspace = workspaceRef.current;
            const step = KEYBOARD_STEPS[event.key];
            if (!workspace || !step) {
                return;
            }
            event.preventDefault();
            setWorkspaceHeight(id, clampHeight(workspace.offsetHeight + step));
        },
        [id, setWorkspaceHeight, workspaceRef]
    );

    const handleReset = useCallback(() => setWorkspaceHeight(id, null), [id, setWorkspaceHeight]);

    return (
        <div className="app-workspace">
            <div ref={workspaceRef} className={className} style={height ? { height } : undefined}>
                {children}
            </div>
            <div
                role="separator"
                tabIndex={0}
                aria-orientation="horizontal"
                aria-valuemin={minHeight}
                aria-valuemax={maxHeight}
                aria-valuenow={height}
                className={`app-workspace-handle ${dragHeight !== null ? 'app-workspace-handle-active' : ''}`}
                onPointerDown={handlePointerDown}
                onPointerMove={handlePointerMove}
                onPointerUp={handlePointerUp}
                onPointerCancel={handlePointerUp}
                onKeyDown={handleKeyDown}
                onDoubleClick={handleReset}
            >
                <Icon path={mdiDragHorizontalVariant} size={0.8} />
            </div>
        </div>
    );
};

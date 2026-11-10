import 'react-resizable/css/styles.css';

import { ReactNode, RefObject, useEffect, useRef, useState } from 'react';
import Draggable from 'react-draggable';
import { ResizableBox } from 'react-resizable';

interface IDraggableBox {
    readonly locked: boolean;
    readonly children: ReactNode;
    readonly layout: {
        position: { x: number; y: number };
        size: { width: number; height: number };
    };
    readonly constraints: {
        minWidth: number;
        minHeight: number;
        maxWidth: number;
        maxHeight: number;
    };
    readonly onDragStart?: () => void;
    readonly onDragStop?: (x: number, y: number) => void;
    readonly onWindowResize?: (x: number, y: number) => void;
    readonly onResizeStop?: (width: number, height: number) => void;
}

export const DraggableBox = ({
    locked,
    children,
    layout,
    constraints,
    onDragStart,
    onDragStop,
    onResizeStop,
    onWindowResize
}: IDraggableBox) => {
    const [position, setPosition] = useState(layout.position);
    const dragRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        setPosition(layout.position);
    }, [layout.position]);

    useEffect(() => {
        const parent = dragRef.current?.parentElement;
        if (locked || !parent) {
            return;
        }
        const observer = new ResizeObserver(() => {
            const parentRect = parent.getBoundingClientRect();
            const newX = Math.max(0, Math.min(position.x, parentRect.width - layout.size.width));
            const newY = Math.max(0, Math.min(position.y, parentRect.height - layout.size.height));
            if (newX !== position.x || newY !== position.y) {
                setPosition({ x: newX, y: newY });
                onWindowResize?.(newX, newY);
            }
        });
        observer.observe(parent);
        return () => observer.disconnect();
    }, [position, layout.size.width, layout.size.height, onWindowResize, locked]);

    return (
        <Draggable
            nodeRef={dragRef as RefObject<HTMLElement>}
            disabled={locked}
            cancel=".react-resizable-handle, .app-draggable-cancel"
            defaultClassNameDragging="app-draggable-dragging"
            bounds="parent"
            grid={[1, 1]}
            position={position}
            onStart={() => {
                onDragStart?.();
            }}
            onStop={(_, { x, y }) => {
                setPosition({ x, y });
                onDragStop?.(x, y);
            }}
        >
            <div
                ref={dragRef}
                className={`app-draggable ${!locked ? 'app-draggable-unlocked' : ''}`}
            >
                <ResizableBox
                    axis="both"
                    resizeHandles={locked ? [] : ['se', 'ne', 'sw']}
                    draggableOpts={{ grid: [1, 1] }}
                    onResizeStop={(_, { size }) => {
                        onResizeStop?.(size.width, size.height);
                    }}
                    width={layout.size.width}
                    height={layout.size.height}
                    minConstraints={[constraints.minWidth, constraints.minHeight]}
                    maxConstraints={[constraints.maxWidth, constraints.maxHeight]}
                >
                    {children}
                </ResizableBox>
            </div>
        </Draggable>
    );
};

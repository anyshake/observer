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
        const handleResize = () => {
            if (locked) {
                return;
            }
            if (dragRef.current) {
                const parent = dragRef.current.parentElement;
                if (parent) {
                    const parentRect = parent.getBoundingClientRect();
                    let newX = Math.min(position.x, parentRect.width - layout.size.width);
                    let newY = Math.min(position.y, parentRect.height - layout.size.height);
                    newX = Math.max(0, newX);
                    newY = Math.max(0, newY);
                    setPosition({ x: newX, y: newY });
                    onWindowResize?.(newX, newY);
                }
            }
        };

        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, [position, layout.size.width, layout.size.height, onWindowResize, locked]);

    return (
        <Draggable
            nodeRef={dragRef as RefObject<HTMLElement>}
            disabled={locked}
            cancel=".react-resizable-handle"
            defaultClassNameDragging="border-2 border-dashed border-gray-300"
            bounds="parent"
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
                className={`absolute m-2 overflow-hidden rounded-lg bg-white p-2 shadow-lg ${!locked ? 'cursor-move' : ''}`}
            >
                <ResizableBox
                    axis="both"
                    resizeHandles={locked ? [] : ['se', 'ne', 'sw']}
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

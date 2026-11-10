import { format } from 'date-fns';
import { saveAs } from 'file-saver';

interface IRect {
    readonly x: number;
    readonly y: number;
    readonly width: number;
    readonly height: number;
}

const PADDING = 16;
const ELLIPSIS = '…';

const getRect = (element: Element, origin: DOMRect): IRect => {
    const rect = element.getBoundingClientRect();
    return {
        x: rect.left - origin.left,
        y: rect.top - origin.top,
        width: rect.width,
        height: rect.height
    };
};

const isVisible = ({ width, height }: IRect) => width > 0 && height > 0;

const getFont = (style: CSSStyleDeclaration) =>
    `${style.fontStyle} ${style.fontWeight} ${style.fontSize} ${style.fontFamily}`;

const fitText = (ctx: CanvasRenderingContext2D, text: string, maxWidth: number) => {
    if (ctx.measureText(text).width <= maxWidth) {
        return text;
    }
    let end = text.length;
    while (end > 0 && ctx.measureText(`${text.slice(0, end)}${ELLIPSIS}`).width > maxWidth) {
        end--;
    }
    return `${text.slice(0, end)}${ELLIPSIS}`;
};

const drawText = (
    ctx: CanvasRenderingContext2D,
    element: Element,
    origin: DOMRect,
    align: 'left' | 'center'
) => {
    const rect = getRect(element, origin);
    const text = element.textContent?.trim();
    if (!text || !isVisible(rect)) {
        return;
    }
    const style = getComputedStyle(element);
    ctx.font = getFont(style);
    ctx.fillStyle = style.color;
    ctx.textAlign = align;
    ctx.textBaseline = 'middle';
    const x = align === 'center' ? rect.x + rect.width / 2 : rect.x;
    ctx.fillText(fitText(ctx, text, rect.width), x, rect.y + rect.height / 2);
};

const drawBox = (ctx: CanvasRenderingContext2D, element: Element, origin: DOMRect) => {
    const style = getComputedStyle(element);
    const { x, y, width, height } = getRect(element, origin);
    ctx.beginPath();
    ctx.roundRect(x, y, width, height, parseFloat(style.borderTopLeftRadius) || 0);
    ctx.fillStyle = style.backgroundColor;
    ctx.fill();
    const border = parseFloat(style.borderTopWidth) || 0;
    if (border > 0) {
        ctx.lineWidth = border;
        ctx.strokeStyle = style.borderTopColor;
        ctx.stroke();
    }
};

const drawIcon = (ctx: CanvasRenderingContext2D, icon: SVGSVGElement, origin: DOMRect) => {
    const path = icon.querySelector('path');
    const rect = getRect(icon, origin);
    if (!path || !isVisible(rect)) {
        return;
    }
    const viewBox = icon.viewBox.baseVal;
    ctx.save();
    ctx.translate(rect.x, rect.y);
    ctx.scale(rect.width / (viewBox?.width || 24), rect.height / (viewBox?.height || 24));
    ctx.fillStyle = getComputedStyle(path).fill;
    ctx.fill(new Path2D(path.getAttribute('d') ?? ''));
    ctx.restore();
};

const drawHeader = (ctx: CanvasRenderingContext2D, header: Element, origin: DOMRect) => {
    const style = getComputedStyle(header);
    const rect = getRect(header, origin);
    const border = parseFloat(style.borderBottomWidth) || 0;
    if (border > 0) {
        ctx.fillStyle = style.borderBottomColor;
        ctx.fillRect(rect.x, rect.y + rect.height - border, rect.width, border);
    }
    const icon = header.querySelector<SVGSVGElement>(':scope > svg');
    if (icon) {
        drawIcon(ctx, icon, origin);
    }
    const title = header.querySelector('.app-chart-panel-title');
    if (title) {
        drawText(ctx, title, origin, 'left');
    }
};

const drawPanel = (ctx: CanvasRenderingContext2D, panel: Element, origin: DOMRect) => {
    drawBox(ctx, panel, origin);
    const header = panel.querySelector('.app-chart-panel-header');
    if (header) {
        drawHeader(ctx, header, origin);
    }
    panel.querySelectorAll('canvas').forEach((canvas) => {
        const rect = getRect(canvas, origin);
        if (isVisible(rect) && canvas.width > 0 && canvas.height > 0) {
            ctx.drawImage(canvas, rect.x, rect.y, rect.width, rect.height);
        }
    });
    panel
        .querySelectorAll('.app-empty-hint')
        .forEach((hint) => drawText(ctx, hint, origin, 'center'));
};

const getBounds = (panels: Element[], origin: DOMRect) => {
    const rects = panels.map((panel) => getRect(panel, origin));
    const left = Math.min(...rects.map(({ x }) => x)) - PADDING;
    const top = Math.min(...rects.map(({ y }) => y)) - PADDING;
    const right = Math.max(...rects.map(({ x, width }) => x + width)) + PADDING;
    const bottom = Math.max(...rects.map(({ y, height }) => y + height)) + PADDING;
    return { x: left, y: top, width: right - left, height: bottom - top };
};

export const exportWorkspace = (workspace: HTMLElement, panelSelector: string, fileName: string) =>
    new Promise<void>((resolve, reject) => {
        const panels = Array.from(workspace.querySelectorAll(panelSelector));
        if (panels.length === 0) {
            reject(new Error('workspace is empty'));
            return;
        }
        const origin = workspace.getBoundingClientRect();
        const bounds = getBounds(panels, origin);
        const scale = Math.max(1, window.devicePixelRatio || 1);
        const canvas = document.createElement('canvas');
        canvas.width = Math.ceil(bounds.width * scale);
        canvas.height = Math.ceil(bounds.height * scale);
        const ctx = canvas.getContext('2d');
        if (!ctx) {
            reject(new Error('canvas 2d context is unavailable'));
            return;
        }
        ctx.scale(scale, scale);
        ctx.translate(-bounds.x, -bounds.y);
        ctx.fillStyle = getComputedStyle(workspace).backgroundColor;
        ctx.fillRect(bounds.x, bounds.y, bounds.width, bounds.height);
        panels.forEach((panel) => drawPanel(ctx, panel, origin));
        canvas.toBlob((blob) => {
            if (!blob) {
                reject(new Error('failed to encode png'));
                return;
            }
            saveAs(blob, `${fileName}_${format(Date.now(), 'yyyyMMdd_HHmmss')}.png`);
            resolve();
        }, 'image/png');
    });

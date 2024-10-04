export class CircularQueue2D {
	private buffer: ArrayBuffer[];
	private views: Float64Array[];
	private head: number;
	private tail: number;
	private rows: number;
	private columns: number;
	private length: number;

	constructor(rows: number, columns: number) {
		this.buffer = Array(rows)
			.fill(null)
			.map(() => new ArrayBuffer(columns * Float64Array.BYTES_PER_ELEMENT));
		this.views = this.buffer.map((buf) => new Float64Array(buf));
		this.head = 0;
		this.tail = 0;
		this.rows = rows;
		this.columns = columns;
		this.length = 0;
	}

	write = (data: Float64Array): void => {
		if (data.length !== this.columns) {
			throw new Error(`Data must have ${this.columns} elements`);
		}
		for (let i = 0; i < data.length; i++) {
			this.views[this.tail][i] = data[i];
		}
		this.tail = (this.tail + 1) % this.rows;
		if (this.length < this.rows) {
			this.length++;
		} else {
			this.head = (this.head + 1) % this.rows;
		}
	};

	read = (size: number): Float64Array[] => {
		if (size > this.length) {
			throw new Error("Not enough rows to read");
		}
		const result: Float64Array[] = [];
		for (let i = 0; i < size; i++) {
			const currentIndex = (this.head + i) % this.rows;
			const row = new Float64Array(this.buffer[currentIndex]);
			result.push(row);
		}
		return result;
	};

	readAll = (): Float64Array[] => {
		const result: Float64Array[] = [];
		for (let i = 0; i < this.length; i++) {
			const currentIndex = (this.head + i) % this.rows;
			const row = new Float64Array(this.buffer[currentIndex]);
			result.push(row);
		}
		return result;
	};

	getShape = (): [number, number] => [this.rows, this.columns];
}

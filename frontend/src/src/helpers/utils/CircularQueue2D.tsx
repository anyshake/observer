type TypedArray =
	| Int8Array
	| Uint8Array
	| Uint8ClampedArray
	| Int16Array
	| Uint16Array
	| Int32Array
	| Uint32Array
	| Float32Array
	| Float64Array;

export class CircularQueue2D<T extends TypedArray> {
	private buffer: ArrayBuffer[];
	private views: T[];
	private head: number;
	private tail: number;
	private rows: number;
	private columns: number;
	private length: number;
	private TypedArrayConstructor: { new (buffer: ArrayBuffer): T };

	constructor(
		rows: number,
		columns: number,
		TypedArrayConstructor: { new (buffer: ArrayBuffer): T }
	) {
		this.buffer = Array(rows)
			.fill(null)
			.map(
				() =>
					new ArrayBuffer(
						columns * new TypedArrayConstructor(new ArrayBuffer(0)).BYTES_PER_ELEMENT
					)
			);
		this.TypedArrayConstructor = TypedArrayConstructor;
		this.views = this.buffer.map((buf) => new this.TypedArrayConstructor(buf));
		this.head = 0;
		this.tail = 0;
		this.rows = rows;
		this.columns = columns;
		this.length = 0;
	}

	write = (data: T): void => {
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

	read = (size: number): T[] => {
		if (size > this.length) {
			throw new Error("Not enough rows to read");
		}
		const result: T[] = [];
		for (let i = 0; i < size; i++) {
			const currentIndex = (this.head + i) % this.rows;
			const row = new this.TypedArrayConstructor(this.buffer[currentIndex]);
			result.push(row);
		}
		return result;
	};

	readAll = (): T[] => {
		const result: T[] = [];
		for (let i = 0; i < this.length; i++) {
			const currentIndex = (this.head + i) % this.rows;
			const row = new this.TypedArrayConstructor(this.buffer[currentIndex]);
			result.push(row);
		}
		return result;
	};

	getShape = (): [number, number] => [this.rows, this.columns];
}

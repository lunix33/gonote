import { Component, OnInit, Input, EventEmitter, Output } from '@angular/core';
import { FormControl } from '@angular/forms';
import { debounceTime, count } from 'rxjs/operators';

@Component({
	selector: 'app-markdown-editor',
	templateUrl: './markdown-editor.component.html',
	styleUrls: ['./markdown-editor.component.scss']
})
export class MarkdownEditorComponent implements OnInit {

	private _data: string = '';

	@Input()
	public set data(v: string) {
		if (v) {
			// Convert back encoded content of data.
			v = v
				.replace(/&amp;/g, '&')
				.replace(/&#39;/g, `'`)
				.replace(/&lt;/g, '<')
				.replace(/&gt;/g, '>')
				.replace(/&#34;/g, '"');

			this._data = v;

			if (this.input)
				this.input.value = v;
		}
	};
	public get data(): string { return this._data; }

	@Input('tab-height')
	public tabHeight: number = 0;

	@Output()
	public update: EventEmitter<string> = new EventEmitter<string>();

	@Output()
	public save: EventEmitter<void> = new EventEmitter<void>();

	public inputControl: FormControl = new FormControl();
	public input: HTMLTextAreaElement;

	constructor() { }

	ngOnInit() {
		if (!this.data)
			this.data = "";

		// Get input.
		this.input = document.querySelector('app-markdown-editor textarea');
		this.input.value = this.data;

		// Listen for input change.
		this.inputControl.valueChanges
			.pipe(debounceTime(500))
			.subscribe(this.onInputChange.bind(this));
	}

	// **** Manipulations *****************************************************
	/**
	 * List of input transformation.
	 */
	Apply(type: string) {
		const d: Manip = {
			str: this.input.value,
			start: this.input.selectionStart,
			end: this.input.selectionEnd
		}

		switch(type) {
			case 'header':
				this.ApplyHeader(d);
				break;
			case 'ulist':
				this.ApplyUList(d);
				break;
			case 'olist':
				this.ApplyOList(d);
				break;
			case 'bold':
				this.ApplyBold(d);
				break;
			case 'italic':
				this.ApplyItalic(d);
				break;
			case 'icode':
				this.ApplyICode(d);
				break;
			case 'link':
				this.ApplyLink(d);
				break;
			case 'image':
				this.ApplyImage(d);
				break;
			case 'codeblock':
				this.ApplyCodeBlock(d);
				break;
			case 'hr':
				this.ApplyHorizontalLine(d);
		}

		// Apply.
		this.data = d.str;
		this.input.value = d.str;
		this.input.focus();
		this.input.selectionStart = d.start;
		this.input.selectionEnd = d.end;
	}

	/**
	 * Apply a header on the manipulation object.
	 * @param d The manipulation object.
	 */
	private ApplyHeader(d: Manip) {
		const addChar = '#'
		this.mutateBeginingOfLine(d, (m: Manip) => {
			const follow = m.str.substring(m.seek, m.seek+1);
			if (follow === addChar)
				this.insertAtPosition(m, addChar);
			else
				this.insertAtPosition(m, `${addChar} `);
		});
	}

	/**
	 * Apply a unordered list on the manipulation object.
	 * @param d The manipulation object.
	 */
	private ApplyUList(d: Manip) {
		const addChar = '*'
		this.mutateBeginingOfLine(d, (m: Manip) => {
			const follow = m.str.substring(m.seek, m.seek+1);
			if (follow === addChar) {
				m.str = `${m.str.substr(0, m.seek)}${m.str.substr(m.seek + 2, m.str.length)}`;
				m.end -= 2;
			} else {
				this.insertAtPosition(m, `${addChar} `);
			}
		});
	}

	/**
	 * Apply a ordered list on the manipulation object.
	 * @param d The manipulation object.
	 */
	private ApplyOList(d: Manip) {
		this.mutateBeginingOfLine(d, (m: Manip) => {
			this.insertAtPosition(m, `${d.count}. `);
		});
	}

	/**
	 * Apply bold on the manipulation object.
	 * @param d The manipulation object.
	 */
	private ApplyBold(d: Manip) {
		this.ignoreLastWhitespace(d);
		this.insertBeforeAfter(d, '**');
	}

	/**
	 * Apply italic on the manipulation object.
	 * @param d The manipulation object.
	 */
	private ApplyItalic(d: Manip) {
		this.ignoreLastWhitespace(d);
		this.insertBeforeAfter(d, '_');
	}

	/**
	 * Apply inline code on the manipulation object.
	 * @param d The manipulation object.
	 */
	private ApplyICode(d: Manip) {
		this.ignoreLastWhitespace(d);
		this.insertBeforeAfter(d, '`');
	}

	/**
	 * Apply a link to the manipulation object.
	 * @param d The manipulation object.
	 */
	private ApplyLink(d: Manip) {
		let link: string = 'Destination Link',
			name: string = 'Display';

		if (d.start !== d.end) {
			link = d.str.substring(d.start, d.end);
			name = link;
		}

		const insert = `[${name}](${link})`;
		d.str = `${d.str.substring(0, d.start)}${d.str.substring(d.end, d.str.length)}`;

		d.seek = d.start;
		this.insertAtPosition(d, insert);

		d.end = d.start + insert.length;
	}

	/**
	 * Apply an image to the manipulation object.
	 * @param d The manipulation object.
	 */
	private ApplyImage(d: Manip) {
		let link: string = 'Picture Link',
			name: string = 'Alt. Text';

		if (d.start !== d.end) {
			link = d.str.substring(d.start, d.end);
			name = link;
		}

		const insert = `![Picture](${link})`;
		d.str = `${d.str.substring(0, d.start)}${d.str.substring(d.end, d.str.length)}`;

		d.seek = d.start;
		this.insertAtPosition(d, insert);

		d.end = d.start + insert.length;	
	}

	/**
	 * Apply a code block to the manipulation object.
	 * @param d The manipulation object.
	 */
	private ApplyCodeBlock(d: Manip) {
		let insertBefore = '```md\n',
			insertAfter = '\n```';

		// Insert new lines before.
		for (let i = 1; i < 3; i++) {
			const char = (d.start != 0) ? d.str.substring(d.start - i, d.start - (i - 1)) : '\n'
			if (char != '\n') {
				for (let j = 0; j < (3 - i); j++)
					insertBefore = `\n${insertBefore}`;
				break;
			}
		}
		if (d.start == d.end) {
			insertBefore = `${insertBefore}Code Here`
		}

		// Insert new lines after.
		for (let i = 1; i < 3; i++) {
			const char = (d.end != d.str.length) ? d.str.substring(d.end + (i - 1), d.end + i) : '\n'
			if (char != '\n') {
				for (let j = 0; j < (3 - i); j++)
					insertAfter = `${insertAfter}\n`;
				break;
			}
		}

		// Insert component
		d.seek = d.start;
		this.insertAtPosition(d, insertBefore);
		d.seek = d.end;
		this.insertAtPosition(d, insertAfter);
	}

	/**
	 * Apply a horizontal line to the manipulation object.
	 * @param d The manipulation object.
	 */
	private ApplyHorizontalLine(d: Manip) {
		let insert = `---`;

		// Insert new lines before
		for (let i = 1; i < 3; i++) {
			const char = (d.start != 0) ? d.str.substring(d.start - i, d.start - (i - 1)) : '\n'
			if (char != '\n') {
				for (let j = 0; j < (3 - i); j++)
					insert = `\n${insert}`;
				break;
			}
		}

		// Insert new lines after
		for (let i = 1; i < 3; i++) {
			const char = (d.start != d.str.length) ? d.str.substring(d.start + (i - 1), d.start + i) : '\n'
			if (char != '\n') {
				for (let j = 0; j < (3 - i); j++)
					insert = `${insert}\n`;
				break;
			}
		}

		// Insert component
		d.seek = d.start;
		this.insertAtPosition(d, insert);
	}
	
	// **** Events ************************************************************
	/**
	 * Event fired when the editor input change value.
	 * @param val The new value of the input.
	 */
	public onInputChange(val: string) {
		this.data = val;
		this.update.emit(this.data);
	}

	public onInputKeydown(e: KeyboardEvent) {
		switch(e.code) {
			case 'Tab':
				// Insert an actual tab character.
				e.preventDefault();
				const d: Manip = {
					str: this.input.value,
					seek: this.input.selectionStart
				}
				this.insertAtPosition(d, '	');
				this.input.value = d.str;
				break;
			case 'KeyS':
				// Annoying CTRL+S when typing!
				if (e.ctrlKey) {
					e.preventDefault();
					this.save.emit();
				}
				break;
		}
	}

	// **** Utilities *********************************************************
	/**
	 * Find the position of the beguinning of the line preceding a position.
	 * @param pos The position of the cursor in the line.
	 */
	private findBeginLine(pos): number {
		let char: string = '',
		    i: number;
		// Move position back until we reach the new line (or 0)
		for(i = 0; (pos - i) > 0 && char != '\n'; i++)
			char = this.input.value.slice((pos - i) - 1, (pos - i));
		pos -= i

		// Reajust for new lines.
		pos = (pos > 0) ? (pos + 1) : 0;

		// Move forward until we reach first letter (skip white spaces).
		char = ' ';
		for (i = 0; /\s/.test(char); i++)
			char = this.input.value.slice((pos + i), (pos + i) + 1);
		pos += (i - 1);

		return pos;
	}

	/**
	 * Add a string into another at the specified position.
	 * @param m The manipulation object.
	 * @param add The string to add.
	 */
	private insertAtPosition(m: Manip, add: string): void {
		m.str = `${m.str.substring(0, m.seek)}${add}${m.str.substring(m.seek, m.str.length)}`;

		// Correc the seek, start and end positions.
		if (m.seek < m.start)
			m.start += add.length;
		m.end += add.length;
		m.seek += add.length;
	}

	/**
	 * Insert a string both after and before of a selection.
	 * @param init The initial string.
	 * @param add The string to add.
	 * @param before The starting position.
	 * @param after The ending position.
	 */
	private insertBeforeAfter(m: Manip, add: string): void {
		m.seek = m.start;
		this.insertAtPosition(m, add);
		m.seek = m.end;
		this.insertAtPosition(m, add);
	}

	/**
	 * Test if the last character of the selection is a whitespace, if so it will be ignored.
	 * @param m The manipulation object.
	 */
	private ignoreLastWhitespace(m: Manip) {
		const char: string = m.str.substring(m.end - 1, m.end);
		if (/\s/.test(char))
			m.end--;
	}

	/**
	 * Run a mutation function at the beguinning of each line.
	 * @param m The manipulation object.
	 * @param mutator The mutation function.
	 */
	private mutateBeginingOfLine(m: Manip, mutator: MutatorFn) {
		let triggered: boolean = true;
		for (m.seek = this.findBeginLine(m.start), m.count = 0; m.seek <= m.end; m.seek++) {
			const char = m.str.slice((m.seek), m.seek + 1);
			if (triggered) {
				if (!/\s/.test(char)) {
					m.count++;
					mutator(m);
					triggered = false;
				}
			} else
				triggered = (char.includes('\n'));
		}
	}
}

type MutatorFn = (m: Manip)=>void;
type Manip = {
	str?: string,
	start?: number,
	end?: number,
	seek?: number,
	count?: number
}

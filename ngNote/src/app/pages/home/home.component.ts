import { Component } from '@angular/core';
import { FormControl } from '@angular/forms';
import { debounceTime } from 'rxjs/operators'

import { Note } from '../../class/note';
import { HttpStatusCode } from 'src/app/class/httpstatus';

@Component({
	selector: 'app-home',
	templateUrl: './home.component.html',
	styleUrls: ['./home.component.scss']
})
export class HomeComponent {
	public notes: Note[] = [];

	private _author: string;
	public get author(): string { return this._author; }
	public set author(v: string) {
		this._author = v;

		// Get all the notes to display.
		Note.GetAll(v)
			.then((notes: Note[]) => {
				for (let n of notes) {
					n.Added = new Date(n.Added).toLocaleString();
					n.Updated = new Date(n.Updated).toLocaleString();
				}
				this.notes = notes
			})
			.catch(err => console.error(err));
	}

	public authorControl: FormControl = new FormControl();

	constructor() {
		// Get the content of the author field.
		this.author = localStorage.getItem('author');

		// On author update
		this.authorControl.valueChanges
			.pipe(debounceTime(500))
			.subscribe(this.onAuthorInput.bind(this));
	}

	/**
	 * Event fired when the delete icon is clicked.
	 * @param e The event fired.
	 * @param n The note on the row
	 */
	async onDeleteBtnClick(e: Event, n: Note) {
		e.preventDefault();

		try {
			await n.Delete()
			const idx = this.notes.findIndex(x => x.ID === n.ID)
			this.notes.splice(idx, 1);
		}
		catch(err) {
			if (err instanceof Error) {
				console.error(err)
			} else {
				if (err.Status = HttpStatusCode.NotFound) {
					// Model
					console.warn(`(${n.ID}) ${n.Title} not found.`)
				} else {
					// Model
					console.error(`unable to delete (${n.ID}) ${n.Title}`)
				}
			}
		}
	}

	/**
	 * Event fired when the author field is modified.
	 * @param v The new value of the field.
	 */
	onAuthorInput(v: string) {
		this.author = v;
		localStorage.setItem('author', this.author);
	}
}

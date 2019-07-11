import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { Note } from 'src/app/class/note';

@Component({
	selector: 'app-note-edit',
	templateUrl: './note-edit.component.html',
	styleUrls: ['./note-edit.component.scss']
})
export class NoteEditComponent implements OnInit {

	private route: ActivatedRoute;
	private router: Router;

	public note: Note;
	private author: string = localStorage.getItem('author');

	constructor(route: ActivatedRoute, router: Router) {
		this.route = route;
		this.router = router;

		this.note = new Note({
			Author: this.author,
			Added: 'now'
		});
	}

	ngOnInit() {
		this.route.params.subscribe(async (p) => {
			const id = p['id'];
			if (id != null) {
				try {
					this.note = await Note.Get(id);
					if (this.note.Author != this.author)
						this.router.navigate(['/note', this.note.ID]);
				} catch (err) {
					this.note = null;
				}
				
				if (this.note == null)
					this.router.navigateByUrl('/not-found');
			}
		});
	}

	onMarkdownChange(md: string) {
		this.note.Content = md;
	}

	/**
	 * Event fired when the form is submitted.
	 */
	async onFormSubmit() {
		if (this.note.Title && this.note.Author && this.note.Content) {
			try {
				if (this.note.ID)
					await this.note.Update();
				else
					await this.note.Add();
				
				this.router.navigate(['/note', this.note.ID]);
			} catch(err) {
				console.error(err);
				// Model
				// An error occure while saving.
			}
		} else {
			// Model
			// Some info is not present.
		}
	}
}

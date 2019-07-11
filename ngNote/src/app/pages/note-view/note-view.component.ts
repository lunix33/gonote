import { Component, OnInit, Input } from '@angular/core';
import { ActivatedRoute } from '@angular/router'
import { Note } from 'src/app/class/note';

@Component({
	selector: 'app-note-view',
	templateUrl: './note-view.component.html',
	styleUrls: ['./note-view.component.scss']
})
export class NoteViewComponent implements OnInit {
	public note: Note = new Note({});

	private route: ActivatedRoute;
	public author: string = localStorage.getItem('author');

	constructor(route: ActivatedRoute) { 
		this.route = route;
	}

	ngOnInit() {
		this.route.params.subscribe(async (p) => {
			this.note = await Note.Get(p['id']);
		});
	}
}

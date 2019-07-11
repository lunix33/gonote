import { Component, Output, EventEmitter, HostListener } from '@angular/core';

import { Note } from '../../class/note'

@Component({
	selector: 'app-navigation',
	templateUrl: './navigation.component.html',
	styleUrls: ['./navigation.component.scss']
})
export class NavigationComponent {

	public notes: Note[] = [];
	public collapsed: boolean = false;
	private breakpoint: number = 768;

	constructor() {
		this.onWindowResize(null);

		Note.GetAll(localStorage.getItem('author'))
			.then((notes: Note[]) => { this.notes = notes })
			.catch(err => console.error(err));
	}

	/**
	 * Event fired when the window is resized.
	 * @param event The event fired
	 */
	@HostListener('window:resize', ['$event'])
	onWindowResize(event: Event) {
		// Get the client width and show/hide the menu based on the breakpoint.
		const width = document.querySelector('html').clientWidth;
		this.collapsed = width < this.breakpoint;
	}
}

import { HttpClient } from '@angular/common/http';

import { InjectorInstance, ApiUrl } from './common'
import { HTTPStatus, HttpStatusCode } from './httpstatus';

export class Note {
	// **** Static ************************************************************
	/**
	 * Fetch all the notes of the author from the API.
	 */
	static async GetAll(author: string): Promise<Note[]> {
		const http = InjectorInstance.get<HttpClient>(HttpClient);
		const data: HTTPStatus = ((await http.get(`${ApiUrl}/notes/${author}`).toPromise()) as HTTPStatus)
		
		if (data.Status === HttpStatusCode.Ok) {
			const rtn: Note[] = data.Message.map((x: any) => new Note(x))
			return rtn;
		} else
			throw data;
	}

	/**
	 * Fetch a specific note from the API.
	 * @param id The Id of the note to fetch.
	 */
	static async Get(id: string): Promise<Note> {
		const http = InjectorInstance.get<HttpClient>(HttpClient);
		const data: HTTPStatus = ((await http.get(`${ApiUrl}/note/${id}`).toPromise()) as HTTPStatus)

		if (data.Status === HttpStatusCode.Ok) {
			return new Note(data.Message);
		} else
			throw data
	}

	// **** Fields ************************************************************
	public ID: string;
	public Author: string;
	public Title: string;
	public Content: string;
	public Added: string;
	public Updated: string;

	/**
	 * Create a new note.
	 * @param o The options to create the new note.
	 */
	constructor(o: any) {
		this.ID = o.ID;
		this.Author = o.Author;
		this.Title = o.Title;
		this.Content = o.Content;
		this.Added = o.Added;
		this.Updated = o.Updated;
	}

	// **** Methods ***********************************************************
	/**
	 * Add the current note to the database.
	 * If the ID already exists, then the note will be cloned.
	 */
	async Add(): Promise<void> {
		const http = InjectorInstance.get<HttpClient>(HttpClient);
		const payload = {
			Author: this.Author,
			Title: this.Title,
			Content: this.Content
		}
		const data: HTTPStatus = ((await http.put(`${ApiUrl}/note/`, payload).toPromise()) as HTTPStatus)

		if (data.Status === HttpStatusCode.Ok) {
			// Update the current object if successful.
			const nNote: Note = data.Message;
			this.Added = nNote.Added;
			this.Updated = nNote.Updated;
			this.ID = nNote.ID
			return
		}

		throw data;
	}

	/**
	 * Update the current note on the database.
	 */
	async Update(): Promise<void> {
		if (!this.ID)
			throw new Error('Unable to update a non-existing note.');

		const http = InjectorInstance.get<HttpClient>(HttpClient);
		const payload = {
			Author: this.Author,
			Title: this.Title,
			Content: this.Content
		}
		const data = ((await http.patch(`${ApiUrl}/note/${this.ID}`, payload).toPromise()) as HTTPStatus)

		if (data.Status === HttpStatusCode.Ok) {
			const uNote: Note = data.Message;
			this.Updated = uNote.Updated;
			return;
		}

		throw data;
	}

	/**
	 * Delete the current note from the database.
	 */
	async Delete(): Promise<void> {
		if (!this.ID)
			throw new Error('Unable to delete a non-existing note.');

		const http = InjectorInstance.get<HttpClient>(HttpClient);
		const data: HTTPStatus = ((await http.delete(`${ApiUrl}/note/${this.ID}`).toPromise()) as HTTPStatus)

		if (data.Status === HttpStatusCode.Ok) {
			this.ID = undefined;
			return;
		}

		throw data;
	}
}

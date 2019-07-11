import { Component, Input, OnInit } from '@angular/core';

@Component({
	selector: 'app-card',
	templateUrl: './card.component.html',
	styleUrls: ['./card.component.scss']
})
export class CardComponent implements OnInit {

	@Input()
	public header: string;

	@Input()
	public footer: string;

	@Input()
	public nopad: string|boolean;

	constructor() { }

	ngOnInit() {
		this.nopad = (this.nopad === '');
	}
}

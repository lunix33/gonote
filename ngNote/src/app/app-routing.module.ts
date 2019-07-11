import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { HomeComponent } from './pages/home/home.component'
import { NotFoundComponent } from './pages/not-found/not-found.component'
import { NoteViewComponent } from './pages/note-view/note-view.component';
import { NoteEditComponent } from './pages/note-edit/note-edit.component';

const routes: Routes = [
  { path: 'home', component: HomeComponent },
  { path: '', redirectTo: '/home', pathMatch: 'full' },

  { path: 'note/add', component: NoteEditComponent},
  { path: 'note/:id', component: NoteViewComponent },
  { path: 'note/:id/edit', component: NoteEditComponent },

  { path: '**', component: NotFoundComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

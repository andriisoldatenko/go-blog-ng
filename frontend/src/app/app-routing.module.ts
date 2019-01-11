import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from 'src/app/home/home.component';
import { PostListComponent } from 'src/app/post-list/post-list.component';
import { EditPostComponent } from 'src/app/edit-post/edit-post.component';
import { OktaCallbackComponent } from '@okta/okta-angular';
import { LoginComponent } from 'src/app/login/login.component';

const routes: Routes = [
  {path: '', component: HomeComponent, pathMatch: 'full'},
  {path: 'implicit/callback', component: OktaCallbackComponent},
  {path: 'login', component: LoginComponent},
  {path: 'posts', component: PostListComponent},
  {path: 'posts/new', component: EditPostComponent},
  {path: 'posts/edit/:id', component: EditPostComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

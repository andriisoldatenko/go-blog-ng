import { Component, OnInit } from '@angular/core';
import { PostService } from 'src/app/shared/services/post.service';
import { IPost } from 'src/app/shared/models/post';
import { OktaAuthService } from '@okta/okta-angular';

@Component({
  selector: 'app-post-list',
  templateUrl: './post-list.component.html',
  styleUrls: ['./post-list.component.scss']
})
export class PostListComponent implements OnInit {
  private posts: IPost[];
  isAuthenticated: boolean;

  constructor(
    private postService: PostService,
    public oktaAuth: OktaAuthService
  ) { }

  async ngOnInit() {
    // get authentication state for immediate use
    this.isAuthenticated = await this.oktaAuth.isAuthenticated();

    // subscribe to authentication state changes
    this.oktaAuth.$authenticationState.subscribe(
      (isAuthenticated: boolean)  => this.isAuthenticated = isAuthenticated
    );
    this.postService.getPosts()
      .subscribe(
        (data: IPost[]) => this.posts = data['posts'],
        err => console.log(err)
      );
  }

  private deletePost(postId: number) {
    this.postService.deletePost(postId)
      .subscribe(() => {
        this.posts = this.posts.filter(post => post.id !== postId);
      });
  }
}

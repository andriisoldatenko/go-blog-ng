import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs/internal/Observable';
import { IPost } from '../models/post';

@Injectable({
  providedIn: 'root'
})
export class PostService {
  private apiUrl = '/api';

  constructor(
    private http: HttpClient
  ) { }

  /** GET: get all posts from the server */
  public getPosts() {
    return this.http.get(`${this.apiUrl}/posts`);
  }

  /** GET: get post by ID */
  public getPost(id: number): Observable<any> {
    return this.http.get<IPost>(`${this.apiUrl}/posts/${id}`);
  }

  /** POST: add a new post to the server */
  public addPost(post: IPost): Observable<any> {
    return this.http.post<IPost>(`${this.apiUrl}/posts`, {post});
  }

  /** PUT: update the post on the server */
  public savePost(post: IPost): Observable<any> {
    return this.http.put<IPost>(`${this.apiUrl}/posts/${post.id}`, post);
  }

  /** DELETE: delete the post from the server */
  public deletePost(id: number) {
    return this.http.delete(`${this.apiUrl}/posts/${id}`);
  }
}

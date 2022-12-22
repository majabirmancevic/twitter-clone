import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { PostRequest } from '../payloads/request/post';
import { PostResponse } from '../payloads/response/post';

@Injectable({
  providedIn: 'root'
})
export class PostService {
  
  constructor(private http: HttpClient) { }
  
  
  tweet(payload: PostRequest) {
    return this.http.post("https://localhost:8000/tweet_service/tweets", payload);
  }
  
 
  
  getTweetsByUsername(username: string): Observable<Array<PostResponse>>{
    return this.http.get<Array<PostResponse>>(`https://localhost:8000/tweet_service/tweets/${username}`);
  }
  
  getLikesByTweet(tweetId: string): Observable<Array<string>>{
    return this.http.get<Array<string>>(`https://localhost:8000/tweet_service/tweets/users/${tweetId}`);
  }
  
  // getAll(): Observable<PostResponse[]> {
  //   return this.http.get<PostResponse[]>("http://localhost:8080/posts");
  // }
}

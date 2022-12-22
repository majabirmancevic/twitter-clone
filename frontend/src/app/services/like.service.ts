import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class LikeService {

  constructor(private http: HttpClient) { }

  // isLiked(tweetId: number): Observable<boolean>{
  //   return this.http.post<boolean>("http://localhost:8080/likes/is-liked", {tweetId: tweetId});
  // }

  like(tweetId: string,username:string){
    return this.http.post(`https://localhost:8000/tweet_service/like/${tweetId}`, {username: username})
  }

  dislike(tweetId: string,username:string){
    return this.http.delete(`https://localhost:8000/tweet_service/dislike/${tweetId}/${username}`)
  }

  getLikeCounter(tweetId: string): Observable<number>{
    return this.http.get<number>(`https://localhost:8000/tweet_service/likes/count/${tweetId}`)
  }
}

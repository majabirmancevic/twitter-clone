import { Component, Input, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { PostResponse } from '../payloads/response/post';
import { PostService } from '../services/post.service';
import { TweetInput } from '../tweet-input';

@Component({
  selector: 'app-tweet',
  templateUrl: './tweet.component.html',
  styleUrls: ['./tweet.component.css']
})
export class TweetComponent  implements OnInit {

  @Input() footer: Boolean = true;
  @Input() counters: Boolean = false;
  @Input() tweet!: PostResponse ;



  likes: string[] = [];

  constructor() {
  
  }

  ngOnInit(): void {
    
  }

  
  
  // goToTweet(tweetId: number){
    
  //   this.router.navigate(["/tweet/" + tweetId], {state: {data: this.tweet}});
  // }

}

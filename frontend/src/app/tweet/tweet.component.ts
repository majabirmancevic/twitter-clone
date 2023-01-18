import { Component, Input, OnInit } from '@angular/core';
import { PostResponse } from '../payloads/response/post';


@Component({
  selector: 'app-tweet',
  templateUrl: './tweet.component.html',
  styleUrls: ['./tweet.component.css']
})
export class TweetComponent  implements OnInit {

  @Input() footer: boolean = true;
  @Input() counters: boolean = false;
  @Input() tweet!: PostResponse ;



  likes: string[] = [];

  constructor() {
  
  }

  ngOnInit(): void {
    
  }

  


}

import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-like-list',
  templateUrl: './like-list.component.html',
  styleUrls: ['./like-list.component.css']
})
export class LikeListComponent implements OnInit {

  //@Output() setTweetId = new EventEmitter<string | undefined>();
  @Input() likes! : string[];

  //tweetId!:string;

  constructor(private activatedRoute:ActivatedRoute) { }

  ngOnInit(): void {
    //this.tweetId = this.activatedRoute.snapshot.params['tweetId'];
  }

}

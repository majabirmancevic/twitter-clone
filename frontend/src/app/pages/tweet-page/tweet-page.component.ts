import { Component, OnInit } from '@angular/core';
import { NavigationEnd, Router } from '@angular/router';
import { OverlayForm } from 'src/app/overlay-form';
import { PostResponse } from 'src/app/payloads/response/post';
import { PostService } from 'src/app/services/post.service';

@Component({
  selector: 'app-tweet-page',
  templateUrl: './tweet-page.component.html',
  styleUrls: ['./tweet-page.component.css']
})
export class TweetPageComponent extends OverlayForm implements OnInit {

  replies: Array<PostResponse>;
  isFocus: boolean = false;
  navigationSubscription;

  constructor(private router: Router, private postService: PostService) { 
    super();
    
    this.tweet = history.state.data;
    console.log(history.state.data);
    this.replies = new Array();

    this.router.routeReuseStrategy.shouldReuseRoute = function () {
      return false;
    };
    
    this.navigationSubscription = this.router.events.subscribe((event: any)=>{
      if(event instanceof NavigationEnd){
        this.router.navigated = false;
      }
      
    })

  }

  ngOnInit(): void {
    
  }

  ngOnDestroy() {
    if (this.navigationSubscription) {
      this.navigationSubscription.unsubscribe();
    }
  }

  focus(status: boolean){
    this.isFocus = status;
  }

  reply(){
    const self = this;
    this.postService.tweet(this.payload).subscribe({
      complete(){
        self.form.reset();
        self.router.navigate(["/tweet/" + self.tweet.id], {state: {data: self.tweet}});
      },
      error(error){console.log(error)}
    });
  }
}
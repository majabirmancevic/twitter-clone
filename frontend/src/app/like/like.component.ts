import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Router } from '@angular/router';
import { faHeart } from '@fortawesome/free-regular-svg-icons';
import { faHeart as faSolidHeart } from '@fortawesome/free-solid-svg-icons';
import { PostResponse } from '../payloads/response/post';
import { AuthService } from '../services/auth.service';
import { LikeService } from '../services/like.service';
import { PostService } from '../services/post.service';
import { SingleAction } from '../single-action';


@Component({
  selector: 'app-like',
  templateUrl: './like.component.html',
  styleUrls: ['./like.component.css']
})
export class LikeComponent extends SingleAction implements OnInit {
  
  likeCounter! : number;
  liked = false;
 
  display = false;
  @Input() override tweet! : PostResponse;
  //@Input() 
  likes : Array<string> = [];

  //@Output() viewLike = new EventEmitter<string>();
  
  constructor(private likeService: LikeService, private router: Router, private authService:AuthService,private postService: PostService) { 
    super();
    this.faIcon = faHeart;
    this.faSolidIcon = faSolidHeart;
    
  }

  ngOnInit(): void {
    //this.likeService.isLiked(this.tweet.id).subscribe(response => this.isActive = response);
    this.likeService.getLikeCounter(this.tweet.id).subscribe(counter => {this.likeCounter = counter});
    this.viewLike()
  }

  like() {  
    const username :string = this.authService.getUsername()
    this.likeService.like(this.tweet.id,username).subscribe(data => {
      this.updateCounter();
    });
    this.liked = true;
  }

  dislike(){
    const username :string = this.authService.getUsername()
    this.likeService.dislike(this.tweet.id,username).subscribe(data => {
      this.updateCounter();
    });
    this.liked = false;
  }

  viewLike() {
    this.postService.getLikesByTweet(this.tweet.id).subscribe(data => {
      this.likes = data ;
    })
  }

  onPress(){
    this.display = !this.display
    console.log("KLIKNUTO NA KAUNTER")
    console.log(this.display)
  }

  updateCounter(){
    this.likeService.getLikeCounter(this.tweet.id).subscribe(counter => {this.likeCounter = counter});
  }
}

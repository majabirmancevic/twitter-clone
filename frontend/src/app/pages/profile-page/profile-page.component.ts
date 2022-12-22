import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, NavigationEnd, Router } from '@angular/router';
import { OverlayForm } from 'src/app/overlay-form';
import { PostResponse } from 'src/app/payloads/response/post';
import { AuthService } from 'src/app/services/auth.service';
import { PostService } from 'src/app/services/post.service';
import { UserService } from 'src/app/services/user.service';
import { RegularUser } from 'src/app/user-model';
import { ImageCroppedEvent, base64ToFile } from 'ngx-image-cropper';
import { BusinessUser } from 'src/app/user-model-business';

@Component({
  selector: 'app-profile-page',
  templateUrl: './profile-page.component.html',
  styleUrls: ['./profile-page.component.css']
})
export class ProfilePageComponent extends OverlayForm implements OnInit {

  user: RegularUser;
  tweets: Array<PostResponse> = [];
  currentUsername! : string;
  

  constructor(
    private userService: UserService,
    private authService: AuthService,
    private activateRoute: ActivatedRoute,
    private router: Router,
    private postService: PostService) {

    super();

   
    this.user = {
      name: "",
      lastname: "",
      gender : "",
      age : 0,
      placeOfLiving: "",
      email: "",
      username: "",
      verified: false,
      role: "regular",
    };

    this.userService.getRegularUser(this.activateRoute.snapshot.params['username']).subscribe(user => {
      this.user = user;
    });
    
    this.postService.getTweetsByUsername(this.activateRoute.snapshot.params['username']).subscribe(posts => {
      this.tweets = posts;
    })

  }

    ngOnInit(): void {
      
    }
    

  getUrl(){
    return this.router.url;
  }

  goToChangePass(){
    this.router.navigate(['/change-password/',this.activateRoute.snapshot.params['username']]);
  }

  fetchTweets() {
    this.postService.getTweetsByUsername(this.user.username).subscribe(response => this.tweets = response);
  }



}
function ngOnInit() {
  throw new Error('Function not implemented.');
}


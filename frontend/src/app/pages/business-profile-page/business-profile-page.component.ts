import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { OverlayForm } from 'src/app/overlay-form';
import { PostResponse } from 'src/app/payloads/response/post';
import { AuthService } from 'src/app/services/auth.service';
import { PostService } from 'src/app/services/post.service';
import { UserService } from 'src/app/services/user.service';
import { BusinessUser } from 'src/app/user-model-business';

@Component({
  selector: 'app-business-profile-page',
  templateUrl: './business-profile-page.component.html',
  styleUrls: ['./business-profile-page.component.css']
})
export class BusinessProfilePageComponent extends OverlayForm implements OnInit {

  user: BusinessUser;
  tweets: Array<PostResponse> = [];
  currentUsername! : string;
  
  constructor(private userService: UserService,
    private authService: AuthService,
    private activateRoute: ActivatedRoute,
    private router: Router,
     private postService: PostService) {  

      super();

      this.user = {
        companyName: "",
        webSite: "",
        email: "",
        username: "",
        verified: false,
        role: "business",
      };

      this.userService.getBusinessUser(this.activateRoute.snapshot.params['username']).subscribe(user => {
        this.user = user;
      });

      this.postService.getTweetsByUsername(this.activateRoute.snapshot.params['username']).subscribe(posts => {
        this.tweets = posts;
      })

    }

    getUrl(){
      return this.router.url;
    }

    fetchTweets() {
      this.postService.getTweetsByUsername(this.user.username).subscribe(response => this.tweets = response);
    }

    goToChangePass(){
      this.router.navigate(['/change-password/',this.activateRoute.snapshot.params['username']]);
    }

  ngOnInit(): void {
  }

}

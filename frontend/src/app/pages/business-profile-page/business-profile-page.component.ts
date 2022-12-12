import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from 'src/app/services/auth.service';
import { PostService } from 'src/app/services/post.service';
import { UserService } from 'src/app/services/user.service';
import { BusinessUser } from 'src/app/user-model-business';

@Component({
  selector: 'app-business-profile-page',
  templateUrl: './business-profile-page.component.html',
  styleUrls: ['./business-profile-page.component.css']
})
export class BusinessProfilePageComponent implements OnInit {

  user: BusinessUser;
  // tweets: Array<PostResponse>;
  currentUsername! : string;
  
  constructor(private userService: UserService,
    private authService: AuthService,
    private activateRoute: ActivatedRoute,
    private router: Router) {  

      //super();

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

    }

    getUrl(){
      return this.router.url;
    }

  ngOnInit(): void {
  }

}

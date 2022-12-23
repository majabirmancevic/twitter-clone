import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css']
})
export class NavbarComponent implements OnInit {

  constructor(private router: Router, private authService: AuthService) { }
  username!: string | null;
  user:any;
  ngOnInit(): void {
    this.username = this.authService.getUsername();
    this.authService.getUser(localStorage.getItem("username")!).subscribe(data => {
      this.user = data
    })

  }

  getUrl(){
    return this.router.url;
  }

  navigate(){
    if (this.user.role == "regular" ){
      this.router.navigate(['/profile/',this.username]);
    }else {
      this.router.navigate(['/profile-business/',this.username]);
    }
  }
  
  logout(){
    this.authService.logout();
  }

}

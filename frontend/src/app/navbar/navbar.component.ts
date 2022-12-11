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

  ngOnInit(): void {
    this.username = this.authService.getUsername()
  }

  getUrl(){
    return this.router.url;
  }

  navigate(){
    this.router.navigate(['/profile/',this.username]);
  }
  // logout(){
  //   this.authService.logout();
  // }

}

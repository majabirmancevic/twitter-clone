import { Component, OnInit } from '@angular/core';
import { PostResponse } from 'src/app/payloads/response/post';

@Component({
  selector: 'app-business-profile-page',
  templateUrl: './business-profile-page.component.html',
  styleUrls: ['./business-profile-page.component.css']
})
export class BusinessProfilePageComponent implements OnInit {

  user: RegularUser;
  tweets: Array<PostResponse>;
  currentUsername! : string;
  
  constructor() { }

  ngOnInit(): void {
  }

}

import { Component, Input, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-like-list',
  templateUrl: './like-list.component.html',
  styleUrls: ['./like-list.component.css']
})
export class LikeListComponent implements OnInit {


  @Input() likes! : string[];

  

  constructor(private activatedRoute:ActivatedRoute) { }

  ngOnInit(): void {
   
  }

}

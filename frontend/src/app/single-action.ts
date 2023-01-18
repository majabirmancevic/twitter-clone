import { Component } from "@angular/core";
import { TweetInput } from "./tweet-input";



@Component({
    template: ''
  })

export abstract class SingleAction extends TweetInput {
    faIcon: any;
    faSolidIcon: any;
    isActive: boolean = false;
}
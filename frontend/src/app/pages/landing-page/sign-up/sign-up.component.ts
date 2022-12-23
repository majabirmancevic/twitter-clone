import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { ToastrService } from 'ngx-toastr';
import { SignUpPayload } from 'src/app/payloads/request/sign-up';
import { AuthService } from 'src/app/services/auth.service';


@Component({
  selector: 'app-sign-up',
  templateUrl: './sign-up.component.html',
  styleUrls: ['./sign-up.component.css']
})
export class SignUpComponent implements OnInit {

  form: FormGroup;
  payload: SignUpPayload;
  captcha: string;
  constructor(private authService: AuthService, private toastr: ToastrService, private router: Router) {

    this.form = new FormGroup({
      name: new FormControl("", Validators.required),
      lastname: new FormControl("", Validators.required),
      gender:new FormControl("",Validators.required),
      age: new FormControl("", [Validators.min(13), Validators.required]),
      placeOfLiving: new FormControl("",Validators.required),
      email: new FormControl("", [Validators.required, Validators.email]),
      username: new FormControl("", {
        validators: [Validators.required],
        updateOn: "blur"
      }), 
      password: new FormControl("", Validators.required)
    });

    this.payload = {
      name: "",
      lastname: "",
      gender : "",
      age : 0,
      placeOfLiving: "",
      email: "",
      username: "", 
      password: ""
    }

    this.captcha = '';

  }

  ngOnInit(): void {
  }

  resolved(captchaResponse: string){
    this.captcha = captchaResponse;
    console.log('resovled capctha with response: ' + this.captcha);
  }

  signUp() {
    this.payload.name = this.form.get('name')?.value;
    this.payload.lastname = this.form.get('lastname')?.value;
    this.payload.gender = this.form.get('gender')?.value;
    this.payload.age = this.form.get('age')?.value;
    this.payload.placeOfLiving = this.form.get('placeOfLiving')?.value;
    this.payload.email = this.form.get('email')?.value;
    this.payload.username = this.form.get('username')?.value;  
    this.payload.password = this.form.get('password')?.value;
    const self = this;
    
    
    self.router.navigate(['/verify-email'], { queryParams: { registered: 'true', verified: 'false' } });
    //self.toastr.info("Check your email address!");
    this.authService.signUp(this.payload).subscribe({
      next() {},
      complete(){},
      error(error) {
        console.log(error);
        console.log("------- GRESKA -------")
        console.log(error.code)
        //self.toastr.error("Something went wrong!");
      }
    })
  }

  get username() {
    return this.form.controls['username'];
  }

}

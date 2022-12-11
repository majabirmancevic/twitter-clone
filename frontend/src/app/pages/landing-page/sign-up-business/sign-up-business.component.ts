import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { ToastrService } from 'ngx-toastr';
import { SignUpPayloadBusiness } from 'src/app/payloads/request/sign-up-business';
import { AuthService } from 'src/app/services/auth.service';


@Component({
  selector: 'app-sign-up-business',
  templateUrl: './sign-up-business.component.html',
  styleUrls: ['./sign-up-business.component.css']
})
export class SignUpBusinessComponent implements OnInit {

  form: FormGroup;
  payload: SignUpPayloadBusiness;
  captcha: string;

  constructor(private authService: AuthService, private toastr: ToastrService, private router: Router) {

    this.form = new FormGroup({
      companyName: new FormControl("", Validators.required),
      email: new FormControl("", [Validators.required, Validators.email]),
      webSite: new FormControl("", Validators.required),
      username: new FormControl("", {
        validators: [Validators.required],
        updateOn: "blur"
      }), 
      password: new FormControl("", Validators.required)
    });

    this.payload = {
      companyName: "",
      webSite: "",
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

  signUpBusiness(){

    this.payload.companyName = this.form.get('companyName')?.value;
    this.payload.webSite = this.form.get('webSite')?.value;
    this.payload.email = this.form.get('email')?.value;
    this.payload.username = this.form.get('username')?.value;  
    this.payload.password = this.form.get('password')?.value;
    const self = this;
    
    
    self.router.navigate(['/verify-email'], { queryParams: { registered: 'true', verified: 'false' } });
    this.authService.signUpBusiness(this.payload).subscribe({
      next() {},
      complete(){},
      error(error) {
        console.log(error);
        console.log("------- GRESKA -------")
        console.log(error.code)
      }
    })

  }

  get username() {
    return this.form.controls

}
}

import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { ToastrService } from 'ngx-toastr';
import { ResetPasswordPayload } from 'src/app/payloads/request/reset-password';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'app-reset-password',
  templateUrl: './reset-password.component.html',
  styleUrls: ['./reset-password.component.css']
})
export class ResetPasswordComponent implements OnInit {

  form: FormGroup;
  payload : ResetPasswordPayload;
  username!: string;
  constructor(private activatedRoute:ActivatedRoute, private toastr: ToastrService, private authService: AuthService, private router: Router) {

    this.form = new FormGroup({
      password: new FormControl("", [Validators.required, Validators.min(8)]),
      repeatPassword: new FormControl("", [Validators.required, Validators.min(8)])
    })

    this.payload = {
      newPassword: "",
      repeatNewPassword: ""
    }
   
    
  }

  ngOnInit(): void {
    this.username = this.activatedRoute.snapshot.params['username'];
  }

  resetPassword(){

    this.payload.newPassword = this.form.get("password")?.value;
    this.payload.repeatNewPassword = this.form.get("repeatPassword")?.value;

    this.authService.resetPassword(this.payload,this.username).subscribe(data => {
    })
    this.toastr.success("Success reset password");
    this.router.navigate(['/sign-in']); 
  }

}

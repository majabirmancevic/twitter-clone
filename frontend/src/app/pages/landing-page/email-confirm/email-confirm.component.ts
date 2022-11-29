import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { ToastrService } from 'ngx-toastr';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'app-email-confirm',
  templateUrl: './email-confirm.component.html',
  styleUrls: ['./email-confirm.component.css']
})
export class EmailConfirmComponent implements OnInit {

  form: FormGroup;
  code: any;

  constructor(private activatedRoute:ActivatedRoute, private toastr: ToastrService, private authService: AuthService, private router: Router) {
    this.form = new FormGroup({
      verificationCode: new FormControl("", Validators.required)    
    })

   }

  ngOnInit(): void {
    this.activatedRoute.queryParams.subscribe(params =>{
      if(params['registered'] !== undefined && params['registered'] === "true"){
        this.toastr.success("An email with a verification code has been sent to your email");
        console.log("Sign Up Successful");
      }
    })
    this.code = this.activatedRoute.snapshot.paramMap.get('code')
  }

  verifyEmail(){
    this.code = this.form.get('code')?.value;
    this.authService.verifyEmail(this.code).subscribe({
      complete() {
        self.router.navigate(['/verifyemail'], { queryParams: { registered: 'true' , verified: 'false' } });
      },
      error(error) {
        console.log(error);
        self.toastr.error("Something went wrong! Error : ", error.code);
      }
    })
  }
}

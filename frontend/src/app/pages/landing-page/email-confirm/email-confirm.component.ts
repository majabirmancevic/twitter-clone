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
  inputcode: any;

  constructor(private activatedRoute:ActivatedRoute, private toastr: ToastrService, private authService: AuthService, private router: Router) {
    
    this.form = new FormGroup({
      code: new FormControl("", Validators.required)    
    })

    this.inputcode = ""

   }

  ngOnInit(): void {
    
      this.activatedRoute.queryParams.subscribe(params =>{
        if(params['registered'] !== undefined && params['registered'] === "true"){
          this.toastr.success("An email with a verification code has been sent to your email");
          console.log("Sent verification code to email");
        }
      })
  
  }

  verifyEmail(){
    this.inputcode = this.form.get('code')?.value;
    const self = this;
    this.authService.verifyEmail(this.inputcode).subscribe({
      next(){
       
      },
      complete() {
        self.router.navigate(['/sign-in'], { queryParams: { registered: 'true' , verified: 'true' } });
      },
      error(error) {
        console.log(error);
        self.toastr.error("Something went wrong! Error : ", error.code);
      }
    })
  }
}

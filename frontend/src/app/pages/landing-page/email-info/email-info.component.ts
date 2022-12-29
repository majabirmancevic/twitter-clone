import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators} from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { ToastrService } from 'ngx-toastr';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'app-email-info',
  templateUrl: './email-info.component.html',
  styleUrls: ['./email-info.component.css']
})
export class EmailInfoComponent implements OnInit {

  form: FormGroup;
  emailStr : string;
  user :string;

  constructor(private activatedRoute:ActivatedRoute, private toastr: ToastrService, private authService: AuthService, private router: Router) {

    this.form = new FormGroup({
      email: new FormControl("", Validators.required),    
      username : new FormControl("", Validators.required)  
    }),

    this.emailStr = "" 
    this.user = ""
   }

  ngOnInit(): void {
  }

  verifyEmail(){
    this.emailStr = this.form.get('email')?.value;
    this.user = this.form.get('username')?.value;
    const self = this;
    this.authService.sendverifyEmail(this.emailStr,this.user).subscribe(data =>{  
        this.toastr.info("Check your email address ..");
        return data
    })
    self.router.navigate(['/info']);
    
  }
  }



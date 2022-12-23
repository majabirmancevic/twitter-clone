import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from 'src/app/services/auth.service';
import { UserService } from 'src/app/services/user.service';
import { PasswordDto } from './password-dto';

@Component({
  selector: 'app-change-password',
  templateUrl: './change-password.component.html',
  styleUrls: ['./change-password.component.css']
})
export class ChangePasswordComponent implements OnInit {

  form: FormGroup;

  username!: string ;
  passwordDto: PasswordDto ;
  
  constructor( private authService: AuthService,
    private userService : UserService,
    private route: ActivatedRoute,
    private router: Router) { 

      this.form = new FormGroup({
        oldPassword: new FormControl('', [Validators.required]),
        newPassword: new FormControl('', [Validators.required, Validators.maxLength(32), Validators.minLength(8)])
      });
     
      this.passwordDto = {
        oldPassword : '',
        newPassword : ''
      }

   
    }

  ngOnInit(): void {
    this.username = this.route.snapshot.params['username'];
  }

  change()  {
    this.passwordDto.oldPassword = this.form.get('oldPassword')?.value;
    this.passwordDto.newPassword = this.form.get('newPassword')?.value;

    this.userService.changePassword(this.username, this.passwordDto).subscribe(data => {     
      return data;   
    });

    this.router.navigate(['/sign-in']);    
}
}

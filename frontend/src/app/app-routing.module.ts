import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomePageComponent } from './pages/home-page/home-page.component';
import { EmailConfirmComponent } from './pages/landing-page/email-confirm/email-confirm.component';
import { LandingPageComponent } from './pages/landing-page/landing-page.component';
import { SignInComponent } from './pages/landing-page/sign-in/sign-in.component';
import { SignUpComponent } from './pages/landing-page/sign-up/sign-up.component';
import { ProfilePageComponent } from './pages/profile-page/profile-page.component';
import { TweetPageComponent } from './pages/tweet-page/tweet-page.component';
import { AuthGuard } from './services/auth.guard';
import { EmailInfoComponent } from './pages/landing-page/email-info/email-info.component';
import { SignUpBusinessComponent } from './pages/landing-page/sign-up-business/sign-up-business.component';
import {ReactiveFormsModule,FormsModule} from '@angular/forms';
import { ChangePasswordComponent } from './pages/change-password/change-password.component';
import { BusinessProfilePageComponent } from './pages/business-profile-page/business-profile-page.component';
import { ResetPasswordComponent } from './pages/reset-password/reset-password.component';
import { InfoMessageComponent } from './info-message/info-message.component';



const routes: Routes = [
  {path : "" , component: LandingPageComponent},
  {path : "sign-in" , component: SignInComponent},
  {path : "sign-up" , component: SignUpComponent},
  {path : "sign-up-business", component: SignUpBusinessComponent},
  {path : "verify-email" , component: EmailConfirmComponent },
    // children: [{path : ":code" , component: EmailConfirmComponent}]
  {path : "email-info" , component: EmailInfoComponent },
  {path : "home" , component: HomePageComponent, canActivate: [AuthGuard]},
  {path : "tweet/:id" , component: TweetPageComponent, canActivate: [AuthGuard]},
  {path : "profile/:username" , component: ProfilePageComponent, canActivate: [AuthGuard]},
  {path : "profile-business/:username" , component:BusinessProfilePageComponent, canActivate: [AuthGuard]},
  {path : "profile/likes" , component: ProfilePageComponent, canActivate: [AuthGuard]},
  {path : "profile/update" , component: ProfilePageComponent, canActivate: [AuthGuard]},
  {path : "change-password/:username", component : ChangePasswordComponent, canActivate: [AuthGuard]},
  {path : "reset-password/:username/:code", component : ResetPasswordComponent},
  {path : "info", component : InfoMessageComponent}



];

@NgModule({
  imports: [RouterModule.forRoot(routes),  ReactiveFormsModule,FormsModule],
  exports: [RouterModule]
})
export class AppRoutingModule { }

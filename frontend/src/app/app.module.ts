import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule} from '@angular/forms';
import {ReactiveFormsModule} from '@angular/forms';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LikeComponent } from './like/like.component';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http'

import { TweetComponent } from './tweet/tweet.component';
import { TweetPageComponent } from './pages/tweet-page/tweet-page.component';
import { HomePageComponent } from './pages/home-page/home-page.component';
import { ProfilePageComponent } from './pages/profile-page/profile-page.component';
import { LandingPageComponent } from './pages/landing-page/landing-page.component';
import { ImageCropperModule } from 'ngx-image-cropper';
import { SignInComponent } from './pages/landing-page/sign-in/sign-in.component';
import { SignUpComponent } from './pages/landing-page/sign-up/sign-up.component';
import { SignUpBusinessComponent } from './pages/landing-page/sign-up-business/sign-up-business.component';
import { ToastrModule } from 'ngx-toastr';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NavbarComponent } from './navbar/navbar.component';
import { AuthGuard } from './services/auth.guard';
import { TokenInterceptor } from './services/token-interceptor';
import { NgxWebstorageModule } from 'ngx-webstorage';
import { EmailConfirmComponent } from './pages/landing-page/email-confirm/email-confirm.component';
import { RecaptchaModule} from 'ng-recaptcha';
import { EmailInfoComponent } from './pages/landing-page/email-info/email-info.component';
import { BusinessProfilePageComponent } from './pages/business-profile-page/business-profile-page.component';


@NgModule({
  declarations: [
    AppComponent,
    LikeComponent,
    TweetComponent,
    TweetPageComponent,
    HomePageComponent,
    ProfilePageComponent,
    LandingPageComponent,
    SignInComponent,
    SignUpComponent,
    NavbarComponent,
    EmailConfirmComponent,
    EmailInfoComponent,

    BusinessProfilePageComponent,

    SignUpBusinessComponent

  ],

    
  imports: [
    BrowserModule,
    AppRoutingModule,
    FontAwesomeModule,
    HttpClientModule,
    ImageCropperModule,
    ToastrModule.forRoot(),
    BrowserAnimationsModule,
    NgxWebstorageModule.forRoot(),
    RecaptchaModule,
    FormsModule,
    ReactiveFormsModule,
  ],
  providers: [
    AuthGuard,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: TokenInterceptor,
      multi: true
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }

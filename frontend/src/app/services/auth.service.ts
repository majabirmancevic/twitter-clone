import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { LocalStorageService } from 'ngx-webstorage';
import { map, tap, Observable } from 'rxjs';
import { SignInRequestPayload } from '../payloads/request/sign-in';
import { SignUpPayload } from '../payloads/request/sign-up';
import { SignInResponsePayload } from '../payloads/response/sign-in';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  
  
  constructor(private http: HttpClient, private localStorage: LocalStorageService, private router: Router) { }
  
  
  signUp(payload: SignUpPayload) {

    const headers = new HttpHeaders({
      'Accept': 'application/json',
      'Content-Type': 'application/json',

    }); 
    const options = { headers: headers };
       
    return this.http.post("https://localhost:8001/", JSON.stringify(payload), options);
  }
  
  signIn(payload: SignInRequestPayload) {

    const headers = new HttpHeaders({
      'Accept': 'application/json',
      'Content-Type': 'application/json',
    }); 
    const options = { headers: headers };
    const self = this;

    return this.http.post<SignInResponsePayload>("https://localhost:8001/login", JSON.stringify(payload), options).pipe<SignInResponsePayload>(
      map(response=>{
        self.localStorage.store("token", response.token);
        self.localStorage.store("username", response.user.username);
        self.localStorage.store("isLoggedIn", true);
        return response;
      })
    )
  }

  verifyEmail(code: any) { 
    const headers = new HttpHeaders({
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    }); 
    const options = { headers: headers };
    return this.http.get(`https://localhost:8001/verifyEmail/${code}`,  options);
    
  }
  
  logout() {
    const self = this;
    return this.http.post("http://localhost:8080/auth/logout", {username: this.getUsername(), refreshToken: this.localStorage.retrieve("refreshToken")}).subscribe({complete(){

      self.localStorage.clear('accessToken');
      self.localStorage.clear('refreshToken');
      self.localStorage.clear('expiresAt');
      self.localStorage.clear('username');
      self.localStorage.store('isLoggedIn', false);
      self.router.navigateByUrl("");
    }});
  }
  
  getUsername(){
    return this.localStorage.retrieve("username");
  }
  
  isLoggedIn() {
    return this.localStorage.retrieve("isLoggedIn");
  }

  refreshAccessToken(): Observable<SignInResponsePayload> {
    return this.http.post<SignInResponsePayload>("http://localhost:8080/auth/refresh-token", {
      refreshToken: this.localStorage.retrieve("refreshToken"),
      username: this.getUsername()
    }).pipe(tap(response=>{
      this.localStorage.store('accessToken', response.token);
      //this.localStorage.store('expiresAt', response.expiresAt);
  }))
  }

  getAccessToken() {
    return this.localStorage.retrieve("accessToken");
  }
  
  // findAllUsernames(): Observable<String[]> {
  //   return this.http.get<String[]>("http://localhost:8080/auth/usernames");
  // }
}

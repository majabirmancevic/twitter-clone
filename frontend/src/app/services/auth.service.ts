import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { LocalStorageService } from 'ngx-webstorage';
import { map, tap, Observable } from 'rxjs';
import { SignInRequestPayload } from '../payloads/request/sign-in';
import { SignUpPayload } from '../payloads/request/sign-up';
import { SignUpPayloadBusiness } from '../payloads/request/sign-up-business';
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
       
    return this.http.post("https://localhost:8002/", JSON.stringify(payload), { headers: headers });
  }
  
  
  signUpBusiness(payload: SignUpPayloadBusiness) {

    const headers = new HttpHeaders({
      'Accept': 'application/json',
      'Content-Type': 'application/json',

    }); 
    const options = { headers: headers };
       
    return this.http.post("https://localhost:8002/business", JSON.stringify(payload), { headers: headers });
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
        localStorage.setItem("token", response.token);
        localStorage.setItem("username", response.username);
        localStorage.setItem("isLoggedIn", "true");
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
    return this.http.get(`https://localhost:8002/verifyEmail/${code}`,  options);
    
  }

  
  getUsername() : string{
    return localStorage.getItem("username")!;
  }

  getUser(username: string){
    return this.http.get(`https://localhost:8000/profile_service/user/${username}`);
  }

  getJwtToken() {
    return localStorage.getItem('token');
  }

  // getUsernameByToken(){
  //   return JSON.parse(this.localStorage.retrieve("token")).username
  // }
  
  isLoggedIn() {
    return localStorage.getItem("isLoggedIn") == "true";
  }

  getAccessToken() {
    return localStorage.getItem("token");
  }
  
  logout() {
      this.router.navigateByUrl("");   
      //this.router.navigate(['']);
      localStorage.clear()

  }

  private parseToken(){
    let jwt = localStorage.getItem('token');
    if(jwt !== null){
      let jwtData = jwt.split('.')[1];
      let decodedJwtJsonData = atob(jwtData);
      let decodedJwtData = JSON.parse(decodedJwtJsonData);
      return decodedJwtData
    }
  }

  getUsernameFromToken(): string{
    let token = this.parseToken();

    if(token) {
      return this.parseToken()['sub']
    }
    return "";
  }

 
  // findAllUsernames(): Observable<String[]> {
  //   return this.http.get<String[]>("http://localhost:8080/auth/usernames");
  // }
}

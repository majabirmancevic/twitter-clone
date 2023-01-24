import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { PasswordDto } from '../pages/change-password/password-dto';
import { RegularUser } from '../user-model';
import { BusinessUser } from '../user-model-business';
import { AuthService } from './auth.service';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  
  payload: FormData;
  username! : string|null ;

  constructor(private http: HttpClient, private authService: AuthService) {
    this.payload = new FormData;
    this.username = this.authService.getUsername()
  }

  getRegularUser(username: string): Observable<RegularUser>{
    return this.http.get<RegularUser>(`https://localhost:8000/profile_service/user/${username}`);
  }


  changePassword(username: string, payload: PasswordDto){
    return this.http.post(`https://localhost:8000/profile_service/changePassword/${username}`, JSON.stringify(payload))
  }

  getBusinessUser(username: string): Observable<BusinessUser>{
    return this.http.get<BusinessUser>(`https://localhost:8000/profile_service/business/user/${username}`);
  }

 



  // editProfile(payload: any){
  //   console.log(payload)
  //   return this.http.post("http://localhost:8080/user/edit-profile", payload)
  // }

}

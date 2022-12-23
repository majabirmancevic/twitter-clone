import { HttpInterceptor, HttpRequest, HttpHandler, HttpEvent, HttpErrorResponse } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { BehaviorSubject, Observable, catchError, throwError, switchMap, filter, take } from "rxjs";
import { SignInResponsePayload } from "../payloads/response/sign-in";
import { AuthService } from "./auth.service";

@Injectable({
    providedIn: 'root'
})

export class TokenInterceptor implements HttpInterceptor {

    constructor(public authService: AuthService) { }

    intercept(request: HttpRequest<any>, next: HttpHandler):Observable<HttpEvent<any>> {

    if (this.authService.isLoggedIn()) {
        request = request.clone({
            setHeaders: {
                Authorization: `Bearer ${this.authService.getJwtToken()}` 
            }
        });
    }
              return next.handle(request);
            }    

   

}

 // addToken(req: HttpRequest<any>, jwtToken: any) {
    //     return req.clone({
    //         headers: req.headers.set('Authorization',
    //             'Bearer ' + jwtToken)
    //     });
    // }
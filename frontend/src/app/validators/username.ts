


// export function usernameValidator(authService:AuthService): AsyncValidatorFn{
//     return (control:AbstractControl):Observable<ValidationErrors|null>=>{
//         return authService.findAllUsernames().pipe(
//             map(usernames => {
//                 const username = usernames.find(username => username.toLowerCase() == control.value.toLowerCase());
//                 return username ? {usernameExists:true} : null;
//             })
//         )
//     }
// }
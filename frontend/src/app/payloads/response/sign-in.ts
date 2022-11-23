import { User } from "src/app/user-model";

export interface SignInResponsePayload{
    token: string;
    user: User;
    //expiresAt: Date
}
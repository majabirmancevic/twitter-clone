import { RegularUser } from "src/app/user-model";

export interface SignInResponsePayload{
    token: string;
    username: string;
    //expiresAt: Date
}
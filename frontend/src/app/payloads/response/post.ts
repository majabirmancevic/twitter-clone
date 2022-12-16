import { RegularUser } from "src/app/user-model";

export interface PostResponse{
    id?: number,  
    regular_username: string, 
    description: string,
    likeCounter: number
}
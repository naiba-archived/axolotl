import axios from "@/utils/axios";
import { Authority } from "@/store";

export async function fetchUser() {
    const resp: any = await axios.get("/api/user")
    let authority;
    switch (resp.authority) {
        default:
            authority = Authority.User
            break;
    }
    return {
        id: resp.id,
        nickname: resp.nickname,
        githubId: resp.github_id,
        authority
    }
}

export async function logout() {
    try {
        await axios.post("/api/user/logout")
    } catch (error) {
    }
}
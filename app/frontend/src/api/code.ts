import axios from "@/utils/axios";

export async function fetchCodeList() {
    const resp: string = await axios.get("/api/code/list");
    return JSON.parse(resp);
}

export async function executeCode(room: string, code: string, container: string) {
    const resp: string = await axios.post("/api/code/run", { room, code, container });
    return resp;
}
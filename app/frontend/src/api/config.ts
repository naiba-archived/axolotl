import axios from "@/utils/axios";

let configCache: any;

export async function getConfig() {
    if (configCache) {
        return configCache;
    }
    configCache = await axios.get("/api/config")
    return configCache
}
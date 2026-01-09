import { retrieveToken } from "@/utils/auth";

/**
 * Formats the authorization header that is to be sent
 * in any requests to the backend
 * @param token the auth token
 * @param contentType content type, will default to application/json
 * @returns the formatted header
 */
export const authHeader = (token: string, contentType: string = "application/json") => {
    return {
        "Content-Type": contentType,
        Authorization: `Bearer ${token}`,
    };
};

/**
 * Wraps a function that needs authentication token
 * Uses retrieveToken() which accesses cookies
 */
export const authWrapper =
    <T>() =>
    async (fn: (token: string) => Promise<T>) => {
        const token = await retrieveToken();
        return fn(token);
    };

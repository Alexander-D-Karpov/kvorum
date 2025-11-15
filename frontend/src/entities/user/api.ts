import fetcher from "@/shared/api/fetcher";

export interface User {
    id: string;
    display_name?: string;
    email?: string;
    phone?: string;
    roles?: string[];
}

export async function fetchCurrentUser(): Promise<User> {
    return fetcher<User>("/api/v1/me");
}

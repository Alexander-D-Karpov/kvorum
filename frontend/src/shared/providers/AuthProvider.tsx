import { createContext, useContext, useMemo, ReactNode, useEffect, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { APIError } from "@/shared/api/fetcher";
import { fetchCurrentUser, User } from "@/entities/user/api";
import { getInitData, isMaxWebApp, ready } from "@/shared/lib/maxBridge";

interface AuthContextValue {
    user: User | null;
    isLoading: boolean;
    isAuthenticated: boolean;
    isOrganizer: boolean;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [initDataSent, setInitDataSent] = useState(false);

    useEffect(() => {
        if (isMaxWebApp()) {
            ready();

            const initData = getInitData();
            if (initData && !initDataSent) {
                fetch('/api/v1/auth/max/webapp', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    credentials: 'include',
                    body: JSON.stringify({ initData }),
                })
                    .then(() => {
                        setInitDataSent(true);
                    })
                    .catch((err) => {
                        console.error('Failed to exchange webapp data:', err);
                    });
            }
        }
    }, [initDataSent]);

    const { data, isLoading } = useQuery({
        queryKey: ["me"],
        queryFn: async () => {
            try {
                const user = await fetchCurrentUser();
                return user;
            } catch (e) {
                if (e instanceof APIError && e.status === 401) {
                    return null;
                }
                throw e;
            }
        },
        retry: false,
        enabled: !isMaxWebApp() || initDataSent,
    });

    const value = useMemo<AuthContextValue>(() => {
        const user = (data ?? null) as User | null;
        const roles = user?.roles ?? [];
        const isOrganizer =
            roles.includes("organizer") ||
            roles.includes("coorganizer") ||
            roles.includes("admin");
        return {
            user,
            isLoading,
            isAuthenticated: !!user,
            isOrganizer,
        };
    }, [data, isLoading]);

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
    const ctx = useContext(AuthContext);
    if (!ctx) {
        throw new Error("useAuth must be used within AuthProvider");
    }
    return ctx;
}
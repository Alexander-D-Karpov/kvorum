interface MaxWebApp {
    initData: string
    initDataUnsafe: {
        query_id?: string
        user?: {
            id: number
            first_name: string
            last_name?: string
            username?: string
            language_code?: string
            photo_url?: string
        }
        auth_date?: number
        hash?: string
        start_param?: string
    }
    platform: string
    version: string
    ready: () => void
    close: () => void
    expand: () => void
    BackButton: {
        isVisible: boolean
        onClick: (callback: () => void) => void
        offClick: (callback: () => void) => void
        show: () => void
        hide: () => void
    }
}

declare global {
    interface Window {
        WebApp?: MaxWebApp
    }
}

export function getMaxWebApp(): MaxWebApp | null {
    if (typeof window === 'undefined') return null
    return window.WebApp || null
}

export function isMaxWebApp(): boolean {
    return getMaxWebApp() !== null
}

export function getInitData(): string {
    const webApp = getMaxWebApp()
    return webApp?.initData || ''
}

export function getStartParam(): string | undefined {
    const webApp = getMaxWebApp()
    return webApp?.initDataUnsafe?.start_param
}

export function ready(): void {
    const webApp = getMaxWebApp()
    webApp?.ready()
}

export function close(): void {
    const webApp = getMaxWebApp()
    webApp?.close()
}

export function expand(): void {
    const webApp = getMaxWebApp()
    webApp?.expand()
}
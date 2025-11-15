import { ReactNode } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/shared/providers/AuthProvider'

interface MainLayoutProps {
    children: ReactNode
}

export default function MainLayout({ children }: MainLayoutProps) {
    const { t, i18n } = useTranslation()
    const { user, isAuthenticated } = useAuth()
    const location = useLocation()

    const changeLang = () => {
        i18n.changeLanguage(i18n.language === 'ru' ? 'en' : 'ru')
    }

    return (
        <div className="min-h-screen bg-background text-foreground flex flex-col">
            <header className="border-b bg-card">
                <div className="container mx-auto px-4 py-3 flex items-center justify-between gap-4">
                    <div className="flex items-center gap-4">
                        <Link to="/" className="font-semibold text-lg">
                            kvorum
                        </Link>
                        <nav className="hidden md:flex items-center gap-3 text-sm">
                            <Link
                                to="/"
                                className={
                                    location.pathname === '/'
                                        ? 'font-medium text-primary'
                                        : 'text-muted-foreground hover:text-foreground'
                                }
                            >
                                {t('nav.home')}
                            </Link>
                            <Link
                                to="/me"
                                className={
                                    location.pathname.startsWith('/me')
                                        ? 'font-medium text-primary'
                                        : 'text-muted-foreground hover:text-foreground'
                                }
                            >
                                {t('nav.me')}
                            </Link>
                            <Link
                                to="/organizer"
                                className={
                                    location.pathname.startsWith('/organizer') ||
                                    location.pathname.includes('/e/') ||
                                    location.pathname.startsWith('/analytics')
                                        ? 'font-medium text-primary'
                                        : 'text-muted-foreground hover:text-foreground'
                                }
                            >
                                {t('nav.organizer')}
                            </Link>
                        </nav>
                    </div>

                    <div className="flex items-center gap-2">
                        <Button variant="outline" size="sm" onClick={changeLang}>
                            {i18n.language === 'ru' ? 'EN' : 'RU'}
                        </Button>
                        {isAuthenticated && (
                            <div className="hidden sm:flex items-center gap-2 text-sm text-muted-foreground">
                                <span>{user?.display_name || user?.email || user?.id}</span>
                            </div>
                        )}
                    </div>
                </div>
            </header>
            <main className="flex-1">
                {children}
            </main>
        </div>
    )
}

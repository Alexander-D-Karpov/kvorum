import { ReactNode } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useAuth } from '@/shared/providers/AuthProvider'

interface OrganizerRouteProps {
    children: ReactNode
}

export function OrganizerRoute({ children }: OrganizerRouteProps) {
    const { isOrganizer, isLoading } = useAuth()
    const location = useLocation()

    if (isLoading) {
        return <div className="flex items-center justify-center min-h-[200px]">Загрузка...</div>
    }

    if (!isOrganizer) {
        return <Navigate to="/" state={{ from: location }} replace />
    }

    return <>{children}</>
}

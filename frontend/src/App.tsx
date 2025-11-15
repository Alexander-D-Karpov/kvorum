import { Routes, Route } from 'react-router-dom'
import MainLayout from './layouts/MainLayout'
import EventPublicPage from './pages/EventPublicPage'
import AttendeeDashboard from './pages/AttendeeDashboard'
import OrganizerConsole from './pages/OrganizerConsole'
import EventFormsPage from './pages/EventFormsPage'
import CampaignsPage from './pages/CampaignsPage'
import AnalyticsPage from './pages/AnalyticsPage'
import CheckinPage from './pages/CheckinPage'
import NotFoundPage from './pages/NotFoundPage'
import { ProtectedRoute } from './shared/routing/ProtectedRoute'
import { OrganizerRoute } from './shared/routing/OrganizerRoute'

function HomePage() {
    return (
        <div className="container mx-auto max-w-3xl px-4 py-12 space-y-4">
            <h1 className="text-3xl font-bold">kvorum</h1>
            <p className="text-sm text-muted-foreground">
                Платформа для событий, регистрации, чек-ина и аналитики
            </p>
        </div>
    )
}

export default function App() {
    return (
        <MainLayout>
            <Routes>
                <Route path="/" element={<HomePage />} />
                <Route path="/e/:eventId" element={<EventPublicPage />} />
                <Route
                    path="/me"
                    element={
                        <ProtectedRoute>
                            <AttendeeDashboard />
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/organizer"
                    element={
                        <OrganizerRoute>
                            <OrganizerConsole />
                        </OrganizerRoute>
                    }
                />
                <Route
                    path="/e/:eventId/forms"
                    element={
                        <OrganizerRoute>
                            <EventFormsPage />
                        </OrganizerRoute>
                    }
                />
                <Route
                    path="/e/:eventId/campaigns"
                    element={
                        <OrganizerRoute>
                            <CampaignsPage />
                        </OrganizerRoute>
                    }
                />
                <Route
                    path="/analytics/:eventId"
                    element={
                        <OrganizerRoute>
                            <AnalyticsPage />
                        </OrganizerRoute>
                    }
                />
                <Route
                    path="/checkin/:eventId"
                    element={
                        <OrganizerRoute>
                            <CheckinPage />
                        </OrganizerRoute>
                    }
                />
                <Route path="*" element={<NotFoundPage />} />
            </Routes>
        </MainLayout>
    )
}

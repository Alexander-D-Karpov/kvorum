import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useEventAnalytics } from '@/entities/analytics/api'

export default function AnalyticsPage() {
    const { eventId } = useParams<{ eventId: string }>()
    const [from, setFrom] = useState('')
    const [to, setTo] = useState('')

    const { data, isLoading } = useEventAnalytics(eventId || '', from || undefined, to || undefined)

    const handleDownloadCSV = () => {
        if (!eventId) return
        const params = new URLSearchParams()
        if (from) params.set('from', from)
        if (to) params.set('to', to)
        const query = params.toString()
        const url = query
            ? `/api/v1/events/${eventId}/analytics.csv?${query}`
            : `/api/v1/events/${eventId}/analytics.csv`
        window.location.href = url
    }

    if (!eventId) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                Не указан идентификатор события
            </div>
        )
    }

    return (
        <div className="container mx-auto max-w-5xl px-4 py-8 space-y-6">
            <div className="flex items-center justify-between gap-4">
                <h1 className="text-2xl font-bold">Аналитика события</h1>
                <Button onClick={handleDownloadCSV}>Экспорт CSV</Button>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle className="text-base">Фильтры</CardTitle>
                </CardHeader>
                <CardContent className="grid md:grid-cols-4 gap-3">
                    <div className="space-y-1">
                        <p className="text-xs text-muted-foreground">С даты</p>
                        <Input
                            type="date"
                            value={from}
                            onChange={(e) => setFrom(e.target.value)}
                        />
                    </div>
                    <div className="space-y-1">
                        <p className="text-xs text-muted-foreground">По дату</p>
                        <Input
                            type="date"
                            value={to}
                            onChange={(e) => setTo(e.target.value)}
                        />
                    </div>
                </CardContent>
            </Card>

            {isLoading && (
                <div className="text-sm text-muted-foreground">Загрузка аналитики...</div>
            )}

            {data && (
                <>
                    <div className="grid gap-4 md:grid-cols-3">
                        <MetricCard title="Всего регистраций" value={data.total_registrations} />
                        <MetricCard title="Идут" value={data.going} />
                        <MetricCard title="Не идут" value={data.not_going} />
                        <MetricCard title="Возможно" value={data.maybe} />
                        <MetricCard title="Лист ожидания" value={data.waitlist} />
                        <MetricCard title="Пришли на событие" value={data.checked_in} />
                    </div>

                    {data.by_source && Object.keys(data.by_source).length > 0 && (
                        <Card>
                            <CardHeader>
                                <CardTitle className="text-base">Регистрации по источникам</CardTitle>
                            </CardHeader>
                            <CardContent className="space-y-2">
                                {Object.entries(data.by_source).map(([source, count]) => (
                                    <div key={source} className="flex items-center justify-between text-sm">
                                        <span className="text-muted-foreground">{source || 'unknown'}</span>
                                        <span className="font-medium">{count}</span>
                                    </div>
                                ))}
                            </CardContent>
                        </Card>
                    )}
                </>
            )}
        </div>
    )
}

function MetricCard({ title, value }: { title: string; value: number }) {
    return (
        <Card>
            <CardHeader>
                <CardTitle className="text-sm text-muted-foreground">{title}</CardTitle>
            </CardHeader>
            <CardContent>
                <div className="text-3xl font-bold">{value}</div>
            </CardContent>
        </Card>
    )
}

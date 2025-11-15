import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useMyOrganizedEvents, useCreateEvent } from '@/entities/event/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useToast } from '@/components/ui/use-toast'

export default function OrganizerConsole() {
    const { data: events, isLoading } = useMyOrganizedEvents()
    const createEvent = useCreateEvent()
    const { toast } = useToast()
    const navigate = useNavigate()

    const [title, setTitle] = useState('')
    const [description, setDescription] = useState('')

    const handleCreate = async () => {
        if (!title.trim()) return
        try {
            const ev = await createEvent.mutateAsync({ title, description })
            setTitle('')
            setDescription('')
            navigate(`/e/${ev.id}/edit`)
        } catch {
            toast({
                title: 'Ошибка',
                description: 'Не удалось создать событие',
                variant: 'destructive',
            })
        }
    }

    return (
        <div className="container mx-auto max-w-5xl px-4 py-8 space-y-8">
            <div className="flex items-center justify-between gap-4">
                <h1 className="text-2xl font-bold">Консоль организатора</h1>
                <Button variant="outline" size="sm" asChild>
                    <a href="/me">Кабинет участника</a>
                </Button>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle className="text-base">Создать событие</CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                    <Input
                        placeholder="Название события"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                    />
                    <Input
                        placeholder="Краткое описание"
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                    />
                    <Button onClick={handleCreate} disabled={createEvent.isPending || !title.trim()}>
                        Создать
                    </Button>
                </CardContent>
            </Card>

            <div className="space-y-3">
                <h2 className="text-xl font-semibold">Мои события</h2>
                {isLoading && (
                    <div className="text-sm text-muted-foreground">Загрузка событий...</div>
                )}
                {!isLoading && (!events || events.length === 0) && (
                    <div className="text-sm text-muted-foreground">У вас пока нет событий</div>
                )}
                {events && events.length > 0 && (
                    <div className="space-y-3">
                        {events.map((event) => (
                            <Card key={event.id}>
                                <CardHeader className="flex flex-row items-center justify-between gap-3">
                                    <div>
                                        <CardTitle className="text-base">{event.title}</CardTitle>
                                        <p className="text-xs text-muted-foreground">
                                            {event.starts_at
                                                ? new Date(event.starts_at).toLocaleString('ru-RU')
                                                : 'Дата не указана'}
                                        </p>
                                    </div>
                                    <div className="flex flex-wrap gap-2">
                                        <Button size="sm" variant="outline" onClick={() => navigate(`/e/${event.id}`)}>
                                            Публичная страница
                                        </Button>
                                        <Button size="sm" onClick={() => navigate(`/e/${event.id}/edit`)}>
                                            Редактировать
                                        </Button>
                                        <Button
                                            size="sm"
                                            variant="outline"
                                            onClick={() => navigate(`/analytics/${event.id}`)}
                                        >
                                            Аналитика
                                        </Button>
                                    </div>
                                </CardHeader>
                            </Card>
                        ))}
                    </div>
                )}
            </div>
        </div>
    )
}

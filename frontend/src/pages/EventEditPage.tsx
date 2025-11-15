import { useParams } from 'react-router-dom'
import { useEvent, useUpdateEvent, usePublishEvent } from '@/entities/event/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { useState, useEffect } from 'react'
import { useToast } from '@/components/ui/use-toast'

export default function EventEditPage() {
    const { eventId } = useParams<{ eventId: string }>()
    const { data: event, isLoading } = useEvent(eventId!)
    const updateEvent = useUpdateEvent(eventId!)
    const publishEvent = usePublishEvent(eventId!)
    const { toast } = useToast()

    const [title, setTitle] = useState('')
    const [description, setDescription] = useState('')
    const [location, setLocation] = useState('')
    const [onlineUrl, setOnlineUrl] = useState('')

    useEffect(() => {
        if (event) {
            setTitle(event.title)
            setDescription(event.description)
            setLocation(event.location || '')
            setOnlineUrl(event.online_url || '')
        }
    }, [event])

    const handleSave = async () => {
        try {
            await updateEvent.mutateAsync({
                title,
                description,
                location,
                online_url: onlineUrl,
            })
            toast({
                title: 'Сохранено',
                description: 'Изменения успешно сохранены',
            })
        } catch (error) {
            toast({
                title: 'Ошибка',
                description: 'Не удалось сохранить изменения',
                variant: 'destructive',
            })
        }
    }

    const handlePublish = async () => {
        try {
            await publishEvent.mutateAsync()
            toast({
                title: 'Опубликовано',
                description: 'Событие успешно опубликовано',
            })
        } catch (error) {
            toast({
                title: 'Ошибка',
                description: 'Не удалось опубликовать событие',
                variant: 'destructive',
            })
        }
    }

    if (isLoading) {
        return <div className="flex items-center justify-center min-h-screen">Загрузка...</div>
    }

    return (
        <div className="container mx-auto py-8 px-4 max-w-4xl">
            <Card>
                <CardHeader>
                    <CardTitle>Редактирование события</CardTitle>
                </CardHeader>

                <CardContent className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium mb-2">Название</label>
                        <Input
                            value={title}
                            onChange={(e) => setTitle(e.target.value)}
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium mb-2">Описание</label>
                        <Textarea
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            rows={5}
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium mb-2">Место проведения</label>
                        <Input
                            value={location}
                            onChange={(e) => setLocation(e.target.value)}
                            placeholder="Адрес или название места"
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium mb-2">Онлайн-ссылка</label>
                        <Input
                            value={onlineUrl}
                            onChange={(e) => setOnlineUrl(e.target.value)}
                            placeholder="https://..."
                        />
                    </div>

                    <div className="flex gap-3 pt-4">
                        <Button onClick={handleSave} disabled={updateEvent.isPending}>
                            Сохранить
                        </Button>

                        {event?.status === 'draft' && (
                            <Button
                                variant="default"
                                onClick={handlePublish}
                                disabled={publishEvent.isPending}
                            >
                                Опубликовать
                            </Button>
                        )}
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}
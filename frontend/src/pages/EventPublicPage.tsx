import { useEffect, useMemo, useRef, useState } from 'react'
import { useParams, useSearchParams } from 'react-router-dom'
import { Calendar, MapPin, Link as LinkIcon, Users } from 'lucide-react'
import { format } from 'date-fns'
import { ru } from 'date-fns/locale'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import RegistrationForm from '@/widgets/RegistrationForm/RegistrationForm'
import { useEvent } from '@/entities/event/api'
import { useUpdateRSVP } from '@/entities/registration/api'
import { useToast } from '@/components/ui/use-toast'
import fetcher from '@/shared/api/fetcher'
import { useEventPolls } from '@/entities/poll/api'
import PollCard from '@/widgets/PollCard/PollCard'

export default function EventPublicPage() {
    const { eventId } = useParams<{ eventId: string }>()
    const [searchParams] = useSearchParams()
    const formRef = useRef<HTMLDivElement | null>(null)
    const { toast } = useToast()

    const { data: event, isLoading } = useEvent(eventId || '')
    const updateRSVP = useUpdateRSVP(eventId || '')
    const { data: polls } = useEventPolls(eventId || '')

    const [timeLeft, setTimeLeft] = useState<string | null>(null)
    const [isGoogleLoading, setIsGoogleLoading] = useState(false)

    const startDate = useMemo(() => {
        if (!event?.starts_at) return null
        return new Date(event.starts_at)
    }, [event?.starts_at])

    useEffect(() => {
        if (!startDate) {
            setTimeLeft(null)
            return
        }
        const update = () => {
            const now = new Date()
            const diff = startDate.getTime() - now.getTime()
            if (diff <= 0) {
                setTimeLeft(null)
                return
            }
            const totalMinutes = Math.floor(diff / (1000 * 60))
            const hours = Math.floor(totalMinutes / 60)
            const minutes = totalMinutes % 60
            const days = Math.floor(hours / 24)
            const hoursRemainder = hours % 24
            if (days > 0) {
                setTimeLeft(`${days} д ${hoursRemainder} ч`)
            } else {
                setTimeLeft(`${hoursRemainder} ч ${minutes} мин`)
            }
        }
        update()
        const id = window.setInterval(update, 60000)
        return () => window.clearInterval(id)
    }, [startDate])

    useEffect(() => {
        const focusForm = searchParams.get('focus')
        if (focusForm === 'registration' && formRef.current) {
            formRef.current.scrollIntoView({ behavior: 'smooth', block: 'start' })
        }
    }, [searchParams, formRef, isLoading])

    if (!eventId) {
        return <div className="flex items-center justify-center min-h-screen">Нет ID события</div>
    }

    if (isLoading) {
        return <div className="flex items-center justify-center min-h-screen">Загрузка события...</div>
    }

    if (!event) {
        return <div className="flex items-center justify-center min-h-screen">Событие не найдено</div>
    }

    const formattedDate = startDate
        ? format(startDate, 'dd MMMM yyyy, HH:mm', { locale: ru })
        : ''

    const statusColor =
        event.status === 'published'
            ? 'bg-green-100 text-green-800'
            : event.status === 'cancelled'
                ? 'bg-red-100 text-red-800'
                : 'bg-yellow-100 text-yellow-800'

    const handleAddToGoogle = async () => {
        try {
            setIsGoogleLoading(true)
            const data = await fetcher<{ link: string }>(`/api/v1/events/${eventId}/google-calendar`)
            if (data.link) {
                window.open(data.link, '_blank', 'noopener,noreferrer')
            }
        } catch {
            toast({
                title: 'Ошибка',
                description: 'Не удалось открыть Google Календарь',
                variant: 'destructive',
            })
        } finally {
            setIsGoogleLoading(false)
        }
    }

    const handleRSVP = async (status: 'going' | 'not_going' | 'maybe') => {
        try {
            await updateRSVP.mutateAsync(status)
            toast({
                title: 'Ответ сохранён',
            })
        } catch {
            toast({
                title: 'Ошибка',
                description: 'Не удалось обновить статус участия',
                variant: 'destructive',
            })
        }
    }

    const isCancelled = event.status === 'cancelled'
    const isDraft = event.status === 'draft'

    return (
        <div className="container mx-auto py-8 px-4 max-w-4xl space-y-8">
            <Card>
                <CardHeader className="space-y-2">
                    <div className="flex items-center justify-between gap-2">
                        <CardTitle className="text-3xl">{event.title}</CardTitle>
                        <Badge className={statusColor}>{event.status}</Badge>
                    </div>
                    <p className="text-muted-foreground text-base">{event.description}</p>
                </CardHeader>
                <CardContent className="space-y-6">
                    <div className="space-y-3">
                        <div className="flex items-center gap-2 text-muted-foreground">
                            <Calendar className="w-5 h-5" />
                            <span>{formattedDate}</span>
                        </div>
                        {event.location && (
                            <div className="flex items-center gap-2 text-muted-foreground">
                                <MapPin className="w-5 h-5" />
                                <span>{event.location}</span>
                            </div>
                        )}
                        {event.online_url && (
                            <div className="flex items-center gap-2 text-muted-foreground">
                                <LinkIcon className="w-5 h-5" />
                                <a
                                    href={event.online_url}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="text-blue-600 hover:underline"
                                >
                                    Онлайн-ссылка
                                </a>
                            </div>
                        )}
                        {event.capacity > 0 && (
                            <div className="flex items-center gap-2 text-muted-foreground">
                                <Users className="w-5 h-5" />
                                <span>Вместимость: {event.capacity}</span>
                                {event.waitlist && (
                                    <Badge variant="outline" className="ml-2">
                                        Доступен лист ожидания
                                    </Badge>
                                )}
                            </div>
                        )}
                        {timeLeft && event.status === 'published' && (
                            <div className="mt-2 inline-flex items-center rounded-md bg-secondary px-3 py-1 text-sm text-secondary-foreground">
                                До начала события: {timeLeft}
                            </div>
                        )}
                    </div>

                    {!isCancelled && !isDraft && (
                        <>
                            <div className="flex flex-wrap gap-3 pt-4">
                                <Button
                                    size="lg"
                                    onClick={() => {
                                        if (formRef.current) {
                                            formRef.current.scrollIntoView({ behavior: 'smooth', block: 'start' })
                                        }
                                    }}
                                >
                                    Записаться
                                </Button>
                                <Button
                                    variant="outline"
                                    onClick={() => handleRSVP('going')}
                                    disabled={updateRSVP.isPending}
                                >
                                    ✅ Иду
                                </Button>
                                <Button
                                    variant="outline"
                                    onClick={() => handleRSVP('not_going')}
                                    disabled={updateRSVP.isPending}
                                >
                                    ❌ Не иду
                                </Button>
                                <Button
                                    variant="outline"
                                    onClick={() => handleRSVP('maybe')}
                                    disabled={updateRSVP.isPending}
                                >
                                    ❓ Возможно
                                </Button>
                            </div>
                            <div className="flex flex-wrap gap-3 pt-4">
                                <Button asChild variant="secondary">
                                    <a href={`/api/v1/events/${eventId}/ics`} download>
                                        📅 Скачать ICS
                                    </a>
                                </Button>
                                <Button
                                    variant="secondary"
                                    onClick={handleAddToGoogle}
                                    disabled={isGoogleLoading}
                                >
                                    📆 Добавить в Google Календарь
                                </Button>
                            </div>
                        </>
                    )}

                    {isCancelled && (
                        <div className="mt-4 text-sm text-red-600">
                            Событие отменено организатором
                        </div>
                    )}

                    <div ref={formRef} className="mt-8">
                        {!isCancelled && !isDraft && (
                            <>
                                <h2 className="text-xl font-semibold mb-4">Регистрация</h2>
                                <RegistrationForm eventId={eventId} />
                            </>
                        )}
                    </div>
                </CardContent>
            </Card>

            {polls && polls.length > 0 && (
                <div className="space-y-4">
                    <h2 className="text-2xl font-semibold">Опросы</h2>
                    <div className="space-y-3">
                        {polls.map((p) => (
                            <PollCard
                                key={p.id}
                                poll={{
                                    id: p.id,
                                    question: p.question,
                                    options: p.options,
                                    type: p.type,
                                }}
                            />
                        ))}
                    </div>
                </div>
            )}
        </div>
    )
}

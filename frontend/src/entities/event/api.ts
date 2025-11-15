import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import fetcher from '@/shared/api/fetcher'

export type EventStatus = 'draft' | 'published' | 'cancelled'
export type EventVisibility = 'public' | 'private' | 'by_link' | 'by_membership' | 'request'

export interface Event {
    id: string
    title: string
    description: string
    status: EventStatus
    visibility: EventVisibility
    starts_at: string
    ends_at?: string | null
    timezone?: string | null
    location?: string | null
    online_url?: string | null
    capacity: number
    waitlist: boolean
    settings?: Record<string, unknown>
}

export interface MyEvent {
    event: Event
    registration_status: 'going' | 'not_going' | 'maybe' | 'waitlist'
    checked_in: boolean
}

export function useEvent(eventId: string) {
    return useQuery({
        queryKey: ['event', eventId],
        queryFn: () => fetcher<Event>(`/api/v1/events/${eventId}`),
        enabled: !!eventId,
    })
}

export function useCreateEvent() {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (data: { title: string; description: string }) =>
            fetcher<Event>('/api/v1/events', {
                method: 'POST',
                body: JSON.stringify(data),
            }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['my-organized-events'] })
        },
    })
}

export function useUpdateEvent(eventId: string) {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (data: Partial<Event>) =>
            fetcher<{ status: string }>(`/api/v1/events/${eventId}`, {
                method: 'PUT',
                body: JSON.stringify(data),
            }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['event', eventId] })
            queryClient.invalidateQueries({ queryKey: ['my-organized-events'] })
        },
    })
}

export function usePublishEvent(eventId: string) {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: () =>
            fetcher<{ status: string }>(`/api/v1/events/${eventId}/publish`, {
                method: 'POST',
            }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['event', eventId] })
            queryClient.invalidateQueries({ queryKey: ['my-organized-events'] })
        },
    })
}

export function useCancelEvent(eventId: string) {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: () =>
            fetcher<{ status: string }>(`/api/v1/events/${eventId}/cancel`, {
                method: 'POST',
            }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['event', eventId] })
            queryClient.invalidateQueries({ queryKey: ['my-organized-events'] })
        },
    })
}

export function useMyEvents() {
    return useQuery({
        queryKey: ['my-events'],
        queryFn: () => fetcher<MyEvent[]>('/api/v1/me/events'),
    })
}

export function useMyOrganizedEvents() {
    return useQuery({
        queryKey: ['my-organized-events'],
        queryFn: () => fetcher<Event[]>('/api/v1/me/organized-events'),
    })
}

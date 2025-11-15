import { useMutation, useQueryClient } from '@tanstack/react-query'
import fetcher from '@/shared/api/fetcher'

export function useRegisterForEvent(eventId: string) {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (data: { source?: string; utm?: Record<string, unknown> }) =>
            fetcher(`/api/v1/events/${eventId}/register`, {
                method: 'POST',
                body: JSON.stringify({
                    source: data.source ?? 'web',
                    utm: data.utm ?? {},
                }),
            }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['event', eventId] })
            queryClient.invalidateQueries({ queryKey: ['my-events'] })
        },
    })
}

export function useUpdateRSVP(eventId: string) {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (status: 'going' | 'not_going' | 'maybe') =>
            fetcher(`/api/v1/events/${eventId}/rsvp`, {
                method: 'POST',
                body: JSON.stringify({ status }),
            }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['event', eventId] })
            queryClient.invalidateQueries({ queryKey: ['my-events'] })
        },
    })
}

export function useCancelRegistration(eventId: string) {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: () =>
            fetcher(`/api/v1/events/${eventId}/register`, {
                method: 'DELETE',
            }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['event', eventId] })
            queryClient.invalidateQueries({ queryKey: ['my-events'] })
        },
    })
}

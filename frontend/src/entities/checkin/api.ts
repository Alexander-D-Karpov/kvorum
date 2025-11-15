import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import fetcher from '@/shared/api/fetcher'

export interface CheckinResponse {
    id: string
    event_id: string
    user_id: string
    method: 'qr' | 'manual'
    at: string
}

export function useScanCheckin(eventId: string) {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (qrCode: string) =>
            fetcher<CheckinResponse>(`/api/v1/events/${eventId}/checkin/scan`, {
                method: 'POST',
                body: JSON.stringify({ qr_code: qrCode }),
            }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['event-checkins', eventId] })
        },
    })
}

export function useManualCheckin(eventId: string) {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (userId: string) =>
            fetcher<CheckinResponse>(`/api/v1/events/${eventId}/checkin/manual`, {
                method: 'POST',
                body: JSON.stringify({ user_id: userId }),
            }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['event-checkins', eventId] })
        },
    })
}

export function useTicketQRCode(eventId: string) {
    return useQuery({
        queryKey: ['ticket-qr', eventId],
        queryFn: () => fetcher<{ token: string }>(`/api/v1/tickets/${eventId}/qr`),
        enabled: !!eventId,
    })
}

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import fetcher from '@/shared/api/fetcher'

export interface Campaign {
    id: string
    event_id: string
    name: string
    segment: string
    channel: string
    message: string
    scheduled_at?: string
    status: string
}

export function useCampaigns(eventId: string) {
    return useQuery({
        queryKey: ['campaigns', eventId],
        queryFn: () => fetcher<Campaign[]>(`/api/v1/events/${eventId}/campaigns`),
        enabled: !!eventId,
    })
}

export function useCreateCampaign(eventId: string) {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (data: {
            name: string
            segment: string
            channel: string
            message: string
            scheduled_at?: string
        }) =>
            fetcher<Campaign>(`/api/v1/events/${eventId}/campaigns`, {
                method: 'POST',
                body: JSON.stringify(data),
            }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['campaigns', eventId] })
        },
    })
}
import { useQuery } from '@tanstack/react-query'
import fetcher from '@/shared/api/fetcher'

export interface EventAnalytics {
    event_id: string
    period_from: string
    period_to: string
    total_registrations: number
    going: number
    not_going: number
    maybe: number
    waitlist: number
    checked_in: number
    by_source: Record<string, number>
}

export function useEventAnalytics(eventId: string, from?: string, to?: string) {
    return useQuery({
        queryKey: ['analytics', eventId, { from, to }],
        queryFn: () => {
            const params = new URLSearchParams()
            if (from) params.set('from', from)
            if (to) params.set('to', to)
            const query = params.toString()
            const url = query
                ? `/api/v1/events/${eventId}/analytics?${query}`
                : `/api/v1/events/${eventId}/analytics`
            return fetcher<EventAnalytics>(url)
        },
        enabled: !!eventId,
    })
}

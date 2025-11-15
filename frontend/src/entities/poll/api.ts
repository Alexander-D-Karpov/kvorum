import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import fetcher from '@/shared/api/fetcher'

export interface Poll {
    id: string
    event_id: string
    question: string
    options: Record<string, string>
    type: 'single' | 'multiple' | 'rating' | 'nps'
    created_at: string
    updated_at: string
}

export function useEventPolls(eventId: string) {
    return useQuery({
        queryKey: ['polls', eventId],
        queryFn: () => fetcher<Poll[]>(`/api/v1/events/${eventId}/polls`),
        enabled: !!eventId,
    })
}

export function useCreatePoll(eventId: string) {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (data: { question: string; options: Record<string, string>; type: string }) =>
            fetcher<Poll>(`/api/v1/events/${eventId}/polls`, {
                method: 'POST',
                body: JSON.stringify(data),
            }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['polls', eventId] })
        },
    })
}

export function useVote(pollId: string) {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: (optionKey: string) =>
            fetcher(`/api/v1/polls/${pollId}/vote`, {
                method: 'POST',
                body: JSON.stringify({ option_key: optionKey }),
            }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['poll', pollId, 'results'] })
        },
    })
}

export function usePollResults(pollId: string) {
    return useQuery({
        queryKey: ['poll', pollId, 'results'],
        queryFn: () => fetcher<Record<string, number>>(`/api/v1/polls/${pollId}/results`),
        enabled: !!pollId,
    })
}

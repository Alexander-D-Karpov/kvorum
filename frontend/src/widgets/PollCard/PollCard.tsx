import { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { useVote, usePollResults } from '@/entities/poll/api'
import { useToast } from '@/components/ui/use-toast'

interface PollCardProps {
    poll: {
        id: string
        question: string
        options: Record<string, string>
        type: string
    }
}

export default function PollCard({ poll }: PollCardProps) {
    const [selectedOption, setSelectedOption] = useState<string>('')
    const [hasVoted, setHasVoted] = useState(false)
    const vote = useVote(poll.id)
    const { data: results } = usePollResults(poll.id)
    const { toast } = useToast()

    const totalVotes = results ? Object.values(results).reduce((a, b) => a + b, 0) : 0

    const handleVote = async () => {
        if (!selectedOption) return
        try {
            await vote.mutateAsync(selectedOption)
            setHasVoted(true)
            toast({
                title: 'Спасибо за ваш голос',
            })
        } catch {
            toast({
                title: 'Ошибка',
                description: 'Не удалось проголосовать',
                variant: 'destructive',
            })
        }
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle>{poll.question}</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
                {!hasVoted ? (
                    <>
                        <div className="space-y-2">
                            {Object.entries(poll.options).map(([key, label]) => (
                                <label key={key} className="flex items-center gap-2 text-sm cursor-pointer">
                                    <input
                                        type="radio"
                                        name={poll.id}
                                        value={key}
                                        checked={selectedOption === key}
                                        onChange={() => setSelectedOption(key)}
                                    />
                                    <span>{label}</span>
                                </label>
                            ))}
                        </div>
                        <Button onClick={handleVote} disabled={!selectedOption || vote.isPending}>
                            Проголосовать
                        </Button>
                    </>
                ) : (
                    <div className="space-y-3">
                        <p className="text-sm text-muted-foreground">Всего голосов: {totalVotes}</p>
                        {results &&
                            Object.entries(poll.options).map(([key, label]) => {
                                const count = results[key] || 0
                                const percentage = totalVotes > 0 ? (count / totalVotes) * 100 : 0
                                return (
                                    <div key={key} className="space-y-1">
                                        <div className="flex justify-between text-sm">
                                            <span>{label}</span>
                                            <span className="text-muted-foreground">
                        {count} ({percentage.toFixed(0)}%)
                      </span>
                                        </div>
                                        <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
                                            <div
                                                className="h-full bg-blue-600 transition-all"
                                                style={{ width: `${percentage}%` }}
                                            />
                                        </div>
                                    </div>
                                )
                            })}
                    </div>
                )}
            </CardContent>
        </Card>
    )
}

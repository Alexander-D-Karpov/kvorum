import { useState } from "react";
import { useParams } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import fetcher from "@/shared/api/fetcher";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/components/ui/use-toast";

interface Campaign {
    id: string;
    name: string;
    status: string;
    scheduled_at?: string | null;
    sent_at?: string | null;
    channel: string;
    segment: string;
}

export default function CampaignsPage() {
    const { eventId } = useParams<{ eventId: string }>();
    const { toast } = useToast();
    const queryClient = useQueryClient();

    const { data: campaigns } = useQuery({
        queryKey: ["campaigns", eventId],
        queryFn: () =>
            fetcher<Campaign[]>(`/api/v1/events/${eventId}/campaigns`),
        enabled: !!eventId,
    });

    const [name, setName] = useState("");
    const [segment, setSegment] = useState("all");
    const [channel, setChannel] = useState("bot");
    const [message, setMessage] = useState("");
    const [scheduledAt, setScheduledAt] = useState("");

    const createMutation = useMutation({
        mutationFn: () =>
            fetcher<Campaign>(`/api/v1/events/${eventId}/campaigns`, {
                method: "POST",
                body: JSON.stringify({
                    name,
                    segment,
                    channel,
                    message,
                    scheduled_at: scheduledAt || null,
                }),
            }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["campaigns", eventId] });
            setName("");
            setMessage("");
            setScheduledAt("");
            toast({
                title: "Рассылка сохранена",
            });
        },
    });

    if (!eventId) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                Нет ID события
            </div>
        );
    }

    return (
        <div className="container mx-auto max-w-5xl px-4 py-8 space-y-6">
            <h1 className="text-2xl font-bold">Рассылки по событию</h1>

            <Card>
                <CardHeader>
                    <CardTitle className="text-base">Новая рассылка</CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                    <Input
                        placeholder="Название кампании"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                    />
                    <div className="grid md:grid-cols-3 gap-3">
                        <div className="space-y-1">
                            <p className="text-xs text-muted-foreground">Сегмент</p>
                            <select
                                className="w-full border rounded-md px-3 py-2 text-sm bg-background"
                                value={segment}
                                onChange={(e) => setSegment(e.target.value)}
                            >
                                <option value="all">Все зарегистрированные</option>
                                <option value="going">Идут</option>
                                <option value="not_going">Не идут</option>
                                <option value="maybe">Возможно</option>
                                <option value="waitlist">Лист ожидания</option>
                            </select>
                        </div>
                        <div className="space-y-1">
                            <p className="text-xs text-muted-foreground">Канал</p>
                            <select
                                className="w-full border rounded-md px-3 py-2 text-sm bg-background"
                                value={channel}
                                onChange={(e) => setChannel(e.target.value)}
                            >
                                <option value="bot">Бот</option>
                                <option value="email">Email</option>
                            </select>
                        </div>
                        <div className="space-y-1">
                            <p className="text-xs text-muted-foreground">Время отправки</p>
                            <Input
                                type="datetime-local"
                                value={scheduledAt}
                                onChange={(e) => setScheduledAt(e.target.value)}
                            />
                        </div>
                    </div>
                    <Textarea
                        placeholder="Текст сообщения"
                        value={message}
                        onChange={(e) => setMessage(e.target.value)}
                    />
                    <Button
                        onClick={() => createMutation.mutate()}
                        disabled={
                            !name.trim() || !message.trim() || createMutation.isPending
                        }
                    >
                        Сохранить кампанию
                    </Button>
                </CardContent>
            </Card>

            <Card>
                <CardHeader>
                    <CardTitle className="text-base">История рассылок</CardTitle>
                </CardHeader>
                <CardContent className="space-y-2 text-sm">
                    {!campaigns || campaigns.length === 0 ? (
                        <div className="text-muted-foreground">Рассылок пока нет</div>
                    ) : (
                        campaigns.map((c) => (
                            <div
                                key={c.id}
                                className="flex items-center justify-between border-b last:border-b-0 py-2"
                            >
                                <div>
                                    <div className="font-medium">{c.name}</div>
                                    <div className="text-xs text-muted-foreground">
                                        Канал: {c.channel}, сегмент: {c.segment}
                                    </div>
                                </div>
                                <div className="text-xs text-muted-foreground text-right">
                                    <div>Статус: {c.status}</div>
                                    {c.scheduled_at && (
                                        <div>Запланировано: {c.scheduled_at}</div>
                                    )}
                                    {c.sent_at && <div>Отправлено: {c.sent_at}</div>}
                                </div>
                            </div>
                        ))
                    )}
                </CardContent>
            </Card>
        </div>
    );
}

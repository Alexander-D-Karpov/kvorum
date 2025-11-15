import { useMemo } from "react";
import { useMyEvents } from "@/entities/event/api";
import {
    Tabs,
    TabsList,
    TabsTrigger,
    TabsContent,
} from "@/components/ui/tabs";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useTicketQRCode } from "@/entities/checkin/api";
import QRCode from "@/widgets/QRCode/QRCode";

function EventsList({ type }: { type: "upcoming" | "past" }) {
    const { data, isLoading } = useMyEvents();
    const now = new Date();

    const items = useMemo(() => {
        if (!data) return [];
        return data.filter((item) => {
            const startsAt = new Date(item.event.starts_at);
            if (type === "upcoming") {
                return startsAt >= now;
            }
            return startsAt < now;
        });
    }, [data, type, now]);

    if (isLoading) {
        return (
            <div className="py-4 text-sm text-muted-foreground">Загрузка...</div>
        );
    }

    if (!items.length) {
        return (
            <div className="py-4 text-sm text-muted-foreground">Список пуст</div>
        );
    }

    return (
        <div className="space-y-3">
            {items.map((item) => (
                <Card key={item.event.id}>
                    <CardHeader className="flex flex-row items-center justify-between gap-3">
                        <div>
                            <CardTitle className="text-base">{item.event.title}</CardTitle>
                            <p className="text-xs text-muted-foreground">
                                {new Date(item.event.starts_at).toLocaleString("ru-RU")}
                            </p>
                        </div>
                        <div className="text-xs text-muted-foreground text-right">
                            <p>Статус: {item.registration_status}</p>
                            {item.checked_in && <p>Билет использован на входе</p>}
                        </div>
                    </CardHeader>
                </Card>
            ))}
        </div>
    );
}

function TicketsTab() {
    const { data, isLoading } = useMyEvents();

    if (isLoading) {
        return (
            <div className="py-4 text-sm text-muted-foreground">Загрузка...</div>
        );
    }

    if (!data || !data.length) {
        return (
            <div className="py-4 text-sm text-muted-foreground">
                Нет активных билетов
            </div>
        );
    }

    return (
        <div className="space-y-4">
            {data
                .filter((item) => item.registration_status === "going")
                .map((item) => (
                    <TicketCard
                        key={item.event.id}
                        eventId={item.event.id}
                        title={item.event.title}
                    />
                ))}
        </div>
    );
}

function TicketCard({ eventId, title }: { eventId: string; title: string }) {
    const { data, isLoading } = useTicketQRCode(eventId);

    return (
        <Card>
            <CardHeader>
                <CardTitle className="text-base">{title}</CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col items-center gap-3">
                {isLoading && (
                    <div className="text-xs text-muted-foreground">
                        Генерация билета...
                    </div>
                )}
                {data?.token && (
                    <>
                        <QRCode value={data.token} size={160} />
                        <div className="text-[10px] text-muted-foreground break-all">
                            {data.token}
                        </div>
                    </>
                )}
            </CardContent>
        </Card>
    );
}

export default function AttendeeDashboard() {
    return (
        <div className="container mx-auto max-w-4xl px-4 py-8 space-y-6">
            <div className="flex items-center justify-between gap-4">
                <h1 className="text-2xl font-bold">Мои события</h1>
                <Button asChild variant="outline" size="sm">
                    <a href="/">На главную</a>
                </Button>
            </div>
            <Tabs defaultValue="upcoming">
                <TabsList>
                    <TabsTrigger value="upcoming">Предстоящие</TabsTrigger>
                    <TabsTrigger value="past">Прошедшие</TabsTrigger>
                    <TabsTrigger value="tickets">Билеты</TabsTrigger>
                </TabsList>
                <TabsContent value="upcoming" className="pt-4">
                    <EventsList type="upcoming" />
                </TabsContent>
                <TabsContent value="past" className="pt-4">
                    <EventsList type="past" />
                </TabsContent>
                <TabsContent value="tickets" className="pt-4">
                    <TicketsTab />
                </TabsContent>
            </Tabs>
        </div>
    );
}

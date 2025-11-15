import { useParams } from "react-router-dom";
import QRScanner from "@/features/scan-checkin/QRScanner";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

export default function CheckinPage() {
    const { eventId } = useParams<{ eventId: string }>();

    if (!eventId) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                Нет ID события
            </div>
        );
    }

    return (
        <div className="container mx-auto max-w-3xl px-4 py-8">
            <Card>
                <CardHeader>
                    <CardTitle>Чек-ин события</CardTitle>
                </CardHeader>
                <CardContent>
                    <QRScanner eventId={eventId} />
                </CardContent>
            </Card>
        </div>
    );
}

import { useEffect, useRef, useState } from "react";
import { BrowserMultiFormatReader } from "@zxing/library";
import { useScanCheckin } from "@/entities/checkin/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { useToast } from "@/components/ui/use-toast";
import { Camera } from "lucide-react";
import { APIError } from "@/shared/api/fetcher";

interface QRScannerProps {
    eventId: string;
}

function loadQueue(eventId: string): string[] {
    if (typeof window === "undefined") return [];
    try {
        const raw = window.localStorage.getItem(`checkin_queue_${eventId}`);
        if (!raw) return [];
        const parsed = JSON.parse(raw);
        if (Array.isArray(parsed)) return parsed.filter((x) => typeof x === "string");
        return [];
    } catch {
        return [];
    }
}

function saveQueue(eventId: string, queue: string[]) {
    if (typeof window === "undefined") return;
    try {
        window.localStorage.setItem(`checkin_queue_${eventId}`, JSON.stringify(queue));
    } catch {
    }
}

export default function QRScanner({ eventId }: QRScannerProps) {
    const videoRef = useRef<HTMLVideoElement | null>(null);
    const readerRef = useRef<BrowserMultiFormatReader | null>(null);
    const [isScanning, setIsScanning] = useState(false);
    const [successCount, setSuccessCount] = useState(0);
    const [errorCount, setErrorCount] = useState(0);
    const [lastToken, setLastToken] = useState<string | null>(null);
    const flushingRef = useRef(false);

    const scanCheckin = useScanCheckin(eventId);
    const { toast } = useToast();

    const flushQueue = async () => {
        if (!eventId) return;
        if (flushingRef.current) return;
        const queue = loadQueue(eventId);
        if (queue.length === 0) return;
        flushingRef.current = true;
        let processed = 0;
        const remaining: string[] = [];
        for (const code of queue) {
            try {
                await scanCheckin.mutateAsync(code);
                processed += 1;
            } catch {
                remaining.push(code);
            }
        }
        saveQueue(eventId, remaining);
        flushingRef.current = false;
        if (processed > 0) {
            setSuccessCount((x) => x + processed);
            toast({
                title: "Синхронизация чек-инов",
                description: `Отправлено из очереди: ${processed}`,
            });
        }
    };

    const handleScanResult = async (qrCode: string) => {
        if (!qrCode) return;

        if (lastToken && lastToken === qrCode) {
            toast({
                title: "Повторное сканирование",
                description: "Этот QR уже только что сканировали",
            });
            return;
        }

        setLastToken(qrCode);

        try {
            await scanCheckin.mutateAsync(qrCode);
            setSuccessCount((x) => x + 1);
            toast({
                title: "Успешный чек-ин",
                description: "Участник отмечен",
            });
        } catch (error: unknown) {
            if (error instanceof APIError) {
                if (error.status === 409) {
                    toast({
                        title: "Уже отмечен",
                        description: "Этот участник уже проходил чек-ин",
                    });
                    return;
                }
                toast({
                    title: "Ошибка QR",
                    description: "Невалидный QR-код",
                    variant: "destructive",
                });
                setErrorCount((x) => x + 1);
                return;
            }

            const offline =
                typeof navigator !== "undefined" && navigator.onLine === false;
            if (offline) {
                const queue = loadQueue(eventId);
                queue.push(qrCode);
                saveQueue(eventId, queue);
                toast({
                    title: "Офлайн-режим",
                    description: "QR сохранён в очередь, синхронизируем при появлении сети",
                });
                return;
            }

            setErrorCount((x) => x + 1);
            toast({
                title: "Ошибка",
                description: "Не удалось отметить участника",
                variant: "destructive",
            });
        }
    };

    const startScanning = async () => {
        if (!videoRef.current) return;
        try {
            const reader = new BrowserMultiFormatReader();
            readerRef.current = reader;
            const devices = await reader.listVideoInputDevices();
            const selectedDeviceId = devices[0]?.deviceId;
            setIsScanning(true);
            reader.decodeFromVideoDevice(
                selectedDeviceId,
                videoRef.current,
                async (result) => {
                    if (result) {
                        const qrCode = result.getText();
                        await handleScanResult(qrCode);
                    }
                },
            );
        } catch {
            toast({
                title: "Ошибка",
                description: "Не удалось получить доступ к камере",
                variant: "destructive",
            });
        }
    };

    const stopScanning = () => {
        if (readerRef.current) {
            readerRef.current.reset();
        }
        setIsScanning(false);
    };

    useEffect(() => {
        flushQueue();
    }, [eventId]);

    useEffect(() => {
        const handler = () => {
            if (navigator.onLine) {
                flushQueue();
            }
        };
        window.addEventListener("online", handler);
        return () => {
            window.removeEventListener("online", handler);
            stopScanning();
        };
    }, []);

    return (
        <Card>
            <CardContent className="pt-6 space-y-4">
                <video
                    ref={videoRef}
                    className="w-full rounded-lg border"
                    style={{ maxHeight: "400px" }}
                />
                <div className="flex gap-3">
                    {!isScanning ? (
                        <Button onClick={startScanning} className="flex-1">
                            <Camera className="w-4 h-4 mr-2" />
                            Начать сканирование
                        </Button>
                    ) : (
                        <Button onClick={stopScanning} variant="destructive" className="flex-1">
                            Остановить
                        </Button>
                    )}
                </div>
                <div className="flex items-center justify-between text-xs text-muted-foreground">
                    <span>Успешно: {successCount}</span>
                    <span>Ошибки: {errorCount}</span>
                </div>
            </CardContent>
        </Card>
    );
}

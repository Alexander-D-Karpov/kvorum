import { useEffect, useRef } from 'react'
import { useTicketQRCode } from '@/entities/checkin/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import QRCodeStyling from 'qr-code-styling'

interface Props {
    eventId: string
}

export default function QRCode({ eventId }: Props) {
    const { data: qrData, isLoading } = useTicketQRCode(eventId)
    const ref = useRef<HTMLDivElement>(null)

    useEffect(() => {
        if (!qrData?.token || !ref.current) return

        const qrCode = new QRCodeStyling({
            width: 300,
            height: 300,
            data: qrData.token,
            dotsOptions: {
                color: '#000000',
                type: 'rounded',
            },
            backgroundOptions: {
                color: '#ffffff',
            },
        })

        ref.current.innerHTML = ''
        qrCode.append(ref.current)
    }, [qrData])

    if (isLoading) {
        return <div>Загрузка QR-кода...</div>
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle>Ваш QR-код для входа</CardTitle>
            </CardHeader>
            <CardContent className="flex justify-center">
                <div ref={ref} />
            </CardContent>
        </Card>
    )
}
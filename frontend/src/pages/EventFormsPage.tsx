import { useState } from "react";
import { useParams } from "react-router-dom";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
    useActiveForm,
    useCreateForm,
    FormField,
    FieldRule,
} from "@/entities/form/api";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/components/ui/use-toast";

export default function EventFormsPage() {
    const { eventId } = useParams<{ eventId: string }>();
    const { data: currentForm } = useActiveForm(eventId || "");
    const createForm = useCreateForm(eventId || "");
    const { toast } = useToast();

    const [fields, setFields] = useState<FormField[]>(
        currentForm?.schema.fields ?? [],
    );
    const [rules, setRules] = useState<FieldRule[]>(currentForm?.rules ?? []);

    const [newFieldLabel, setNewFieldLabel] = useState("");
    const [newFieldType, setNewFieldType] = useState<
        "text" | "textarea" | "select" | "checkbox" | "radio"
    >("text");

    const addField = () => {
        if (!newFieldLabel.trim()) return;
        const id = newFieldLabel.toLowerCase().replace(/\s+/g, "_");
        setFields((prev) => [...prev, { id, label: newFieldLabel, type: newFieldType }]);
        setNewFieldLabel("");
    };

    const save = async () => {
        try {
            await createForm.mutateAsync({
                schema: { fields },
                rules,
            });
            toast({
                title: "Форма сохранена",
            });
        } catch {
            toast({
                title: "Ошибка",
                description: "Не удалось сохранить форму",
                variant: "destructive",
            });
        }
    };

    if (!eventId) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                Нет ID события
            </div>
        );
    }

    return (
        <div className="container mx-auto max-w-5xl px-4 py-8 space-y-6">
            <h1 className="text-2xl font-bold">Форма регистрации</h1>

            <Card>
                <CardHeader>
                    <CardTitle className="text-base">Поля формы</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="grid md:grid-cols-[1fr_auto_auto] gap-2 items-center">
                        <Input
                            placeholder="Название поля"
                            value={newFieldLabel}
                            onChange={(e) => setNewFieldLabel(e.target.value)}
                        />
                        <select
                            className="border rounded-md px-3 py-2 text-sm bg-background"
                            value={newFieldType}
                            onChange={(e) =>
                                setNewFieldType(e.target.value as
                                    | "text"
                                    | "textarea"
                                    | "select"
                                    | "checkbox"
                                    | "radio")
                            }
                        >
                            <option value="text">Текст</option>
                            <option value="textarea">Многострочный текст</option>
                            <option value="select">Список</option>
                            <option value="checkbox">Чекбокс</option>
                            <option value="radio">Радио</option>
                        </select>
                        <Button onClick={addField} disabled={!newFieldLabel.trim()}>
                            Добавить
                        </Button>
                    </div>

                    <div className="space-y-2 text-sm">
                        {fields.map((field, index) => (
                            <div
                                key={field.id}
                                className="flex flex-wrap items-center gap-2 border rounded-md px-3 py-2"
                            >
                                <span className="font-medium">{field.label}</span>
                                <span className="text-xs text-muted-foreground">
                  ({field.type})
                </span>
                                <Button
                                    size="sm"
                                    variant="ghost"
                                    onClick={() =>
                                        setFields((prev) => prev.filter((_, i) => i !== index))
                                    }
                                >
                                    Удалить
                                </Button>
                            </div>
                        ))}
                    </div>
                </CardContent>
            </Card>

            <Card>
                <CardHeader>
                    <CardTitle className="text-base">Правила отображения</CardTitle>
                </CardHeader>
                <CardContent className="space-y-2 text-sm">
                    <p className="text-xs text-muted-foreground">
                        Для простоты здесь можно редактировать JSON правил вручную
                    </p>
                    <Textarea
                        rows={8}
                        value={JSON.stringify(rules, null, 2)}
                        onChange={(e) => {
                            try {
                                const parsed = JSON.parse(e.target.value);
                                setRules(parsed);
                            } catch {
                            }
                        }}
                    />
                </CardContent>
            </Card>

            <Button onClick={save}>Сохранить форму</Button>
        </div>
    );
}

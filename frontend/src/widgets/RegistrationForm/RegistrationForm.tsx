import { useEffect, useMemo, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import {
    useActiveForm,
    useDraft,
    useSaveDraft,
    useSubmitForm,
    FieldRule,
} from "@/entities/form/api";
import { useAuth } from "@/shared/providers/AuthProvider";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/components/ui/use-toast";

interface RegistrationFormProps {
    eventId: string;
    onSuccess?: () => void;
}

type FormValues = Record<string, any>;

export default function RegistrationForm({
                                             eventId,
                                             onSuccess,
                                         }: RegistrationFormProps) {
    const { user } = useAuth();
    const { toast } = useToast();

    const { data: form, isLoading: formLoading } = useActiveForm(eventId);
    const formId = form?.id ?? "";
    const { data: draftData } = useDraft(formId);
    const saveDraft = useSaveDraft(formId);
    const submitForm = useSubmitForm(formId);

    const {
        register,
        handleSubmit,
        watch,
        reset,
        formState: { errors, isSubmitting },
    } = useForm<FormValues>({ mode: "onBlur" });

    const [hiddenFields, setHiddenFields] = useState<Record<string, boolean>>({});
    const [requiredOverrides, setRequiredOverrides] = useState<
        Record<string, boolean>
    >({});
    const [lastSavedAt, setLastSavedAt] = useState<Date | null>(null);

    const fields = useMemo(() => form?.schema?.fields ?? [], [form]);
    const values = watch();
    const autoSaveRef = useRef<number | null>(null);

    useEffect(() => {
        if (!form) return;
        const initialValues: FormValues = {};
        if (draftData && draftData.draft) {
            Object.assign(initialValues, draftData.draft);
        } else if (user) {
            for (const field of fields) {
                if (!field.id) continue;
                if (!initialValues[field.id]) {
                    if (field.id === "email" && user.email) initialValues[field.id] = user.email;
                    if (
                        (field.id === "name" || field.id === "full_name") &&
                        user.display_name
                    ) {
                        initialValues[field.id] = user.display_name;
                    }
                    if (field.id === "phone" && user.phone) initialValues[field.id] = user.phone;
                }
            }
        }
        reset(initialValues);
    }, [form, draftData, user, fields, reset]);

    useEffect(() => {
        if (!form?.rules || form.rules.length === 0) return;
        const newHidden: Record<string, boolean> = {};
        const newRequired: Record<string, boolean> = {};
        const rules = form.rules as FieldRule[];
        for (const rule of rules) {
            const satisfied = rule.when.every((cond) => values[cond.field] === cond.equals);
            if (!satisfied) continue;
            if (rule.action === "hide") newHidden[rule.target] = true;
            if (rule.action === "show") newHidden[rule.target] = false;
            if (rule.action === "require") newRequired[rule.target] = true;
            if (rule.action === "optional") newRequired[rule.target] = false;
        }
        setHiddenFields(newHidden);
        setRequiredOverrides(newRequired);
    }, [values, form?.rules]);

    useEffect(() => {
        if (!formId) return;
        if (autoSaveRef.current) {
            window.clearTimeout(autoSaveRef.current);
        }
        autoSaveRef.current = window.setTimeout(() => {
            if (!values || Object.keys(values).length === 0) return;
            saveDraft.mutate(values, {
                onSuccess: () => {
                    setLastSavedAt(new Date());
                },
            });
        }, 800);
        return () => {
            if (autoSaveRef.current) {
                window.clearTimeout(autoSaveRef.current);
            }
        };
    }, [values, formId, saveDraft]);

    const onSubmit = async (data: FormValues) => {
        if (!formId) return;
        try {
            await submitForm.mutateAsync(data);
            toast({
                title: "Регистрация отправлена",
            });
            if (onSuccess) onSuccess();
        } catch {
            toast({
                title: "Ошибка",
                description: "Не удалось отправить форму",
                variant: "destructive",
            });
        }
    };

    if (formLoading || !form) {
        return (
            <div className="py-6 text-sm text-muted-foreground">Загрузка формы...</div>
        );
    }

    return (
        <form className="space-y-6" onSubmit={handleSubmit(onSubmit)}>
            {fields.map((field) => {
                if (hiddenFields[field.id]) {
                    return null;
                }
                const required = requiredOverrides[field.id] ?? field.required ?? false;

                if (field.type === "textarea") {
                    return (
                        <div key={field.id} className="space-y-1">
                            <label className="text-sm font-medium">
                                {field.label}
                                {required && <span className="text-red-500 ml-1">*</span>}
                            </label>
                            <Textarea
                                {...register(field.id, { required })}
                                placeholder={field.placeholder}
                            />
                            {errors[field.id] && (
                                <p className="text-xs text-red-500">Это поле обязательно</p>
                            )}
                        </div>
                    );
                }

                if (field.type === "select" && field.options) {
                    return (
                        <div key={field.id} className="space-y-1">
                            <label className="text-sm font-medium">
                                {field.label}
                                {required && <span className="text-red-500 ml-1">*</span>}
                            </label>
                            <select
                                {...register(field.id, { required })}
                                className="w-full border rounded-md px-3 py-2 text-sm bg-background"
                            >
                                <option value="">
                                    {field.placeholder || "Выберите вариант"}
                                </option>
                                {field.options.map((opt) => (
                                    <option key={opt.value} value={opt.value}>
                                        {opt.label}
                                    </option>
                                ))}
                            </select>
                            {errors[field.id] && (
                                <p className="text-xs text-red-500">Это поле обязательно</p>
                            )}
                        </div>
                    );
                }

                if (field.type === "checkbox") {
                    return (
                        <div key={field.id} className="flex items-center gap-2">
                            <input
                                type="checkbox"
                                {...register(field.id, { required })}
                                className="h-4 w-4"
                            />
                            <span className="text-sm">
                {field.label}
                                {required && <span className="text-red-500 ml-1">*</span>}
              </span>
                            {errors[field.id] && (
                                <p className="text-xs text-red-500">Обязательно для продолжения</p>
                            )}
                        </div>
                    );
                }

                if (field.type === "radio" && field.options) {
                    return (
                        <div key={field.id} className="space-y-1">
                            <p className="text-sm font-medium">
                                {field.label}
                                {required && <span className="text-red-500 ml-1">*</span>}
                            </p>
                            <div className="space-y-2">
                                {field.options.map((opt) => (
                                    <label
                                        key={opt.value}
                                        className="flex items-center gap-2 text-sm cursor-pointer"
                                    >
                                        <input
                                            type="radio"
                                            value={opt.value}
                                            {...register(field.id, { required })}
                                        />
                                        <span>{opt.label}</span>
                                    </label>
                                ))}
                            </div>
                            {errors[field.id] && (
                                <p className="text-xs text-red-500">Выберите вариант</p>
                            )}
                        </div>
                    );
                }

                return (
                    <div key={field.id} className="space-y-1">
                        <label className="text-sm font-medium">
                            {field.label}
                            {required && <span className="text-red-500 ml-1">*</span>}
                        </label>
                        <Input
                            {...register(field.id, { required })}
                            placeholder={field.placeholder}
                        />
                        {errors[field.id] && (
                            <p className="text-xs text-red-500">Это поле обязательно</p>
                        )}
                    </div>
                );
            })}

            <div className="flex items-center justify-between gap-2 text-xs text-muted-foreground">
                {lastSavedAt && (
                    <span>
            Черновик сохранён в{" "}
                        {lastSavedAt.toLocaleTimeString("ru-RU", {
                            hour: "2-digit",
                            minute: "2-digit",
                        })}
          </span>
                )}
                {saveDraft.isPending && <span>Сохранение черновика...</span>}
            </div>

            <Button type="submit" disabled={isSubmitting}>
                Отправить
            </Button>
        </form>
    );
}

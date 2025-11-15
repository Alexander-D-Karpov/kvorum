import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import fetcher from "@/shared/api/fetcher";

export type FormFieldType = "text" | "textarea" | "select" | "checkbox" | "radio";

export interface FormFieldOption {
    value: string;
    label: string;
}

export interface FormField {
    id: string;
    label: string;
    type: FormFieldType;
    placeholder?: string;
    required?: boolean;
    options?: FormFieldOption[];
}

export interface FieldCondition {
    field: string;
    equals: string | number | boolean;
}

export type FieldRuleAction = "show" | "hide" | "require" | "optional";

export interface FieldRule {
    target: string;
    action: FieldRuleAction;
    when: FieldCondition[];
}

export interface Form {
    id: string;
    event_id: string;
    version: number;
    schema: {
        fields: FormField[];
    };
    rules: FieldRule[];
}

export function useActiveForm(eventId: string) {
    return useQuery({
        queryKey: ["form", "active", eventId],
        queryFn: () => fetcher<Form>(`/api/v1/events/${eventId}/forms/active`),
        enabled: !!eventId,
    });
}

export function useSubmitForm(formId: string) {
    return useMutation({
        mutationFn: (answers: Record<string, unknown>) =>
            fetcher(`/api/v1/forms/${formId}/submit`, {
                method: "POST",
                body: JSON.stringify({ answers }),
            }),
    });
}

export function useDraft(formId: string) {
    return useQuery({
        queryKey: ["form", "draft", formId],
        queryFn: () =>
            fetcher<{ draft?: Record<string, unknown> }>(
                `/api/v1/forms/${formId}/draft`,
            ),
        enabled: !!formId,
    });
}

export function useSaveDraft(formId: string) {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: (data: Record<string, unknown>) =>
            fetcher(`/api/v1/forms/${formId}/draft`, {
                method: "PUT",
                body: JSON.stringify({ data }),
            }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["form", "draft", formId] });
        },
    });
}

export function useCreateForm(eventId: string) {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: (payload: { schema: Form["schema"]; rules: Form["rules"] }) =>
            fetcher<Form>(`/api/v1/events/${eventId}/forms`, {
                method: "POST",
                body: JSON.stringify(payload),
            }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["form", "active", eventId] });
        },
    });
}

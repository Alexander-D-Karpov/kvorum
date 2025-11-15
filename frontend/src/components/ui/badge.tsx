import * as React from "react";
import { cn } from "@/lib/utils";

export interface BadgeProps
    extends React.HTMLAttributes<HTMLSpanElement> {
    variant?: "default" | "outline";
}

export function Badge({
                          className,
                          variant = "default",
                          ...props
                      }: BadgeProps) {
    const base =
        "inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-semibold border";
    const variants: Record<string, string> = {
        default: "bg-primary text-primary-foreground border-transparent",
        outline: "border border-border text-foreground",
    };
    return (
        <span
            className={cn(base, variants[variant] ?? variants.default, className)}
            {...props}
        />
    );
}

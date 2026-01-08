"use client";

import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { LucideIcon } from "lucide-react";
import { cn } from "@/lib/utils";

interface MetricCardProps {
  title: string;
  value: string | number;
  icon: LucideIcon;
  trend?: {
    value: number;
    label: string;
    isPositive: boolean;
  };
  isLoading?: boolean;
  className?: string;
  style?: React.CSSProperties;
}

export function MetricCard({
  title,
  value,
  icon: Icon,
  trend,
  isLoading = false,
  className,
  style,
}: MetricCardProps) {
  if (isLoading) {
    return (
      <Card className={cn("hover-lift", className)} style={style}>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium label-muted">
            <Skeleton className="h-4 w-24" />
          </CardTitle>
          <Skeleton className="h-4 w-4" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-8 w-16 mb-2" />
          <Skeleton className="h-3 w-20" />
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className={cn("hover-lift animate-in fade-in slide-in-from-bottom-4", className)} style={style}>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium label-muted">{title}</CardTitle>
        <Icon className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div className="metric-large">{value}</div>
        {trend && (
          <p className={cn(
            "text-xs flex items-center space-x-1",
            trend.isPositive ? "text-green-600" : "text-red-600"
          )}>
            <span>{trend.isPositive ? "↗" : "↘"}</span>
            <span>{trend.value}% {trend.label}</span>
          </p>
        )}
      </CardContent>
    </Card>
  );
}
"use client";

import React from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { DemographicRiskDistribution } from "@/services/report.service";
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Legend, Cell } from "recharts";
import { LucideIcon } from "lucide-react";

/**
 * Props interface for Demographic Risk Chart
 * Follows Interface Segregation Principle
 */
export interface DemographicRiskChartProps {
  title: string;
  description: string;
  icon: LucideIcon;
  data: DemographicRiskDistribution[];
  isLoading?: boolean;
}

/**
 * Risk level colors matching the system standard
 */
const RISK_COLORS = {
  nulo: "#22c55e", // Green
  bajo: "#84cc16", // Light green
  medio: "#eab308", // Yellow
  alto: "#f97316", // Orange
  muy_alto: "#ef4444", // Red
} as const;

/**
 * Risk level labels in Spanish
 */
const RISK_LABELS = {
  nulo: "Nulo",
  bajo: "Bajo",
  medio: "Medio",
  alto: "Alto",
  muy_alto: "Muy Alto",
} as const;

/**
 * Custom Tooltip for Risk Distribution
 */
interface RiskTooltipProps {
  active?: boolean;
  payload?: Array<{
    name: string;
    value: number;
    dataKey: string;
    payload: {
      category: string;
      risk_level: string;
      count: number;
    };
  }>;
}

const RiskTooltip: React.FC<RiskTooltipProps> = ({ active, payload }) => {
  if (!active || !payload || payload.length === 0) return null;

  const data = payload[0].payload;
  const riskLabel = RISK_LABELS[data.risk_level as keyof typeof RISK_LABELS] || data.risk_level;
  const riskColor = RISK_COLORS[data.risk_level as keyof typeof RISK_COLORS] || "#6b7280";

  return (
    <div className="bg-popover p-3 rounded-md shadow-lg border border-border">
      <p className="font-semibold text-sm mb-2">{data.category}</p>
      <div className="flex items-center space-x-2">
        <div
          className="w-3 h-3 rounded-full"
          style={{ backgroundColor: riskColor }}
        />
        <p className="text-xs">
          <span className="font-medium">{riskLabel}:</span>{" "}
          <span className="text-muted-foreground">{data.count} personas</span>
        </p>
      </div>
    </div>
  );
};

/**
 * Translate category values to readable Spanish labels
 */
function translateCategory(category: string): string {
  const translations: Record<string, string> = {
    "18-25": "18-25 años",
    "26-35": "26-35 años",
    "36-45": "36-45 años",
    "46-55": "46-55 años",
    "56+": "56+ años",
    diurno: "Diurno",
    nocturno: "Nocturno",
    mixto: "Mixto",
  };
  return translations[category] || category;
}

/**
 * Demographic Risk Chart Component
 * 
 * Displays risk distribution across demographic categories using stacked bar chart.
 * Follows SOLID principles:
 * - Single Responsibility: Visualizes demographic risk correlation only
 * - Open/Closed: Extensible via props
 * - Liskov Substitution: Can be swapped with other risk visualization components
 * - Interface Segregation: Focused props interface
 * - Dependency Inversion: Depends on props abstraction
 */
export function DemographicRiskChart({
  title,
  description,
  icon: Icon,
  data,
  isLoading = false,
}: DemographicRiskChartProps) {
  // Loading state
  if (isLoading) {
    return (
      <Card className="shadow-sm">
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Skeleton className="h-5 w-5 rounded" />
            <Skeleton className="h-5 w-32" />
          </CardTitle>
          <Skeleton className="h-4 w-48" />
        </CardHeader>
        <CardContent>
          <div className="h-64 flex items-center justify-center">
            <Skeleton className="h-48 w-full rounded" />
          </div>
        </CardContent>
      </Card>
    );
  }

  // Empty state
  if (!data || data.length === 0) {
    return (
      <Card className="shadow-sm">
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Icon className="h-5 w-5 text-muted-foreground" />
            <span>{title}</span>
          </CardTitle>
          <CardDescription>{description}</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-64 flex flex-col items-center justify-center text-center px-4">
            <Icon className="h-12 w-12 text-muted-foreground/50 mb-2" />
            <p className="text-muted-foreground text-sm font-medium">
              Sin datos de riesgo disponibles
            </p>
            <p className="text-muted-foreground/70 text-xs mt-1">
              Se requieren evaluaciones completadas para este análisis
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  // Transform data for stacked bar chart
  // Group by category and create stacked data
  interface CategoryRiskData {
    category: string;
    nulo: number;
    bajo: number;
    medio: number;
    alto: number;
    muy_alto: number;
  }

  const categoryMap = new Map<string, CategoryRiskData>();
  
  // Validate and process data
  const validData = data.filter(item => 
    item && 
    item.category && 
    item.risk_level && 
    item.count !== undefined && 
    item.count !== null &&
    item.count >= 0
  );
  
  validData.forEach((item) => {
    const category = translateCategory(item.category);
    if (!categoryMap.has(category)) {
      categoryMap.set(category, {
        category,
        nulo: 0,
        bajo: 0,
        medio: 0,
        alto: 0,
        muy_alto: 0,
      });
    }
    const categoryData = categoryMap.get(category)!;
    const riskLevel = item.risk_level;
    const count = item.count || 0;
    
    if (riskLevel === "nulo") {
      categoryData.nulo = count;
    } else if (riskLevel === "bajo") {
      categoryData.bajo = count;
    } else if (riskLevel === "medio") {
      categoryData.medio = count;
    } else if (riskLevel === "alto") {
      categoryData.alto = count;
    } else if (riskLevel === "muy_alto") {
      categoryData.muy_alto = count;
    }
  });

  const chartData = Array.from(categoryMap.values());
  
  // If no valid data after processing, show empty state
  if (chartData.length === 0) {
    return (
      <Card className="shadow-sm">
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Icon className="h-5 w-5 text-muted-foreground" />
            <span>{title}</span>
          </CardTitle>
          <CardDescription>{description}</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-64 flex flex-col items-center justify-center text-center px-4">
            <Icon className="h-12 w-12 text-muted-foreground/50 mb-2" />
            <p className="text-muted-foreground text-sm font-medium">
              Sin datos válidos disponibles
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  // Calculate totals for each category
  const categoryTotals = chartData.map((item) => ({
    ...item,
    total: item.nulo + item.bajo + item.medio + item.alto + item.muy_alto,
  }));

  const riskLevels: Array<keyof typeof RISK_COLORS> = ["nulo", "bajo", "medio", "alto", "muy_alto"];

  return (
    <Card className="shadow-sm hover:shadow-md transition-shadow">
      <CardHeader>
        <CardTitle className="flex items-center space-x-2">
          <Icon className="h-5 w-5 text-primary" />
          <span>{title}</span>
        </CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent className="pb-4">
        <ResponsiveContainer width="100%" height={320}>
          <BarChart
            data={categoryTotals}
            margin={{ top: 20, right: 30, left: 20, bottom: 80 }}
          >
            <XAxis
              dataKey="category"
              angle={-45}
              textAnchor="end"
              height={100}
              tick={{ fontSize: 10 }}
              interval={0}
            />
            <YAxis
              tick={{ fontSize: 11 }}
              width={50}
            />
            <Tooltip
              content={<RiskTooltip />}
              cursor={{ fill: 'rgba(0, 0, 0, 0.1)' }}
              formatter={(value: any, name?: string) => {
                // Translate risk level keys to Spanish labels
                const riskLevel = name as keyof typeof RISK_LABELS;
                const label = RISK_LABELS[riskLevel] || name || '';
                return [`${value} personas`, label];
              }}
              labelFormatter={(label) => `Categoría: ${label}`}
            />
            {riskLevels.map((level) => (
              <Bar
                key={level}
                dataKey={level}
                stackId="a"
                fill={RISK_COLORS[level]}
                name={RISK_LABELS[level]}
              />
            ))}
          </BarChart>
        </ResponsiveContainer>
        {/* Risk Level Color Glossary */}
        <div className="mt-3 pt-3 border-t border-border">
          <div className="mb-2">
            <p className="text-xs font-medium text-muted-foreground mb-2">Niveles de Riesgo:</p>
          </div>
          <div className="flex flex-wrap gap-4 text-xs mb-3">
            {riskLevels.map((level) => (
              <div key={level} className="flex items-center space-x-2">
                <div
                  className="w-4 h-4 rounded flex-shrink-0 border border-border/50"
                  style={{ backgroundColor: RISK_COLORS[level] }}
                />
                <span className="text-foreground font-medium">{RISK_LABELS[level]}</span>
              </div>
            ))}
          </div>
        </div>
        
        {/* Summary */}
        <div className="pt-3 border-t border-border">
          <div className="grid grid-cols-2 md:grid-cols-5 gap-3 text-xs">
            {riskLevels.map((level) => {
              const total = categoryTotals.reduce((sum, item) => sum + item[level], 0);
              return (
                <div key={level} className="flex items-center space-x-2">
                  <div
                    className="w-3 h-3 rounded-full flex-shrink-0"
                    style={{ backgroundColor: RISK_COLORS[level] }}
                  />
                  <span className="text-muted-foreground">{RISK_LABELS[level]}:</span>
                  <span className="font-semibold">{total.toLocaleString()}</span>
                </div>
              );
            })}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

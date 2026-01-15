"use client";

import React from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { DemographicDistribution } from "@/services/report.service";
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from "recharts";
import { LucideIcon } from "lucide-react";

/**
 * Props interface for Enhanced Demographic Chart
 * Follows Interface Segregation Principle - focused, minimal interface
 */
export interface EnhancedDemographicChartProps {
  title: string;
  description: string;
  icon: LucideIcon;
  data: DemographicDistribution[];
  isLoading?: boolean;
  chartType?: "donut"; // Only donut chart type supported
}

/**
 * Shared color palette for all demographic charts
 * All charts use the same set of distinct colors to ensure consistency
 * Each segment in each chart gets a different color from this shared palette
 */
const SHARED_COLOR_PALETTE = [
  "#3b82f6", // Blue
  "#a855f7", // Purple
  "#22c55e", // Green
  "#f97316", // Orange
  "#ef4444", // Red
  "#06b6d4", // Cyan
  "#f59e0b", // Amber
  "#8b5cf6", // Violet
  "#10b981", // Emerald
  "#ec4899", // Pink
  "#6366f1", // Indigo
  "#14b8a6", // Teal
] as const;

/**
 * Custom Tooltip Component
 * Single Responsibility: Handles tooltip rendering
 */
interface CustomTooltipProps {
  active?: boolean;
  payload?: Array<{
    name: string;
    value: number;
    payload: {
      category: string;
      count: number;
      percentage: number;
      displayName: string;
    };
  }>;
}

const CustomTooltip: React.FC<CustomTooltipProps> = ({ active, payload }) => {
  if (!active || !payload || payload.length === 0) return null;

  const data = payload[0].payload;
  return (
    <div className="bg-popover p-3 rounded-md shadow-lg border border-border">
      <p className="font-semibold text-sm mb-1">{data.displayName}</p>
      <p className="text-xs text-muted-foreground">
        Cantidad: <span className="font-medium text-foreground">{data.count}</span>
      </p>
      <p className="text-xs text-muted-foreground">
        Porcentaje: <span className="font-medium text-foreground">{data.percentage}%</span>
      </p>
    </div>
  );
};

/**
 * Translate category values to readable Spanish labels
 * Single Responsibility: Category translation logic
 */
function translateCategory(category: string): string {
  const translations: Record<string, string> = {
    // Age ranges
    "18-25": "18-25 años",
    "26-35": "26-35 años",
    "36-45": "36-45 años",
    "46-55": "46-55 años",
    "56+": "56+ años",
    // Marital status
    soltero: "Soltero/a",
    casado: "Casado/a",
    divorciado: "Divorciado/a",
    viudo: "Viudo/a",
    "union_libre": "Unión libre",
    // Shift types
    diurno: "Diurno",
    nocturno: "Nocturno",
    mixto: "Mixto",
    // Experience
    "<1 año": "< 1 año",
    "1-5 años": "1-5 años",
    "5-10 años": "5-10 años",
    "10+ años": "10+ años",
    "<1 anio": "< 1 año",
    "1-3 anios": "1-3 años",
    "3-5 anios": "3-5 años",
    "5+ anios": "5+ años",
  };
  return translations[category] || category;
}

/**
 * Enhanced Demographic Chart Component
 * 
 * Displays demographic distribution with multiple chart types (pie, donut, bar).
 * Follows SOLID principles:
 * - Single Responsibility: Renders demographic charts only
 * - Open/Closed: Extensible via props (chartType, colorScheme)
 * - Liskov Substitution: Can be swapped with other chart components
 * - Interface Segregation: Focused props interface
 * - Dependency Inversion: Depends on props abstraction, not concrete implementations
 */
export function EnhancedDemographicChart({
  title,
  description,
  icon: Icon,
  data,
  isLoading = false,
  chartType = "donut",
}: EnhancedDemographicChartProps) {
  // Use shared color palette for all charts to ensure consistency
  const colors = SHARED_COLOR_PALETTE;

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
              Sin datos registrados
            </p>
            <p className="text-muted-foreground/70 text-xs mt-1">
              Importa personal con este campo para ver el análisis
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  // Calculate total and format data
  const total = data.reduce((sum, item) => sum + (item.count || 0), 0);
  
  // Validate and filter out invalid data
  const validData = data.filter(item => item && item.count !== undefined && item.count !== null);
  
  if (validData.length === 0) {
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
  
  const chartData = validData.map((item, index) => ({
    ...item,
    name: translateCategory(item.category),
    displayName: translateCategory(item.category),
    value: item.count || 0,
    percentage: total > 0 ? Math.round(((item.count || 0) / total) * 100) : 0,
    colorIndex: index,
  }));

  // Render donut chart (only supported chart type)
  const renderChart = () => {
    return (
      <ResponsiveContainer width="100%" height={320}>
        <PieChart margin={{ top: 10, right: 10, bottom: 10, left: 10 }}>
          <Pie
            data={chartData}
            cx="50%"
            cy="45%"
            innerRadius={70}
            outerRadius={110}
            paddingAngle={3}
            dataKey="value"
            label={(entry: any) => {
              const total = chartData.reduce((sum, item) => sum + item.value, 0);
              const percentage = total > 0 ? Math.round((entry.value / total) * 100) : 0;
              // Show all labels regardless of percentage
              return `${percentage}%`;
            }}
            labelLine={false}
          >
            {chartData.map((entry, index) => (
              <Cell key={`cell-${index}`} fill={colors[index % colors.length]} />
            ))}
            </Pie>
          <Tooltip content={<CustomTooltip />} cursor={{ fill: 'rgba(0, 0, 0, 0.1)' }} />
        </PieChart>
      </ResponsiveContainer>
    );
  };

  return (
    <Card className="shadow-sm hover:shadow-md transition-shadow h-full flex flex-col">
      <CardHeader className="pb-3">
        <CardTitle className="flex items-center space-x-2 text-base">
          <Icon className="h-5 w-5 text-primary" />
          <span>{title}</span>
        </CardTitle>
        <CardDescription className="text-xs">{description}</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 flex flex-col pb-0 pt-0">
        <div className="flex-1 min-h-0">
          {renderChart()}
        </div>
        {/* Category Legend & Analysis */}
        <div className="mt-1 pt-1 border-t border-border space-y-3">
          {/* Category Legend */}
          <div className="min-h-[60px]">
            <div className="text-xs text-muted-foreground mb-2 font-medium">Categorías</div>
            <div className="grid grid-cols-2 gap-2">
              {chartData.slice(0, 6).map((entry, index) => (
                <div key={index} className="flex items-center space-x-2">
                  <div
                    className="w-3 h-3 rounded-full flex-shrink-0 border border-border/50"
                    style={{ backgroundColor: colors[index % colors.length] }}
                  />
                  <span className="text-xs text-foreground truncate">{entry.displayName}</span>
                </div>
              ))}
              {chartData.length > 6 && (
                <div className="col-span-2 text-xs text-muted-foreground text-center">
                  +{chartData.length - 6} más
                </div>
              )}
            </div>
          </div>

          {/* Analysis Insights */}
          <div className="min-h-[50px]">
            <div className="text-xs text-muted-foreground mb-2 font-medium">Análisis</div>
            <div className="grid grid-cols-2 gap-3 text-xs">
              <div>
                <span className="text-muted-foreground">Más frecuente:</span>
                <div className="font-semibold text-foreground">
                  {chartData.sort((a, b) => b.percentage - a.percentage)[0]?.displayName || 'N/A'}
                </div>
              </div>
              <div>
                <span className="text-muted-foreground">Uniformidad:</span>
                <div className="font-semibold text-foreground">
                  {(() => {
                    const avg = total / chartData.length;
                    const variance = chartData.reduce((sum, d) => sum + Math.pow(d.value - avg, 2), 0) / chartData.length;
                    const stdDev = Math.sqrt(variance);
                    const uniformity = Math.max(0, 100 - (stdDev / avg) * 100);
                    return Math.round(uniformity) + '%';
                  })()}
                </div>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

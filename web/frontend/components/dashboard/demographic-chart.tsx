"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { DemographicDistribution } from "@/services/report.service";
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Cell } from "recharts";
import { LucideIcon } from "lucide-react";

interface DemographicChartProps {
  title: string;
  description: string;
  icon: LucideIcon;
  data: DemographicDistribution[];
  isLoading?: boolean;
  colorScheme?: "blue" | "green" | "purple" | "orange";
}

const colorSchemes = {
  blue: ["#3b82f6", "#60a5fa", "#93c5fd", "#bfdbfe", "#dbeafe"],
  green: ["#22c55e", "#4ade80", "#86efac", "#bbf7d0", "#dcfce7"],
  purple: ["#a855f7", "#c084fc", "#d8b4fe", "#e9d5ff", "#f3e8ff"],
  orange: ["#f97316", "#fb923c", "#fdba74", "#fed7aa", "#ffedd5"],
};

/**
 * Demographic Chart Component
 * 
 * Displays demographic distribution data in a horizontal bar chart format.
 * Used to visualize age, marital status, shift type, and experience distributions.
 */
export function DemographicChart({
  title,
  description,
  icon: Icon,
  data,
  isLoading = false,
  colorScheme = "blue",
}: DemographicChartProps) {
  const colors = colorSchemes[colorScheme];

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
          <div className="h-48 flex items-center justify-center">
            <Skeleton className="h-40 w-full" />
          </div>
        </CardContent>
      </Card>
    );
  }

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
          <div className="h-48 flex flex-col items-center justify-center text-center px-4">
            <p className="text-muted-foreground text-sm">
              Sin datos registrados
            </p>
            <p className="text-muted-foreground/70 text-xs mt-1">
              Importa personal con este campo para ver el analisis
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  // Calculate total for percentage display
  const total = data.reduce((sum, item) => sum + item.count, 0);

  // Format data for display (add percentage)
  const chartData = data.map((item) => ({
    ...item,
    percentage: total > 0 ? Math.round((item.count / total) * 100) : 0,
    // Translate common Spanish values for display
    displayName: translateCategory(item.category),
  }));

  return (
    <Card className="shadow-sm hover:shadow-md transition-shadow">
      <CardHeader>
        <CardTitle className="flex items-center space-x-2">
          <Icon className="h-5 w-5 text-primary" />
          <span>{title}</span>
        </CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={200}>
          <BarChart
            data={chartData}
            layout="vertical"
            margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
          >
            <XAxis type="number" hide />
            <YAxis
              type="category"
              dataKey="displayName"
              tick={{ fontSize: 12 }}
              width={100}
            />
            <Tooltip
              formatter={(value) => [`${value}`, "Total"]}
              labelStyle={{ fontWeight: "bold" }}
              contentStyle={{
                backgroundColor: "white",
                border: "1px solid #e5e7eb",
                borderRadius: "8px",
                boxShadow: "0 2px 4px rgba(0,0,0,0.1)",
              }}
            />
            <Bar dataKey="count" radius={[0, 4, 4, 0]}>
              {chartData.map((_, index) => (
                <Cell
                  key={`cell-${index}`}
                  fill={colors[index % colors.length]}
                />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
        {/* Legend with counts */}
        <div className="mt-4 grid grid-cols-2 gap-2 text-sm">
          {chartData.slice(0, 4).map((item, index) => (
            <div
              key={item.category}
              className="flex items-center space-x-2"
            >
              <div
                className="w-3 h-3 rounded-full"
                style={{ backgroundColor: colors[index % colors.length] }}
              />
              <span className="text-muted-foreground truncate">
                {item.displayName}: {item.count}
              </span>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

/**
 * Translate category values to more readable Spanish labels
 */
function translateCategory(category: string): string {
  const translations: Record<string, string> = {
    // Age ranges
    "18-25": "18-25 anios",
    "26-35": "26-35 anios",
    "36-45": "36-45 anios",
    "46-55": "46-55 anios",
    "56+": "56+ anios",
    // Marital status
    soltero: "Soltero/a",
    casado: "Casado/a",
    divorciado: "Divorciado/a",
    viudo: "Viudo/a",
    "union_libre": "Union libre",
    // Shift types
    diurno: "Diurno",
    nocturno: "Nocturno",
    mixto: "Mixto",
    // Experience
    "<1 año": "< 1 anio",
    "1-5 años": "1-5 anios",
    "5-10 años": "5-10 anios",
    "10+ años": "10+ anios",
    "<1 anio": "< 1 anio",
    "1-3 anios": "1-3 anios",
    "3-5 anios": "3-5 anios",
    "5+ anios": "5+ anios",
  };
  return translations[category] || category;
}

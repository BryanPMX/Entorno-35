"use client";

import React from "react";
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Cell } from "recharts";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  NOM35_RISK_HEX_COLORS,
  getGuideIIIRiskLevelByScore,
  getNom35RiskLabel,
} from "@/lib/nom35-risk";
import { translations } from "@/lib/translations";

interface DepartmentData {
  department: string;
  risk_level: string;
  count: number;
  avg_score?: number;
}

interface DepartmentHeatmapProps {
  data: DepartmentData[];
  isLoading?: boolean;
}

interface HeatmapTooltipProps {
  active?: boolean;
  payload?: Array<{ payload: { averageRisk: number; totalAssessments: number } }>;
  label?: string;
}

function HeatmapTooltip({ active, payload, label }: HeatmapTooltipProps) {
  if (active && payload && payload.length) {
    const datum = payload[0]?.payload;
    if (!datum) return null;
    return (
      <div className="bg-popover p-3 rounded-md shadow-md border">
        <p className="font-medium">{label}</p>
        <p className="text-sm text-muted-foreground">Riesgo Promedio: {datum.averageRisk}</p>
        <p className="text-sm text-muted-foreground">
          {translations.charts.totalAssessments}: {datum.totalAssessments}
        </p>
      </div>
    );
  }
  return null;
}

export function DepartmentHeatmap({ data, isLoading = false }: DepartmentHeatmapProps) {
  // Group data by department and use actual average scores from backend
  const departmentStats = data.reduce((acc, item) => {
    if (!acc[item.department]) {
      acc[item.department] = {
        department: item.department,
        totalAssessments: 0,
        avgScore: item.avg_score ?? null, // Use actual average score from backend
      };
    }

    acc[item.department].totalAssessments += item.count;
    // avg_score should be the same for all entries of the same department
    if (item.avg_score !== undefined && acc[item.department].avgScore === null) {
      acc[item.department].avgScore = item.avg_score;
    }

    return acc;
  }, {} as Record<string, { department: string; totalAssessments: number; avgScore: number | null }>);

  const chartData = Object.values(departmentStats).map((dept) => {
    const avgScore = dept.avgScore ?? 0;
    const riskLevel = getGuideIIIRiskLevelByScore(dept.avgScore);
    return {
      department: dept.department,
      averageRisk: Math.round(avgScore),
      totalAssessments: dept.totalAssessments,
      color: NOM35_RISK_HEX_COLORS[riskLevel],
    };
  });

  if (isLoading) {
    return (
      <Card className="animate-in fade-in slide-in-from-bottom-4">
        <CardHeader>
          <CardTitle>{translations.charts.departmentHeatmap}</CardTitle>
          <CardDescription>{translations.charts.departmentHeatmapDesc}</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-[400px] flex items-center justify-center">
            <div className="animate-pulse bg-muted rounded w-full h-32"></div>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!chartData.length) {
    return (
      <Card className="animate-in fade-in slide-in-from-bottom-4">
        <CardHeader>
          <CardTitle>{translations.charts.departmentHeatmap}</CardTitle>
          <CardDescription>{translations.charts.departmentHeatmapDesc}</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-[400px] flex items-center justify-center">
            <p className="text-muted-foreground">{translations.charts.noDataAvailable}</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="animate-in fade-in slide-in-from-bottom-4">
      <CardHeader>
        <CardTitle>{translations.charts.departmentHeatmap}</CardTitle>
        <CardDescription>{translations.charts.departmentHeatmapDesc}</CardDescription>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={350}>
          <BarChart data={chartData} margin={{ left: 80, bottom: 20, right: 20, top: 30 }}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis
              dataKey="department"
              angle={-45}
              textAnchor="end"
              height={80}
              fontSize={12}
            />
            <YAxis
              label={{
                value: translations.charts.riskScore,
                angle: -90,
                position: 'left',
                offset: 0,
                style: { textAnchor: 'middle' },
                dy: '50%'
              }}
              domain={[0, 100]}
            />
            <Tooltip content={<HeatmapTooltip />} />
            <Bar
              dataKey="averageRisk"
              radius={[4, 4, 0, 0]}
            >
              {chartData.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={entry.color} />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>

        {/* Risk Level Legend */}
        <div className="mt-2 pt-4 border-t border-border">
          <div className="mb-3">
            <p className="text-sm font-medium text-muted-foreground">Escala de Riesgo</p>
          </div>
          <div className="grid grid-cols-2 md:grid-cols-5 gap-3">
            <div className="flex items-center space-x-2">
                <div
                  className="w-4 h-4 rounded flex-shrink-0 border border-border/50"
                  style={{ backgroundColor: NOM35_RISK_HEX_COLORS.nulo }}
                />
              <div className="flex flex-col">
                <span className="text-sm font-medium text-foreground">{getNom35RiskLabel("nulo")}</span>
                <span className="text-xs text-muted-foreground">0-49</span>
              </div>
            </div>
            <div className="flex items-center space-x-2">
                <div
                  className="w-4 h-4 rounded flex-shrink-0 border border-border/50"
                  style={{ backgroundColor: NOM35_RISK_HEX_COLORS.bajo }}
                />
              <div className="flex flex-col">
                <span className="text-sm font-medium text-foreground">{getNom35RiskLabel("bajo")}</span>
                <span className="text-xs text-muted-foreground">50-74</span>
              </div>
            </div>
            <div className="flex items-center space-x-2">
                <div
                  className="w-4 h-4 rounded flex-shrink-0 border border-border/50"
                  style={{ backgroundColor: NOM35_RISK_HEX_COLORS.medio }}
                />
              <div className="flex flex-col">
                <span className="text-sm font-medium text-foreground">{getNom35RiskLabel("medio")}</span>
                <span className="text-xs text-muted-foreground">75-98</span>
              </div>
            </div>
            <div className="flex items-center space-x-2">
                <div
                  className="w-4 h-4 rounded flex-shrink-0 border border-border/50"
                  style={{ backgroundColor: NOM35_RISK_HEX_COLORS.alto }}
                />
              <div className="flex flex-col">
                <span className="text-sm font-medium text-foreground">{getNom35RiskLabel("alto")}</span>
                <span className="text-xs text-muted-foreground">99-139</span>
              </div>
            </div>
            <div className="flex items-center space-x-2">
                <div
                  className="w-4 h-4 rounded flex-shrink-0 border border-border/50"
                  style={{ backgroundColor: NOM35_RISK_HEX_COLORS.muy_alto }}
                />
              <div className="flex flex-col">
                <span className="text-sm font-medium text-foreground">{getNom35RiskLabel("muy_alto")}</span>
                <span className="text-xs text-muted-foreground">140+</span>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

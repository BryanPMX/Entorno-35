"use client";

import React from "react";
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from "recharts";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { translations } from "@/lib/translations";

interface RiskData {
  risk_level: string;
  count: number;
}

interface RiskDistributionChartProps {
  data: RiskData[];
  isLoading?: boolean;
}

const RISK_COLORS = {
  nulo: "#22c55e", // Green
  bajo: "#84cc16", // Light green
  medio: "#eab308", // Yellow
  alto: "#f97316", // Orange
  muy_alto: "#ef4444", // Red
};

const RISK_LABELS = {
  nulo: "Nulo",
  bajo: "Bajo",
  medio: "Medio",
  alto: "Alto",
  muy_alto: "Muy Alto",
};

export function RiskDistributionChart({ data, isLoading = false }: RiskDistributionChartProps) {
  // Validate and filter data
  const validData = (data || []).filter(item => 
    item && 
    item.risk_level && 
    item.count !== undefined && 
    item.count !== null &&
    item.count >= 0
  );
  
  const chartData = validData.map((item) => ({
    name: RISK_LABELS[item.risk_level as keyof typeof RISK_LABELS] || item.risk_level,
    value: item.count || 0,
    fill: RISK_COLORS[item.risk_level as keyof typeof RISK_COLORS] || "#6b7280",
  }));

  const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div className="bg-popover p-3 rounded-md shadow-md border">
          <p className="font-medium">{data.name}</p>
          <p className="text-sm text-muted-foreground">{data.value} {translations.charts.totalAssessments.toLowerCase()}</p>
        </div>
      );
    }
    return null;
  };

  if (isLoading) {
    return (
      <Card className="animate-in fade-in slide-in-from-bottom-4">
        <CardHeader>
          <CardTitle>{translations.charts.riskDistribution}</CardTitle>
          <CardDescription>{translations.charts.riskDistributionDesc}</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-[240px] flex items-center justify-center">
            <div className="animate-pulse bg-muted rounded-full w-32 h-32"></div>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!chartData.length) {
    return (
      <Card className="animate-in fade-in slide-in-from-bottom-4">
        <CardHeader>
          <CardTitle>{translations.charts.riskDistribution}</CardTitle>
          <CardDescription>{translations.charts.riskDistributionDesc}</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-[240px] flex items-center justify-center">
            <p className="text-muted-foreground">{translations.charts.noDataAvailable}</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="animate-in fade-in slide-in-from-bottom-4">
      <CardHeader>
        <CardTitle>{translations.charts.riskDistribution}</CardTitle>
        <CardDescription>{translations.charts.riskDistributionDesc}</CardDescription>
      </CardHeader>
      <CardContent className="pb-4">
        <ResponsiveContainer width="100%" height={280}>
          <PieChart margin={{ top: 30, right: 20, bottom: 60, left: 20 }}>
            <Pie
              data={chartData}
              cx="50%"
              cy="45%"
              innerRadius={50}
              outerRadius={75}
              paddingAngle={2}
              dataKey="value"
              label={(entry: any) => {
                const total = chartData.reduce((sum, item) => sum + item.value, 0);
                const percentage = total > 0 ? Math.round((entry.value / total) * 100) : 0;
                // Only show label if percentage is significant (> 5%)
                if (percentage < 5) return "";
                return `${percentage}%`;
              }}
              labelLine={false}
            >
              {chartData.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={entry.fill} />
              ))}
            </Pie>
            <Tooltip content={<CustomTooltip />} />
            <Legend
              verticalAlign="bottom"
              height={50}
              iconType="circle"
              wrapperStyle={{ paddingTop: "10px" }}
              formatter={(value, entry) => (
                <span style={{ color: entry.color, fontSize: "11px" }}>{value}</span>
              )}
            />
          </PieChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
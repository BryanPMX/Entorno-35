"use client";

import React from "react";
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from "recharts";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

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
  const chartData = data.map((item) => ({
    name: RISK_LABELS[item.risk_level as keyof typeof RISK_LABELS] || item.risk_level,
    value: item.count,
    fill: RISK_COLORS[item.risk_level as keyof typeof RISK_COLORS] || "#6b7280",
  }));

  const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div className="bg-popover p-3 rounded-md shadow-md border">
          <p className="font-medium">{data.name}</p>
          <p className="text-sm text-muted-foreground">{data.value} assessments</p>
        </div>
      );
    }
    return null;
  };

  if (isLoading) {
    return (
      <Card className="animate-in fade-in slide-in-from-bottom-4">
        <CardHeader>
          <CardTitle>Risk Distribution</CardTitle>
          <CardDescription>Distribution of risk levels across assessments</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-[300px] flex items-center justify-center">
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
          <CardTitle>Risk Distribution</CardTitle>
          <CardDescription>Distribution of risk levels across assessments</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-[300px] flex items-center justify-center">
            <p className="text-muted-foreground">No assessment data available</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="animate-in fade-in slide-in-from-bottom-4">
      <CardHeader>
        <CardTitle>Risk Distribution</CardTitle>
        <CardDescription>Distribution of risk levels across assessments</CardDescription>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={300}>
          <PieChart>
            <Pie
              data={chartData}
              cx="50%"
              cy="50%"
              innerRadius={60}
              outerRadius={100}
              paddingAngle={2}
              dataKey="value"
            >
              {chartData.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={entry.fill} />
              ))}
            </Pie>
            <Tooltip content={<CustomTooltip />} />
            <Legend
              verticalAlign="bottom"
              height={36}
              formatter={(value, entry) => (
                <span style={{ color: entry.color }}>{value}</span>
              )}
            />
          </PieChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
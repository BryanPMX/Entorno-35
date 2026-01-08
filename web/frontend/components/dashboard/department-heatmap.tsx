"use client";

import React from "react";
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from "recharts";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

interface DepartmentData {
  department: string;
  risk_level: string;
  count: number;
}

interface DepartmentHeatmapProps {
  data: DepartmentData[];
  isLoading?: boolean;
}

const RISK_COLORS = {
  nulo: "#22c55e", // Green
  bajo: "#84cc16", // Light green
  medio: "#eab308", // Yellow
  alto: "#f97316", // Orange
  muy_alto: "#ef4444", // Red
};

export function DepartmentHeatmap({ data, isLoading = false }: DepartmentHeatmapProps) {
  // Group data by department and calculate average risk score
  const departmentStats = data.reduce((acc, item) => {
    if (!acc[item.department]) {
      acc[item.department] = {
        department: item.department,
        totalAssessments: 0,
        riskScore: 0,
      };
    }

    acc[item.department].totalAssessments += item.count;

    // Convert risk level to numeric score for averaging
    const riskScores = {
      nulo: 10,
      bajo: 25,
      medio: 35,
      alto: 45,
      muy_alto: 55,
    };

    acc[item.department].riskScore +=
      (riskScores[item.risk_level as keyof typeof riskScores] || 0) * item.count;

    return acc;
  }, {} as Record<string, { department: string; totalAssessments: number; riskScore: number }>);

  const chartData = Object.values(departmentStats).map((dept) => ({
    department: dept.department,
    averageRisk: Math.round(dept.riskScore / dept.totalAssessments),
    totalAssessments: dept.totalAssessments,
    color: dept.riskScore / dept.totalAssessments < 20 ? RISK_COLORS.nulo :
           dept.riskScore / dept.totalAssessments < 30 ? RISK_COLORS.bajo :
           dept.riskScore / dept.totalAssessments < 40 ? RISK_COLORS.medio :
           dept.riskScore / dept.totalAssessments < 50 ? RISK_COLORS.alto :
           RISK_COLORS.muy_alto,
  }));

  const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div className="bg-popover p-3 rounded-md shadow-md border">
          <p className="font-medium">{label}</p>
          <p className="text-sm text-muted-foreground">
            Average Risk: {data.averageRisk}/100
          </p>
          <p className="text-sm text-muted-foreground">
            Total Assessments: {data.totalAssessments}
          </p>
        </div>
      );
    }
    return null;
  };

  if (isLoading) {
    return (
      <Card className="animate-in fade-in slide-in-from-bottom-4">
        <CardHeader>
          <CardTitle>Department Risk Heatmap</CardTitle>
          <CardDescription>Average risk levels by department</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-[300px] flex items-center justify-center">
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
          <CardTitle>Department Risk Heatmap</CardTitle>
          <CardDescription>Average risk levels by department</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-[300px] flex items-center justify-center">
            <p className="text-muted-foreground">No department data available</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="animate-in fade-in slide-in-from-bottom-4">
      <CardHeader>
        <CardTitle>Department Risk Heatmap</CardTitle>
        <CardDescription>Average risk levels by department</CardDescription>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis
              dataKey="department"
              angle={-45}
              textAnchor="end"
              height={80}
              fontSize={12}
            />
            <YAxis
              label={{ value: 'Risk Score', angle: -90, position: 'insideLeft' }}
              domain={[0, 60]}
            />
            <Tooltip content={<CustomTooltip />} />
            <Bar
              dataKey="averageRisk"
              fill="#8884d8"
              radius={[4, 4, 0, 0]}
            >
              {chartData.map((entry, index) => (
                <Bar key={`cell-${index}`} fill={entry.color} />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
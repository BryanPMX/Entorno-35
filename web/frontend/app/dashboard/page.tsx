"use client";

import { useQuery } from "@tanstack/react-query";
import { Users, FileText, TrendingUp, Plus } from "lucide-react";
import { MetricCard } from "@/components/dashboard/metric-card";
import { RiskDistributionChart } from "@/components/dashboard/risk-distribution-chart";
import { DepartmentHeatmap } from "@/components/dashboard/department-heatmap";
import { AssessmentWizard } from "@/components/assessments/assessment-wizard";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { reportService } from "@/services/report.service";

/**
 * Dashboard Page
 *
 * Comprehensive admin dashboard with metrics, charts, and assessment creation.
 * Features professional interactivity, smooth animations, and data visualization.
 */
export default function DashboardPage() {
  // Fetch general report data for dashboard metrics and charts
  const { data: reportData, isLoading: reportLoading } = useQuery({
    queryKey: ["general-report"],
    queryFn: () => reportService.getGeneralReport(),
  });

  const metrics = [
    {
      title: "Total Staff",
      value: reportData?.total_staff || 0,
      icon: Users,
      trend: reportData ? { value: 5, label: "from last period", isPositive: true } : undefined,
    },
    {
      title: "Completed Assessments",
      value: reportData?.completed_assessments || 0,
      icon: FileText,
      trend: reportData ? {
        value: Math.round(((reportData.completed_assessments || 0) / (reportData.total_staff || 1)) * 100),
        label: "participation rate",
        isPositive: true
      } : undefined,
    },
    {
      title: "Risk Distribution",
      value: reportData?.risk_distribution.length || 0,
      icon: TrendingUp,
      trend: reportData ? {
        value: (reportData.risk_distribution.find(r => r.risk_level === 'alto')?.count || 0) +
               (reportData.risk_distribution.find(r => r.risk_level === 'muy_alto')?.count || 0),
        label: "high risk cases",
        isPositive: false
      } : undefined,
    },
  ];

  const hasData = reportData && (reportData.total_staff > 0 || reportData.completed_assessments > 0);

  if (!hasData && !reportLoading) {
    return (
      <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4">
        <div>
          <h1 className="heading-1">Dashboard</h1>
          <p className="label-muted mt-2">
            Welcome to Entorno 35 - NOM-035 Compliance Platform
          </p>
        </div>

        <Card className="animate-in fade-in slide-in-from-bottom-4">
          <CardHeader className="text-center">
            <div className="mx-auto w-16 h-16 bg-muted rounded-full flex items-center justify-center mb-4">
              <FileText className="h-8 w-8 text-muted-foreground" />
            </div>
            <CardTitle className="heading-2">Get Started with NOM-035 Assessments</CardTitle>
            <CardDescription className="text-lg">
              Create your first assessment cycle and begin tracking workplace psychosocial risk factors.
            </CardDescription>
          </CardHeader>
          <CardContent className="flex justify-center">
            <AssessmentWizard
              trigger={
                <Button size="lg" className="flex items-center space-x-2">
                  <Plus className="h-5 w-5" />
                  <span>Create First Assessment</span>
                </Button>
              }
            />
          </CardContent>
        </Card>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <Card className="hover-lift animate-in fade-in slide-in-from-bottom-4">
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Users className="h-5 w-5" />
                <span>Staff Management</span>
              </CardTitle>
              <CardDescription>
                Import and manage your staff members
              </CardDescription>
            </CardHeader>
            <CardContent>
              <p className="label-muted">
                Start by importing your staff list via CSV upload for easy bulk management.
              </p>
            </CardContent>
          </Card>

          <Card className="hover-lift animate-in fade-in slide-in-from-bottom-4">
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <FileText className="h-5 w-5" />
                <span>Assessment Creation</span>
              </CardTitle>
              <CardDescription>
                Generate secure assessment links
              </CardDescription>
            </CardHeader>
            <CardContent>
              <p className="label-muted">
                Create assessment cycles and automatically generate secure links for staff participation.
              </p>
            </CardContent>
          </Card>

          <Card className="hover-lift animate-in fade-in slide-in-from-bottom-4">
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <TrendingUp className="h-5 w-5" />
                <span>Compliance Reports</span>
              </CardTitle>
              <CardDescription>
                Generate NOM-035 compliance reports
              </CardDescription>
            </CardHeader>
            <CardContent>
              <p className="label-muted">
                Access individual assessments and company-wide compliance reports with actionable recommendations.
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="heading-1">Dashboard</h1>
          <p className="label-muted mt-2">
            NOM-035 Compliance Overview - {reportData?.period || new Date().getFullYear()}
          </p>
        </div>

        <AssessmentWizard
          trigger={
            <Button className="flex items-center space-x-2">
              <Plus className="h-4 w-4" />
              <span>Create Assessment</span>
            </Button>
          }
        />
      </div>

      {/* Metrics Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {metrics.map((metric, index) => (
          <MetricCard
            key={metric.title}
            title={metric.title}
            value={metric.value}
            icon={metric.icon}
            trend={metric.trend}
            isLoading={reportLoading}
            className="animate-in fade-in slide-in-from-bottom-4"
            style={{ animationDelay: `${index * 100}ms` }}
          />
        ))}
      </div>

      {/* Charts */}
      <div className="grid gap-6 md:grid-cols-2">
        <RiskDistributionChart
          data={reportData?.risk_distribution || []}
          isLoading={reportLoading}
        />
        <DepartmentHeatmap
          data={reportData?.department_heatmap || []}
          isLoading={reportLoading}
        />
      </div>

      {/* Quick Actions */}
      <Card className="animate-in fade-in slide-in-from-bottom-4">
        <CardHeader>
          <CardTitle>Quick Actions</CardTitle>
          <CardDescription>
            Common tasks to manage your NOM-035 compliance program
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-4">
            <Button variant="outline" className="flex items-center space-x-2">
              <Users className="h-4 w-4" />
              <span>Manage Staff</span>
            </Button>
            <Button variant="outline" className="flex items-center space-x-2">
              <FileText className="h-4 w-4" />
              <span>View Reports</span>
            </Button>
            <AssessmentWizard
              trigger={
                <Button variant="outline" className="flex items-center space-x-2">
                  <Plus className="h-4 w-4" />
                  <span>New Assessment Cycle</span>
                </Button>
              }
            />
          </div>
        </CardContent>
      </Card>
    </div>
  );
}


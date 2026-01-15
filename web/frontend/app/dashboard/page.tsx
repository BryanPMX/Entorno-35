"use client";

import { useQuery } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { Users, FileText, TrendingUp, Plus, Calendar, Heart, Clock, Briefcase } from "lucide-react";
import { MetricCard } from "@/components/dashboard/metric-card";
import { RiskDistributionChart } from "@/components/dashboard/risk-distribution-chart";
import { DepartmentHeatmap } from "@/components/dashboard/department-heatmap";
import { EnhancedDemographicChart } from "@/components/dashboard/enhanced-demographic-chart";
import { DemographicRiskChart } from "@/components/dashboard/demographic-risk-chart";
import { AssessmentWizard } from "@/components/assessments/assessment-wizard";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { reportService } from "@/services/report.service";
import { translations } from "@/lib/translations";

/**
 * Dashboard Page
 *
 * Comprehensive admin dashboard with metrics, charts, and assessment creation.
 * Features professional interactivity, smooth animations, and data visualization.
 */
export default function DashboardPage() {
  const router = useRouter();

  // Fetch general report data for dashboard metrics and charts
  const { data: reportData, isLoading: reportLoading } = useQuery({
    queryKey: ["general-report"],
    queryFn: () => reportService.getGeneralReport(),
  });

  const metrics = [
    {
      title: translations.dashboard.totalStaff,
      value: reportData?.total_staff || 0,
      icon: Users,
      trend: reportData ? { value: 5, label: translations.dashboard.fromLastPeriod, isPositive: true } : undefined,
    },
    {
      title: translations.dashboard.completedAssessments,
      value: reportData?.completed_assessments || 0,
      icon: FileText,
      trend: reportData ? {
        value: Math.round(((reportData.completed_assessments || 0) / (reportData.total_staff || 1)) * 100),
        label: translations.dashboard.participationRate,
        isPositive: true
      } : undefined,
    },
    {
      title: translations.dashboard.riskDistribution,
      value: reportData?.risk_distribution?.length || 0,
      icon: TrendingUp,
      trend: reportData ? {
        value: (reportData.risk_distribution?.find(r => r.risk_level === 'alto')?.count || 0) +
               (reportData.risk_distribution?.find(r => r.risk_level === 'muy_alto')?.count || 0),
        label: translations.dashboard.highRiskCases,
        isPositive: false
      } : undefined,
    },
  ];

  const hasData = reportData && (reportData.total_staff > 0 || reportData.completed_assessments > 0);

  if (!hasData && !reportLoading) {
    return (
      <div className="space-y-12 animate-in fade-in slide-in-from-bottom-4">
        <div>
          <h1 className="heading-1">{translations.dashboard.title}</h1>
          <p className="label-muted mt-2">
            {translations.dashboard.welcome}
          </p>
        </div>

        <Card className="animate-in fade-in slide-in-from-bottom-4">
          <CardHeader className="text-center">
            <div className="mx-auto w-16 h-16 bg-muted rounded-full flex items-center justify-center mb-4">
              <FileText className="h-8 w-8 text-muted-foreground" />
            </div>
            <CardTitle className="heading-2">{translations.dashboard.getStarted}</CardTitle>
            <CardDescription className="text-lg">
              {translations.dashboard.getStartedDesc}
            </CardDescription>
          </CardHeader>
          <CardContent className="flex justify-center">
            <AssessmentWizard
              trigger={
                <Button size="lg" className="flex items-center space-x-2">
                  <Plus className="h-5 w-5" />
                  <span>{translations.dashboard.createFirstAssessment}</span>
                </Button>
              }
            />
          </CardContent>
        </Card>

        <div className="grid gap-8 md:grid-cols-2 lg:grid-cols-3 animate-in fade-in slide-in-from-bottom-4" style={{ animationDelay: '300ms' }}>
          <Card className="hover-lift shadow-sm animate-in fade-in slide-in-from-bottom-4">
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Users className="h-5 w-5" />
                <span>{translations.onboarding.staffManagement}</span>
              </CardTitle>
              <CardDescription>
                {translations.onboarding.staffManagementDesc}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <p className="label-muted">
                {translations.onboarding.staffManagementHint}
              </p>
            </CardContent>
          </Card>

          <Card className="hover-lift shadow-sm animate-in fade-in slide-in-from-bottom-4" style={{ animationDelay: '500ms' }}>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <FileText className="h-5 w-5" />
                <span>{translations.onboarding.assessmentCreation}</span>
              </CardTitle>
              <CardDescription>
                {translations.onboarding.assessmentCreationDesc}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <p className="label-muted">
                {translations.onboarding.assessmentCreationHint}
              </p>
            </CardContent>
          </Card>

          <Card className="hover-lift shadow-sm animate-in fade-in slide-in-from-bottom-4" style={{ animationDelay: '600ms' }}>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <TrendingUp className="h-5 w-5" />
                <span>{translations.onboarding.complianceReports}</span>
              </CardTitle>
              <CardDescription>
                {translations.onboarding.complianceReportsDesc}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <p className="label-muted">
                {translations.onboarding.complianceReportsHint}
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-12 animate-in fade-in slide-in-from-bottom-4">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="heading-1">{translations.dashboard.title}</h1>
          <p className="label-muted mt-2">
            {translations.dashboard.subtitle} - {reportData?.period || new Date().getFullYear()}
          </p>
        </div>

        <AssessmentWizard
          trigger={
            <Button className="flex items-center space-x-2">
              <Plus className="h-4 w-4" />
              <span>{translations.dashboard.createAssessment}</span>
            </Button>
          }
        />
      </div>

      {/* Metrics Cards */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4 animate-in fade-in slide-in-from-bottom-4">
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

      {/* Risk Charts */}
      <div className="grid gap-8 md:grid-cols-2 animate-in fade-in slide-in-from-bottom-4" style={{ animationDelay: '200ms' }}>
        <RiskDistributionChart
          data={reportData?.risk_distribution || []}
          isLoading={reportLoading}
        />
        <DepartmentHeatmap
          data={reportData?.department_heatmap || []}
          isLoading={reportLoading}
        />
      </div>

      {/* Demographic Analysis */}
      <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4" style={{ animationDelay: '300ms' }}>
        <div>
          <h2 className="text-xl font-semibold text-gray-900">{translations.dashboard.demographicAnalysis}</h2>
          <p className="text-sm text-muted-foreground mt-1">
            {translations.dashboard.demographicDesc}
          </p>
        </div>
        
        {/* Distribution Charts */}
        <div>
          <div className="grid gap-6 md:grid-cols-2">
            <EnhancedDemographicChart
              title={translations.dashboard.age}
              description={translations.dashboard.ageDesc}
              icon={Calendar}
              data={reportData?.age_distribution || []}
              isLoading={reportLoading}
              chartType="donut"
            />
            <EnhancedDemographicChart
              title={translations.dashboard.maritalStatus}
              description={translations.dashboard.maritalStatusDesc}
              icon={Heart}
              data={reportData?.marital_status_distribution || []}
              isLoading={reportLoading}
              chartType="donut"
            />
            <EnhancedDemographicChart
              title={translations.dashboard.shiftType}
              description={translations.dashboard.shiftTypeDesc}
              icon={Clock}
              data={reportData?.shift_type_distribution || []}
              isLoading={reportLoading}
              chartType="donut"
            />
            <EnhancedDemographicChart
              title={translations.dashboard.experience}
              description={translations.dashboard.experienceDesc}
              icon={Briefcase}
              data={reportData?.experience_distribution || []}
              isLoading={reportLoading}
              chartType="donut"
            />
          </div>
        </div>

        {/* Risk Correlation Charts */}
        {((reportData?.age_risk_distribution?.length ?? 0) > 0 || (reportData?.shift_risk_distribution?.length ?? 0) > 0) && (
          <div>
            <h3 className="text-lg font-medium text-gray-800 mb-4">{translations.dashboard.riskByDemographics}</h3>
            <div className="grid gap-6 md:grid-cols-2">
              {reportData?.age_risk_distribution && reportData.age_risk_distribution.length > 0 && (
                <DemographicRiskChart
                  title={translations.dashboard.riskByAge}
                  description={translations.dashboard.riskByAgeDesc}
                  icon={Calendar}
                  data={reportData.age_risk_distribution}
                  isLoading={reportLoading}
                />
              )}
              {reportData?.shift_risk_distribution && reportData.shift_risk_distribution.length > 0 && (
                <DemographicRiskChart
                  title={translations.dashboard.riskByShift}
                  description={translations.dashboard.riskByShiftDesc}
                  icon={Clock}
                  data={reportData.shift_risk_distribution}
                  isLoading={reportLoading}
                />
              )}
            </div>
          </div>
        )}
      </div>

      {/* Quick Actions */}
      <Card className="animate-in fade-in slide-in-from-bottom-4 shadow-sm hover:shadow-md transition-all" style={{ animationDelay: '400ms' }}>
        <CardHeader>
          <CardTitle>{translations.dashboard.quickActions}</CardTitle>
          <CardDescription>
            {translations.dashboard.quickActionsDesc}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-6">
            <Button 
              variant="outline" 
              className="flex items-center space-x-2"
              onClick={() => router.push("/dashboard/staff")}
            >
              <Users className="h-4 w-4" />
              <span>{translations.dashboard.manageStaff}</span>
            </Button>
            <Button 
              variant="outline" 
              className="flex items-center space-x-2"
              onClick={() => router.push("/dashboard/assessments")}
            >
              <FileText className="h-4 w-4" />
              <span>{translations.dashboard.viewReports}</span>
            </Button>
            <AssessmentWizard
              trigger={
                <Button variant="outline" className="flex items-center space-x-2">
                  <Plus className="h-4 w-4" />
                  <span>{translations.dashboard.newAssessmentCycle}</span>
                </Button>
              }
            />
          </div>
        </CardContent>
      </Card>
    </div>
  );
}


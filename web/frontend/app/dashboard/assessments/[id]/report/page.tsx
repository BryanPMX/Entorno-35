"use client";

import { useParams, useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { ArrowLeft, FileText, AlertTriangle, CheckCircle, Clock, User, Building, Calendar, Target, Download, RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { reportService, IndividualReportDTO } from "@/services/report.service";

/**
 * Individual Assessment Report Page
 *
 * Displays detailed NOM-035 compliance report for a completed assessment.
 * Shows risk analysis, category scores, recommendations, and staff details.
 */
export default function AssessmentReportPage() {
  const params = useParams();
  const router = useRouter();
  const assessmentId = params.id as string;

  const handleDownloadPDF = async () => {
    if (!report) return;

    try {
      const pdfBlob = await reportService.downloadIndividualReportPDF(assessmentId);
      const url = window.URL.createObjectURL(pdfBlob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `NOM035_Report_${report.staff_name}_${new Date().toISOString().split('T')[0]}.pdf`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Failed to download PDF:', error);
      // You might want to show a toast notification here
    }
  };

  // Fetch individual report data with more specific cache key
  const { data: report, isLoading, error, refetch } = useQuery({
    queryKey: ["individual-report", assessmentId],
    queryFn: () => reportService.getIndividualReport(assessmentId),
    enabled: !!assessmentId,
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes
  });

  const getRiskBadge = (riskLevel: string) => {
    const riskConfig = {
      nulo: { label: "Nulo", color: "bg-green-100 text-green-800", description: "Very Low Risk" },
      bajo: { label: "Bajo", color: "bg-blue-100 text-blue-800", description: "Low Risk" },
      medio: { label: "Medio", color: "bg-yellow-100 text-yellow-800", description: "Medium Risk" },
      alto: { label: "Alto", color: "bg-orange-100 text-orange-800", description: "High Risk" },
      muy_alto: { label: "Muy Alto", color: "bg-red-100 text-red-800", description: "Very High Risk" },
    };

    const config = riskConfig[riskLevel as keyof typeof riskConfig];
    if (!config) return <Badge variant="outline">{riskLevel}</Badge>;

    return (
      <Badge className={`${config.color} border-0 text-sm px-3 py-1`}>
        {config.label}
      </Badge>
    );
  };

  const getRiskDescription = (riskLevel: string) => {
    const descriptions = {
      nulo: "No significant psychosocial risk factors detected",
      bajo: "Minimal psychosocial risk factors present",
      medio: "Moderate psychosocial risk factors requiring attention",
      alto: "High psychosocial risk factors requiring intervention",
      muy_alto: "Severe psychosocial risk factors requiring immediate action",
    };
    return descriptions[riskLevel as keyof typeof descriptions] || "Risk level assessment";
  };

  const formatScore = (score: number, maxScore: number = 100) => {
    return `${score}/${maxScore}`;
  };

  // Use backend-provided risk levels instead of recalculating
  const getRiskColor = (riskLevel: string) => {
    switch (riskLevel) {
      case "nulo": return "bg-green-500";
      case "bajo": return "bg-blue-500";
      case "medio": return "bg-yellow-500";
      case "alto": return "bg-orange-500";
      case "muy_alto": return "bg-red-500";
      default: return "bg-gray-500";
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="container mx-auto px-4 py-8">
          <div className="max-w-4xl mx-auto space-y-6">
            <div className="animate-pulse">
              <div className="h-8 bg-gray-200 rounded w-1/3 mb-4"></div>
              <div className="h-4 bg-gray-200 rounded w-1/4"></div>
            </div>
            <div className="grid gap-6 md:grid-cols-2">
              <Card className="animate-pulse">
                <CardContent className="p-6">
                  <div className="h-20 bg-gray-200 rounded mb-4"></div>
                  <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                </CardContent>
              </Card>
              <Card className="animate-pulse">
                <CardContent className="p-6">
                  <div className="h-20 bg-gray-200 rounded mb-4"></div>
                  <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error || !report) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="container mx-auto px-4 py-8">
          <div className="max-w-4xl mx-auto">
            <Button
              variant="ghost"
              onClick={() => router.back()}
              className="mb-6"
            >
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Assessments
            </Button>

            <Alert variant="destructive">
              <AlertTriangle className="h-4 w-4" />
              <AlertDescription>
                {error?.message || "Failed to load assessment report. The assessment may not be completed yet."}
              </AlertDescription>
            </Alert>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto space-y-6">
          {/* Header */}
          <div className="flex items-center justify-between">
            <div>
              <div className="flex items-center space-x-4 mb-4">
                <Button
                  variant="ghost"
                  onClick={() => router.back()}
                >
                  <ArrowLeft className="h-4 w-4 mr-2" />
                  Back to Assessments
                </Button>
                <Button
                  onClick={() => refetch()}
                  disabled={isLoading}
                  variant="outline"
                  className="mr-2"
                >
                  <RefreshCw className={`h-4 w-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
                  Refresh
                </Button>
                <Button
                  onClick={handleDownloadPDF}
                  disabled={!report}
                  className="bg-blue-600 hover:bg-blue-700"
                >
                  <Download className="h-4 w-4 mr-2" />
                  Download PDF
                </Button>
              </div>
              <h1 className="text-3xl font-bold text-gray-900">Assessment Report</h1>
              <p className="text-gray-600 mt-2">
                NOM-035 Psychosocial Risk Assessment Results
              </p>
            </div>
            <div className="text-right">
              <div className="text-sm text-gray-500">Assessment ID</div>
              <div className="font-mono text-xs bg-gray-100 px-2 py-1 rounded">
                {assessmentId.slice(0, 8)}...
              </div>
            </div>
          </div>

          {/* Staff Information */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <User className="h-5 w-5" />
                <span>Staff Information</span>
              </CardTitle>
              <CardDescription>
                Staff details for this specific assessment period ({report.period})
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-3">
                  <div>
                    <div className="text-sm font-medium text-gray-500">Name</div>
                    <div className="text-lg font-semibold">{report.staff_name}</div>
                  </div>
                  <div>
                    <div className="text-sm font-medium text-gray-500">Department</div>
                    <div>{report.department || "Not specified"}</div>
                  </div>
                </div>
                <div className="space-y-3">
                  <div>
                    <div className="text-sm font-medium text-gray-500">Assessment Period</div>
                    <div className="flex items-center space-x-2">
                      <Calendar className="h-4 w-4" />
                      <span className="font-semibold">{report.period}</span>
                    </div>
                  </div>
                  <div>
                    <div className="text-sm font-medium text-gray-500">Shift</div>
                    <div>{report.shift || "Not specified"}</div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Risk Assessment */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Target className="h-5 w-5" />
                <span>Risk Assessment</span>
              </CardTitle>
              <CardDescription>
                Overall psychosocial risk evaluation based on NOM-035 methodology
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="flex items-center justify-between">
                <div>
                  <div className="text-2xl font-bold mb-1">{report.total_score}/100</div>
                  <div className="text-sm text-gray-600">Total Score</div>
                </div>
                <div className="text-right">
                  <div className="mb-2">{getRiskBadge(report.risk_level)}</div>
                  <div className="text-sm text-gray-600 max-w-xs">
                    {getRiskDescription(report.risk_level)}
                  </div>
                </div>
              </div>

              <div className="w-full bg-gray-200 rounded-full h-3">
                <div
                  className={`h-3 rounded-full transition-all duration-500 ${getRiskColor(report.risk_level)}`}
                  style={{ width: `${Math.min((report.total_score / 100) * 100, 100)}%` }}
                ></div>
              </div>

              {report.requires_medical_attention && (
                <Alert variant="destructive">
                  <AlertTriangle className="h-4 w-4" />
                  <AlertDescription>
                    This assessment indicates a need for medical attention. Please consult with occupational health services.
                  </AlertDescription>
                </Alert>
              )}
            </CardContent>
          </Card>

          {/* Category Scores */}
          <Card>
            <CardHeader>
              <CardTitle>Category Analysis</CardTitle>
              <CardDescription>
                Breakdown by NOM-035 psychosocial risk categories
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {Object.entries(report.category_scores).map(([category, score]) => {
                  const riskLevel = report.category_risk_levels?.[category] || "nulo";
                  const maxScore = report.category_max_scores?.[category] || 20; // Fallback to 20 if not provided
                  return (
                    <div key={category} className="flex items-center justify-between">
                      <div className="flex-1">
                        <div className="flex items-center justify-between mb-2">
                          <div className="font-medium">{category}</div>
                          {getRiskBadge(riskLevel)}
                        </div>
                        <div className="w-full bg-gray-200 rounded-full h-2 mt-1">
                          <div
                            className={`h-2 rounded-full transition-all duration-500 ${getRiskColor(riskLevel)}`}
                            style={{ width: `${Math.min((score / maxScore) * 100, 100)}%` }}
                          ></div>
                        </div>
                      </div>
                      <div className="ml-4 text-right">
                        <div className="font-semibold">{score.toFixed(1)}</div>
                        <div className="text-xs text-gray-500">/{maxScore}</div>
                      </div>
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>

          {/* Domain Scores */}
          <Card>
            <CardHeader>
              <CardTitle>Domain Analysis</CardTitle>
              <CardDescription>
                Detailed breakdown by psychosocial risk domains
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {Object.entries(report.domain_scores).map(([domain, score]) => {
                  const riskLevel = report.domain_risk_levels?.[domain] || "nulo";
                  const maxScore = report.domain_max_scores?.[domain] || 15; // Fallback to 15 if not provided
                  return (
                    <div key={domain} className="flex items-center justify-between">
                      <div className="flex-1">
                        <div className="flex items-center justify-between mb-2">
                          <div className="font-medium">{domain}</div>
                          {getRiskBadge(riskLevel)}
                        </div>
                        <div className="w-full bg-gray-200 rounded-full h-2 mt-1">
                          <div
                            className={`h-2 rounded-full transition-all duration-500 ${getRiskColor(riskLevel)}`}
                            style={{ width: `${Math.min((score / maxScore) * 100, 100)}%` }}
                          ></div>
                        </div>
                      </div>
                      <div className="ml-4 text-right">
                        <div className="font-semibold">{score.toFixed(1)}</div>
                        <div className="text-xs text-gray-500">/{maxScore}</div>
                      </div>
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>

          {/* Recommendations */}
          {report.recommendations && report.recommendations.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <CheckCircle className="h-5 w-5" />
                  <span>Recommendations</span>
                </CardTitle>
                <CardDescription>
                  Actions to address identified psychosocial risk factors
                </CardDescription>
              </CardHeader>
              <CardContent>
                <ul className="space-y-3">
                  {report.recommendations.map((recommendation, index) => (
                    <li key={index} className="flex items-start space-x-3">
                      <div className="flex-shrink-0 w-6 h-6 bg-blue-100 rounded-full flex items-center justify-center mt-0.5">
                        <span className="text-xs font-semibold text-blue-600">{index + 1}</span>
                      </div>
                      <p className="text-gray-700">{recommendation}</p>
                    </li>
                  ))}
                </ul>
              </CardContent>
            </Card>
          )}

          {/* Assessment Metadata */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <FileText className="h-5 w-5" />
                <span>Assessment Timeline</span>
              </CardTitle>
              <CardDescription>
                Important dates and context for this assessment
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                <div>
                  <div className="font-medium text-gray-500">Guide Type</div>
                  <div className="flex items-center space-x-2">
                    <span>NOM-035 {report.guide_type}</span>
                    <Badge variant="outline" className="text-xs">
                      {report.guide_type === "I" ? "Trauma" : report.guide_type === "II" ? "Risk Factors" : "Work Environment"}
                    </Badge>
                  </div>
                </div>
                <div>
                  <div className="font-medium text-gray-500">Assessment Period</div>
                  <div className="flex items-center space-x-2">
                    <Calendar className="h-4 w-4" />
                    <span>{report.period}</span>
                  </div>
                </div>
                <div>
                  <div className="font-medium text-gray-500">Status</div>
                  <div className="flex items-center space-x-2">
                    <CheckCircle className="h-4 w-4 text-green-500" />
                    <span>Completed</span>
                  </div>
                </div>
                <div>
                  <div className="font-medium text-gray-500">Completed At</div>
                  <div className="flex items-center space-x-2">
                    <Clock className="h-4 w-4" />
                    <span>
                      {report.completed_at
                        ? new Date(report.completed_at).toLocaleDateString("es-MX", {
                            year: "numeric",
                            month: "long",
                            day: "numeric",
                            hour: "2-digit",
                            minute: "2-digit",
                          })
                        : "Not completed"
                      }
                    </span>
                  </div>
                </div>
              </div>
              <div className="mt-4 p-3 bg-blue-50 border border-blue-200 rounded-md">
                <p className="text-sm text-blue-800">
                  <strong>Important:</strong> This report reflects the psychosocial risk assessment results for the {report.period} period only.
                  If this staff member has multiple assessments, each report is independent and specific to its assessment period.
                </p>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
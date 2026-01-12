"use client";

import { useParams, useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { ArrowLeft, FileText, AlertTriangle, CheckCircle, Clock, User, Building, Calendar, Target, Download } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { reportService, IndividualReportDTO } from "@/services/report.service";
import { translations, getRiskLevelLabel, getRiskLevelDescription } from "@/lib/translations";

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
  const { data: report, isLoading, error } = useQuery({
    queryKey: ["individual-report", assessmentId],
    queryFn: () => reportService.getIndividualReport(assessmentId),
    enabled: !!assessmentId,
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes
  });

  const getRiskBadge = (riskLevel: string) => {
    const riskConfig = {
      nulo: { color: "bg-emerald-100 text-emerald-800 border-emerald-200" },
      bajo: { color: "bg-sky-100 text-sky-800 border-sky-200" },
      medio: { color: "bg-amber-100 text-amber-800 border-amber-200" },
      alto: { color: "bg-orange-100 text-orange-800 border-orange-200" },
      muy_alto: { color: "bg-red-100 text-red-800 border-red-200" },
    };

    const config = riskConfig[riskLevel as keyof typeof riskConfig];
    if (!config) return <Badge variant="outline">{riskLevel}</Badge>;

    return (
      <Badge className={`${config.color} border text-sm px-3 py-1 font-medium`}>
        {getRiskLevelLabel(riskLevel)}
      </Badge>
    );
  };

  const formatScore = (score: number, maxScore: number = 100) => {
    return `${score}/${maxScore}`;
  };

  // Use backend-provided risk levels instead of recalculating
  const getRiskColor = (riskLevel: string) => {
    switch (riskLevel) {
      case "nulo": return "bg-emerald-500";
      case "bajo": return "bg-sky-500";
      case "medio": return "bg-amber-500";
      case "alto": return "bg-orange-500";
      case "muy_alto": return "bg-red-600";
      default: return "bg-gray-400";
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
              {translations.reports.backToAssessments}
            </Button>
                <Button
                  onClick={handleDownloadPDF}
                  disabled={!report}
                  className="bg-blue-600 hover:bg-blue-700"
                >
                  <Download className="h-4 w-4 mr-2" />
                  {translations.reports.downloadPDF}
                </Button>
              </div>
              <h1 className="text-3xl font-bold text-gray-900">{translations.reports.title}</h1>
              <p className="text-gray-600 mt-2">
                {translations.reports.subtitle}
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
                <span>{translations.reports.staffInfo}</span>
              </CardTitle>
              <CardDescription>
                {translations.reports.assessmentPeriod}: {report.period}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-3">
                  <div>
                    <div className="text-sm font-medium text-gray-500">Nombre</div>
                    <div className="text-lg font-semibold">{report.staff_name}</div>
                  </div>
                  <div>
                    <div className="text-sm font-medium text-gray-500">Departamento</div>
                    <div>{report.department || "No especificado"}</div>
                  </div>
                </div>
                <div className="space-y-3">
                  <div>
                    <div className="text-sm font-medium text-gray-500">{translations.reports.assessmentPeriod}</div>
                    <div className="flex items-center space-x-2">
                      <Calendar className="h-4 w-4" />
                      <span className="font-semibold">{report.period}</span>
                    </div>
                  </div>
                  <div>
                    <div className="text-sm font-medium text-gray-500">Turno</div>
                    <div>{report.shift || "No especificado"}</div>
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
                <span>{translations.reports.riskAssessment}</span>
              </CardTitle>
              <CardDescription>
                {translations.reports.subtitle}
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="flex items-center justify-between">
                <div>
                  <div className="text-2xl font-bold mb-1">
                    {report.total_score.toFixed(1)} / {report.total_max_score.toFixed(0)}
                  </div>
                  <div className="text-sm text-gray-600">{translations.reports.totalScore}</div>
                </div>
                <div className="text-right">
                  <div className="mb-2">{getRiskBadge(report.risk_level)}</div>
                  <div className="text-sm text-gray-600 max-w-xs">
                    {getRiskLevelDescription(report.risk_level)}
                  </div>
                </div>
              </div>

              <div className="w-full bg-gray-200 rounded-full h-3">
                <div
                  className={`h-3 rounded-full transition-all duration-500 ${getRiskColor(report.risk_level)}`}
                  style={{ width: `${Math.min((report.total_score / report.total_max_score) * 100, 100)}%` }}
                ></div>
              </div>

              {report.requires_medical_attention && (
                <Alert variant="destructive">
                  <AlertTriangle className="h-4 w-4" />
                  <AlertDescription>
                    {translations.reports.medicalAttention}
                  </AlertDescription>
                </Alert>
              )}
            </CardContent>
          </Card>

          {/* Category Scores */}
          <Card>
            <CardHeader>
              <CardTitle>{translations.reports.categoryAnalysis}</CardTitle>
              <CardDescription>
                Analisis por categorias de riesgo psicosocial NOM-035
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-6">
                {Object.entries(report.category_scores).map(([category, score], index) => {
                  const riskLevel = report.category_risk_levels?.[category] || "nulo";
                  const maxScore = report.category_max_scores?.[category] || 20; // Fallback to 20 if not provided
                  const percentage = Math.min((score / maxScore) * 100, 100);

                  return (
                    <div key={category} className="space-y-3">
                      {/* Header with category name and risk badge */}
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-3">
                          <span className="text-sm font-medium text-gray-600 w-6">{index + 1}.</span>
                          <span className="font-semibold text-gray-900">{category}</span>
                        </div>
                        {getRiskBadge(riskLevel)}
                      </div>

                      {/* Progress bar and score */}
                      <div className="flex items-center space-x-4">
                        <div className="flex-1">
                          <div className="w-full bg-gray-100 rounded-full h-3 overflow-hidden">
                            <div
                              className={`h-3 rounded-full transition-all duration-700 ease-out ${getRiskColor(riskLevel)} shadow-sm`}
                              style={{ width: `${percentage}%` }}
                            ></div>
                          </div>
                        </div>
                        <div className="flex items-baseline space-x-1 min-w-0">
                          <span className="text-lg font-bold text-gray-900 tabular-nums">
                            {score.toFixed(1)}
                          </span>
                          <span className="text-sm font-medium text-gray-500">
                            /{maxScore}
                          </span>
                          <span className="text-xs text-gray-400 ml-1">
                            ({percentage.toFixed(0)}%)
                          </span>
                        </div>
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
              <CardTitle>{translations.reports.domainAnalysis}</CardTitle>
              <CardDescription>
                Analisis detallado por dominios de riesgo psicosocial
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-6">
                {Object.entries(report.domain_scores).map(([domain, score], index) => {
                  const riskLevel = report.domain_risk_levels?.[domain] || "nulo";
                  const maxScore = report.domain_max_scores?.[domain] || 15; // Fallback to 15 if not provided
                  const percentage = Math.min((score / maxScore) * 100, 100);

                  return (
                    <div key={domain} className="space-y-3">
                      {/* Header with domain name and risk badge */}
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-3">
                          <span className="text-sm font-medium text-gray-600 w-6">{index + 1}.</span>
                          <span className="font-semibold text-gray-900">{domain}</span>
                        </div>
                        {getRiskBadge(riskLevel)}
                      </div>

                      {/* Progress bar and score */}
                      <div className="flex items-center space-x-4">
                        <div className="flex-1">
                          <div className="w-full bg-gray-100 rounded-full h-3 overflow-hidden">
                            <div
                              className={`h-3 rounded-full transition-all duration-700 ease-out ${getRiskColor(riskLevel)} shadow-sm`}
                              style={{ width: `${percentage}%` }}
                            ></div>
                          </div>
                        </div>
                        <div className="flex items-baseline space-x-1 min-w-0">
                          <span className="text-lg font-bold text-gray-900 tabular-nums">
                            {score.toFixed(1)}
                          </span>
                          <span className="text-sm font-medium text-gray-500">
                            /{maxScore}
                          </span>
                          <span className="text-xs text-gray-400 ml-1">
                            ({percentage.toFixed(0)}%)
                          </span>
                        </div>
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
                  <span>{translations.reports.recommendations}</span>
                </CardTitle>
                <CardDescription>
                  Acciones para atender los factores de riesgo psicosocial identificados
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
                <span>{translations.reports.assessmentTimeline}</span>
              </CardTitle>
              <CardDescription>
                Fechas y contexto importante de esta evaluacion
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                <div>
                  <div className="font-medium text-gray-500">Tipo de Guia</div>
                  <div className="flex items-center space-x-2">
                    <span>NOM-035 {report.guide_type}</span>
                    <Badge variant="outline" className="text-xs">
                      {report.guide_type === "I" ? "Trauma" : report.guide_type === "II" ? "Factores de Riesgo" : "Entorno Laboral"}
                    </Badge>
                  </div>
                </div>
                <div>
                  <div className="font-medium text-gray-500">{translations.reports.assessmentPeriod}</div>
                  <div className="flex items-center space-x-2">
                    <Calendar className="h-4 w-4" />
                    <span>{report.period}</span>
                  </div>
                </div>
                <div>
                  <div className="font-medium text-gray-500">Estado</div>
                  <div className="flex items-center space-x-2">
                    <CheckCircle className="h-4 w-4 text-green-500" />
                    <span>Completado</span>
                  </div>
                </div>
                <div>
                  <div className="font-medium text-gray-500">Fecha de Finalizacion</div>
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
                        : "No completado"
                      }
                    </span>
                  </div>
                </div>
              </div>
              <div className="mt-4 p-3 bg-blue-50 border border-blue-200 rounded-md">
                <p className="text-sm text-blue-800">
                  <strong>Importante:</strong> Este reporte refleja los resultados de la evaluacion de riesgo psicosocial solo para el periodo {report.period}.
                  {translations.reports.multipleAssessmentsNote}
                </p>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
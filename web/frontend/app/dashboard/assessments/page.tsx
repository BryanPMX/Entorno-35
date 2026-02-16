"use client";

import { useState, useCallback } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { AssessmentTable } from "@/components/assessments/assessment-table";
import { assessmentService } from "@/services/assessment.service";
import { reportService } from "@/services/report.service";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { FileText } from "lucide-react";

/**
 * Pagina de Evaluaciones
 *
 * Interfaz de gestion de evaluaciones NOM-035.
 * Incluye listado, filtrado, creacion y seguimiento de estado.
 */
export default function AssessmentsPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [currentPage, setCurrentPage] = useState(1);
  const [limit, setLimit] = useState(25);
  const [filters, setFilters] = useState<{ status?: string; period?: number }>({});

  // Obtener evaluaciones con paginacion y filtros
  const { data: assessmentData, isLoading, refetch } = useQuery({
    queryKey: ["assessments", currentPage, limit, filters],
    queryFn: () =>
      assessmentService.list({
        limit,
        offset: (currentPage - 1) * limit,
        ...filters,
      }),
  });

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  const handleLimitChange = (newLimit: number) => {
    setLimit(newLimit);
    setCurrentPage(1);
  };

  const handleCreateAssessment = () => {
    // Invalidate all related queries to refresh dashboard and lists
    queryClient.invalidateQueries({ queryKey: ["general-report"] });
    queryClient.invalidateQueries({ queryKey: ["assessments"] });
    refetch();
  };

  const handleViewReport = (assessmentId: string) => {
    router.push(`/dashboard/assessments/${assessmentId}/report`);
  };

  const copyToClipboard = async (text: string): Promise<boolean> => {
    if (navigator.clipboard && window.isSecureContext) {
      try {
        await navigator.clipboard.writeText(text);
        return true;
      } catch (error) {
        console.warn("Error al copiar:", error);
      }
    }

    try {
      const textArea = document.createElement("textarea");
      textArea.value = text;
      textArea.style.position = "fixed";
      textArea.style.left = "-9999px";
      textArea.style.top = "-9999px";
      document.body.appendChild(textArea);
      textArea.focus();
      textArea.select();

      const successful = document.execCommand("copy");
      document.body.removeChild(textArea);

      return successful;
    } catch (error) {
      console.error("Error al copiar:", error);
      return false;
    }
  };

  const handleGenerateLink = async (assessmentId: string) => {
    try {
      const linkData = await assessmentService.createLink(assessmentId, 7);

      const fullLink = `${window.location.origin}/assessment/${linkData.token}`;
      const copied = await copyToClipboard(fullLink);

      if (copied) {
        toast.success("Enlace copiado al portapapeles", {
          description: "El enlace seguro ha sido generado y copiado.",
        });
      } else {
        toast.success("Enlace generado", {
          description: `Copia este enlace: ${fullLink}`,
          duration: 10000,
        });
      }

      refetch();
    } catch (error) {
      toast.error("Error al generar el enlace", {
        description: "Por favor intenta de nuevo o contacta a soporte.",
      });
      console.error("Error al generar enlace:", error);
    }
  };

  const handleSendEmail = async (assessmentId: string) => {
    try {
      await assessmentService.sendEmail(assessmentId);
      toast.success("Correo enviado exitosamente", {
        description: "El enlace de evaluacion ha sido enviado al correo del empleado.",
      });
      refetch();
    } catch (error: unknown) {
      const errorMessage =
        typeof error === "object" && error && "response" in error && (error as { response?: { data?: { error?: string } } }).response?.data?.error
          ? (error as { response: { data: { error: string } } }).response.data.error
          : "Error desconocido";
      toast.error("Error al enviar el correo", {
        description: errorMessage,
      });
      console.error("Error al enviar correo:", error);
    }
  };

  const handleFiltersChange = useCallback((newFilters: { status?: string; period?: number }) => {
    setFilters(newFilters);
    setCurrentPage(1);
  }, []);

  const [isDownloading, setIsDownloading] = useState(false);

  const handleDownloadGeneralReport = async () => {
    try {
      setIsDownloading(true);
      const pdfBlob = await reportService.downloadGeneralReportPDF(filters.period);
      const url = window.URL.createObjectURL(pdfBlob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `NOM035_Reporte_General_${new Date().toISOString().split('T')[0]}.pdf`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      toast.success("Reporte general descargado exitosamente");
    } catch (error) {
      console.error('Error al descargar reporte:', error);
      toast.error("Error al descargar el reporte general");
    } finally {
      setIsDownloading(false);
    }
  };

  return (
    <div className="page-content-shell page-content-shell-c space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Evaluaciones</h1>
          <p className="mt-2 text-muted-foreground">
            Crea y gestiona evaluaciones NOM-035 para tu personal
          </p>
        </div>
        <Button
          variant="outline"
          onClick={handleDownloadGeneralReport}
          disabled={isDownloading}
          className="flex items-center space-x-2 border-primary/20 bg-background/60 hover:bg-accent/70"
        >
          {isDownloading ? (
            <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
          ) : (
            <FileText className="h-4 w-4" />
          )}
          <span>Descargar Reporte General</span>
        </Button>
      </div>

      <AssessmentTable
        assessmentData={assessmentData || null}
        isLoading={isLoading}
        onPageChange={handlePageChange}
        onLimitChange={handleLimitChange}
        onCreateAssessment={handleCreateAssessment}
        currentPage={currentPage}
        limit={limit}
        onViewReport={handleViewReport}
        onGenerateLink={handleGenerateLink}
        onFiltersChange={handleFiltersChange}
        onSendEmail={handleSendEmail}
      />
    </div>
  );
}

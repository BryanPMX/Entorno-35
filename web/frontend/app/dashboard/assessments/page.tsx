"use client";

import { useState, useCallback } from "react";
import { useQuery } from "@tanstack/react-query";
import { AssessmentTable } from "@/components/assessments/assessment-table";
import { assessmentService } from "@/services/assessment.service";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

/**
 * Assessments Page
 *
 * Comprehensive assessment management interface for NOM-035 compliance.
 * Features assessment listing, filtering, creation, and status tracking.
 */
export default function AssessmentsPage() {
  const router = useRouter();
  const [currentPage, setCurrentPage] = useState(1);
  const [limit, setLimit] = useState(25);
  const [filters, setFilters] = useState<{ status?: string; period?: number }>({});

  // Fetch assessments with pagination and filters
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
    setCurrentPage(1); // Reset to first page when changing limit
  };

  const handleCreateAssessment = () => {
    refetch(); // Refresh the list after creating an assessment
  };

  const handleViewReport = (assessmentId: string) => {
    router.push(`/dashboard/assessments/${assessmentId}/report`);
  };

  const copyToClipboard = async (text: string): Promise<boolean> => {
    // Try modern clipboard API first
    if (navigator.clipboard && window.isSecureContext) {
      try {
        await navigator.clipboard.writeText(text);
        return true;
      } catch (error) {
        console.warn("Clipboard API failed:", error);
      }
    }

    // Fallback: Use textarea method
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
      console.error("Fallback clipboard method failed:", error);
      return false;
    }
  };

  const handleGenerateLink = async (assessmentId: string) => {
    try {
      const linkData = await assessmentService.createLink(assessmentId, 7); // 7 days expiry

      // Copy link to clipboard with fallback
      const fullLink = `${window.location.origin}/assessment/${linkData.token}`;
      const copied = await copyToClipboard(fullLink);

      if (copied) {
        toast.success("Assessment link copied to clipboard!", {
          description: "The secure link has been generated and copied to your clipboard.",
        });
      } else {
        // Fallback: Show link that user can copy manually
        toast.success("Assessment link generated!", {
          description: `Copy this link: ${fullLink}`,
          duration: 10000, // Show for 10 seconds
        });
      }

      refetch(); // Refresh to show any status changes
    } catch (error) {
      toast.error("Failed to generate assessment link", {
        description: "Please try again or contact support if the issue persists.",
      });
      console.error("Link generation error:", error);
    }
  };

  const handleFiltersChange = useCallback((newFilters: { status?: string; period?: number }) => {
    setFilters(newFilters);
    setCurrentPage(1); // Reset to first page when filters change
  }, []);

  return (
    <div className="space-y-8">
      <div>
        <h1 className="heading-1">Assessments</h1>
        <p className="label-muted mt-2">
          Create and manage NOM-035 assessments for your staff members
        </p>
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
      />
    </div>
  );
}


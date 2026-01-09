"use client";

import { useState } from "react";
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

  // Fetch assessments with pagination
  const { data: assessmentData, isLoading, refetch } = useQuery({
    queryKey: ["assessments", currentPage, limit],
    queryFn: () =>
      assessmentService.list({
        limit,
        offset: (currentPage - 1) * limit,
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

  const handleGenerateLink = async (assessmentId: string) => {
    try {
      const linkData = await assessmentService.createLink(assessmentId, 7); // 7 days expiry

      // Copy link to clipboard
      const fullLink = `${window.location.origin}/assessment/${linkData.token}`;
      await navigator.clipboard.writeText(fullLink);

      toast.success("Assessment link copied to clipboard!", {
        description: "The secure link has been generated and copied to your clipboard.",
      });

      refetch(); // Refresh to show any status changes
    } catch (error) {
      toast.error("Failed to generate assessment link");
      console.error("Link generation error:", error);
    }
  };

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
      />
    </div>
  );
}


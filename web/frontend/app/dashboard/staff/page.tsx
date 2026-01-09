"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { StaffTable } from "@/components/staff/staff-table";
import { CsvUploadModal } from "@/components/staff/csv-upload-modal";
import { StaffCreateDialog } from "@/components/staff/staff-create-dialog";
import { staffService } from "@/services/staff.service";
import type { PaginatedResponse, Staff } from "@/types/backend";

/**
 * Staff Management Page
 *
 * Comprehensive staff management with data table and CSV import functionality.
 */
export default function StaffPage() {
  const [currentPage, setCurrentPage] = useState(1);
  const [limit, setLimit] = useState(50);
  const [isUploadModalOpen, setIsUploadModalOpen] = useState(false);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);

  // Fetch staff data with pagination
  const {
    data: staffData,
    isLoading,
    refetch,
  } = useQuery<PaginatedResponse<Staff>>({
    queryKey: ["staff", currentPage, limit],
    queryFn: () => staffService.getAll({ limit, offset: (currentPage - 1) * limit }),
  });

  // Handle page changes
  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  // Handle limit changes
  const handleLimitChange = (newLimit: number) => {
    setLimit(newLimit);
    setCurrentPage(1); // Reset to first page when changing limit
  };

  // Handle CSV upload completion
  const handleUploadComplete = () => {
    setIsUploadModalOpen(false);
    refetch(); // Refresh the staff data
  };

  // Handle CSV upload modal open
  const handleCsvUpload = () => {
    setIsUploadModalOpen(true);
  };

  // Handle manual staff creation
  const handleCreateStaff = () => {
    setIsCreateDialogOpen(true);
  };

  // Handle staff creation completion
  const handleStaffCreated = () => {
    setIsCreateDialogOpen(false);
    refetch(); // Refresh the staff data
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Staff Management</h1>
        <p className="text-gray-600 mt-2">
          Manage your staff members, view their information, and import via CSV
        </p>
      </div>

      <StaffTable
        staffData={staffData || null}
        isLoading={isLoading}
        onPageChange={handlePageChange}
        onLimitChange={handleLimitChange}
        onCsvUpload={handleCsvUpload}
        onCreateStaff={handleCreateStaff}
        currentPage={currentPage}
        limit={limit}
      />

      <CsvUploadModal
        isOpen={isUploadModalOpen}
        onClose={() => setIsUploadModalOpen(false)}
        onUploadComplete={handleUploadComplete}
      />

      <StaffCreateDialog
        onStaffCreated={handleStaffCreated}
      />
    </div>
  );
}


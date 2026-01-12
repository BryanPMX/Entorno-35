"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { StaffTable } from "@/components/staff/staff-table";
import { CsvUploadModal } from "@/components/staff/csv-upload-modal";
import { staffService } from "@/services/staff.service";
import type { PaginatedResponse, Staff } from "@/types/backend";

/**
 * Pagina de Gestion de Personal
 *
 * Administracion integral del personal con tabla de datos y funcionalidad de importacion CSV.
 */
export default function StaffPage() {
  const [currentPage, setCurrentPage] = useState(1);
  const [limit, setLimit] = useState(50);
  const [isUploadModalOpen, setIsUploadModalOpen] = useState(false);

  // Obtener datos del personal con paginacion
  const {
    data: staffData,
    isLoading,
    refetch,
  } = useQuery<PaginatedResponse<Staff>>({
    queryKey: ["staff", currentPage, limit],
    queryFn: () => staffService.getAll({ limit, offset: (currentPage - 1) * limit }),
  });

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  const handleLimitChange = (newLimit: number) => {
    setLimit(newLimit);
    setCurrentPage(1);
  };

  const handleUploadComplete = () => {
    setIsUploadModalOpen(false);
    refetch();
  };

  const handleCsvUpload = () => {
    setIsUploadModalOpen(true);
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Gestion de Personal</h1>
        <p className="text-gray-600 mt-2">
          Administra los miembros del personal, visualiza su informacion e importa mediante CSV
        </p>
      </div>

      <StaffTable
        staffData={staffData || null}
        isLoading={isLoading}
        onPageChange={handlePageChange}
        onLimitChange={handleLimitChange}
        onCsvUpload={handleCsvUpload}
        onDataRefresh={refetch}
        currentPage={currentPage}
        limit={limit}
      />

      <CsvUploadModal
        isOpen={isUploadModalOpen}
        onClose={() => setIsUploadModalOpen(false)}
        onUploadComplete={handleUploadComplete}
      />
    </div>
  );
}

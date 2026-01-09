"use client";

import React, { useState } from "react";
import { ChevronUp, ChevronDown, Upload, Users } from "lucide-react";
import { Button } from "@/components/ui/button";
import { StaffCreateDialog } from "./staff-create-dialog";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import type { PaginatedResponse, Staff } from "@/types/backend";

interface StaffTableProps {
  staffData: PaginatedResponse<Staff> | null;
  isLoading: boolean;
  onPageChange: (page: number) => void;
  onLimitChange: (limit: number) => void;
  onCsvUpload: () => void;
  currentPage: number;
  limit: number;
}

type SortField = "full_name" | "email" | "created_at";
type SortDirection = "asc" | "desc";

export function StaffTable({
  staffData,
  isLoading,
  onPageChange,
  onLimitChange,
  onCsvUpload,
  currentPage,
  limit,
}: StaffTableProps) {
  const [sortField, setSortField] = useState<SortField>("created_at");
  const [sortDirection, setSortDirection] = useState<SortDirection>("desc");
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);

  const handleSort = (field: SortField) => {
    if (sortField === field) {
      setSortDirection(sortDirection === "asc" ? "desc" : "asc");
    } else {
      setSortField(field);
      setSortDirection("asc");
    }
    // In a real implementation, you'd call an API with sort parameters
  };

  const SortButton = ({ field, children }: { field: SortField; children: React.ReactNode }) => (
    <Button
      variant="ghost"
      size="sm"
      onClick={() => handleSort(field)}
      className="h-auto p-0 font-medium hover:bg-transparent"
    >
      {children}
      {sortField === field && (
        sortDirection === "asc" ? (
          <ChevronUp className="ml-1 h-4 w-4" />
        ) : (
          <ChevronDown className="ml-1 h-4 w-4" />
        )
      )}
    </Button>
  );

  const totalPages = staffData ? Math.ceil(staffData.total / limit) : 0;
  const startItem = staffData ? (currentPage - 1) * limit + 1 : 0;
  const endItem = staffData ? Math.min(currentPage * limit, staffData.total) : 0;

  return (
    <div className="space-y-6">
      {/* Header with stats and upload button */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-2">
            <Users className="h-5 w-5 text-muted-foreground" />
            <div>
              <p className="text-sm font-medium">
                {staffData?.total || 0} Staff Members
              </p>
              <p className="text-xs text-muted-foreground">
                Showing {startItem}-{endItem} of {staffData?.total || 0}
              </p>
            </div>
          </div>
        </div>
        <div className="flex space-x-2">
          <Button onClick={() => setIsCreateDialogOpen(true)} variant="default">
            <Users className="h-4 w-4 mr-2" />
            Add Staff Member
          </Button>
          <Button onClick={onCsvUpload} variant="outline" className="flex items-center space-x-2">
            <Upload className="h-4 w-4" />
            <span>Import CSV</span>
          </Button>
        </div>
      </div>

      {/* Staff Table */}
      <Card>
        <CardHeader>
          <CardTitle>Staff Members</CardTitle>
          <CardDescription>
            Manage your staff members and their information
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>
                    <SortButton field="full_name">Name</SortButton>
                  </TableHead>
                  <TableHead>Identifier</TableHead>
                  <TableHead>
                    <SortButton field="email">Email</SortButton>
                  </TableHead>
                  <TableHead>Department</TableHead>
                  <TableHead>Role</TableHead>
                  <TableHead>
                    <SortButton field="created_at">Created</SortButton>
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center py-8">
                      <div className="flex items-center justify-center space-x-2">
                        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-primary"></div>
                        <span>Loading staff members...</span>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : staffData?.data && staffData.data.length > 0 ? (
                  staffData.data.map((staff) => (
                    <TableRow key={staff.id}>
                      <TableCell className="font-medium">
                        {staff.full_name}
                      </TableCell>
                      <TableCell className="font-mono text-sm">
                        {staff.curp || staff.employee_id}
                      </TableCell>
                      <TableCell>{staff.email || "—"}</TableCell>
                      <TableCell>
                        {staff.demographics?.department || "—"}
                      </TableCell>
                      <TableCell>
                        {staff.demographics?.role || "—"}
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {new Date(staff.created_at).toLocaleDateString()}
                      </TableCell>
                    </TableRow>
                  ))
                ) : (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                      No staff members found. Import staff data using the CSV upload feature.
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>

          {/* Pagination */}
          {staffData && staffData.total > limit && (
            <div className="flex items-center justify-between pt-4">
              <div className="flex items-center space-x-2">
                <span className="text-sm text-muted-foreground">Show</span>
                <select
                  value={limit}
                  onChange={(e) => onLimitChange(Number(e.target.value))}
                  className="h-8 w-16 rounded border border-input bg-background px-3 py-1 text-sm"
                >
                  <option value={10}>10</option>
                  <option value={25}>25</option>
                  <option value={50}>50</option>
                  <option value={100}>100</option>
                </select>
                <span className="text-sm text-muted-foreground">per page</span>
              </div>

              <Pagination>
                <PaginationContent>
                  <PaginationItem>
                    <PaginationPrevious
                      onClick={() => currentPage > 1 && onPageChange(currentPage - 1)}
                      className={
                        currentPage <= 1 ? "pointer-events-none opacity-50" : "cursor-pointer"
                      }
                    />
                  </PaginationItem>

                  {/* Page numbers */}
                  {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                    const pageNumber = Math.max(1, Math.min(totalPages - 4, currentPage - 2)) + i;
                    if (pageNumber > totalPages) return null;

                    return (
                      <PaginationItem key={pageNumber}>
                        <PaginationLink
                          onClick={() => onPageChange(pageNumber)}
                          isActive={pageNumber === currentPage}
                          className="cursor-pointer"
                        >
                          {pageNumber}
                        </PaginationLink>
                      </PaginationItem>
                    );
                  })}

                  <PaginationItem>
                    <PaginationNext
                      onClick={() => currentPage < totalPages && onPageChange(currentPage + 1)}
                      className={
                        currentPage >= totalPages ? "pointer-events-none opacity-50" : "cursor-pointer"
                      }
                    />
                  </PaginationItem>
                </PaginationContent>
              </Pagination>
            </div>
          )}
        </CardContent>
      </Card>

      <StaffCreateDialog
        onStaffCreated={() => setIsCreateDialogOpen(false)}
      />
    </div>
  );
}
"use client";

import React, { useState, useEffect } from "react";
import { ChevronUp, ChevronDown, Eye, Link, Plus, FileText, AlertTriangle, CheckCircle, Clock } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
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
import { AssessmentWizard } from "./assessment-wizard";
import type { PaginatedResponse, Assessment } from "@/types/backend";

interface AssessmentTableProps {
  assessmentData: PaginatedResponse<Assessment> | null;
  isLoading: boolean;
  onPageChange: (page: number) => void;
  onLimitChange: (limit: number) => void;
  onCreateAssessment: () => void;
  currentPage: number;
  limit: number;
  onViewReport?: (assessmentId: string) => void;
  onGenerateLink?: (assessmentId: string) => void;
  onFiltersChange?: (filters: { status?: string; period?: number }) => void;
}

type SortField = "created_at" | "staff_name" | "status" | "risk_level" | "period";
type SortDirection = "asc" | "desc";

const statusConfig = {
  pending: { label: "Pending", color: "bg-yellow-100 text-yellow-800", icon: Clock },
  completed: { label: "Completed", color: "bg-green-100 text-green-800", icon: CheckCircle },
  cancelled: { label: "Cancelled", color: "bg-gray-100 text-gray-800", icon: AlertTriangle },
};

const riskLevelConfig = {
  nulo: { label: "Nulo", color: "bg-green-100 text-green-800" },
  bajo: { label: "Bajo", color: "bg-blue-100 text-blue-800" },
  medio: { label: "Medio", color: "bg-yellow-100 text-yellow-800" },
  alto: { label: "Alto", color: "bg-orange-100 text-orange-800" },
  muy_alto: { label: "Muy Alto", color: "bg-red-100 text-red-800" },
};

export function AssessmentTable({
  assessmentData,
  isLoading,
  onPageChange,
  onLimitChange,
  onCreateAssessment,
  currentPage,
  limit,
  onViewReport,
  onGenerateLink,
  onFiltersChange,
}: AssessmentTableProps) {
  const [sortField, setSortField] = useState<SortField>("created_at");
  const [sortDirection, setSortDirection] = useState<SortDirection>("desc");
  const [statusFilter, setStatusFilter] = useState<string>("all");
  const [periodFilter, setPeriodFilter] = useState<string>("all");

  // Notify parent component when filters change
  useEffect(() => {
    if (onFiltersChange) {
      const filters: { status?: string; period?: number } = {};

      if (statusFilter !== "all") {
        filters.status = statusFilter;
      }

      if (periodFilter !== "all") {
        filters.period = parseInt(periodFilter);
      }

      onFiltersChange(filters);
    }
  }, [statusFilter, periodFilter, onFiltersChange]);

  const handleSort = (field: SortField) => {
    if (sortField === field) {
      setSortDirection(sortDirection === "asc" ? "desc" : "asc");
    } else {
      setSortField(field);
      setSortDirection("asc");
    }
  };

  const SortableHeader = ({
    field,
    children,
  }: {
    field: SortField;
    children: React.ReactNode;
  }) => (
    <TableHead>
      <Button
        variant="ghost"
        className="h-auto p-0 font-medium hover:bg-transparent"
        onClick={() => handleSort(field)}
      >
        {children}
        {sortField === field && (
          <span className="ml-1">
            {sortDirection === "asc" ? (
              <ChevronUp className="h-3 w-3" />
            ) : (
              <ChevronDown className="h-3 w-3" />
            )}
          </span>
        )}
      </Button>
    </TableHead>
  );

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("es-MX", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  const getStatusBadge = (status: string) => {
    const config = statusConfig[status as keyof typeof statusConfig];
    if (!config) return <Badge variant="secondary">{status}</Badge>;

    const IconComponent = config.icon;
    return (
      <Badge className={`${config.color} border-0`}>
        <IconComponent className="h-3 w-3 mr-1" />
        {config.label}
      </Badge>
    );
  };

  const getRiskBadge = (riskLevel?: string) => {
    if (!riskLevel) return <Badge variant="outline">-</Badge>;

    const config = riskLevelConfig[riskLevel as keyof typeof riskLevelConfig];
    if (!config) return <Badge variant="outline">{riskLevel}</Badge>;

    return (
      <Badge className={`${config.color} border-0`}>
        {config.label}
      </Badge>
    );
  };

  const getStaffName = (assessment: Assessment) => {
    if (assessment.staff) {
      return assessment.staff.full_name;
    }
    return `Staff ${assessment.staff_id.slice(-4)}`;
  };

  return (
    <div className="space-y-6">
      {/* Header with filters and actions */}
      <div className="flex flex-col sm:flex-row gap-4 justify-between items-start sm:items-center">
        <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center">
          <div className="flex items-center gap-2">
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-32">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="pending">Pending</SelectItem>
                <SelectItem value="completed">Completed</SelectItem>
                <SelectItem value="cancelled">Cancelled</SelectItem>
              </SelectContent>
            </Select>

            <Select value={periodFilter} onValueChange={setPeriodFilter}>
              <SelectTrigger className="w-24">
                <SelectValue placeholder="Period" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Years</SelectItem>
                <SelectItem value="2025">2025</SelectItem>
                <SelectItem value="2024">2024</SelectItem>
                <SelectItem value="2023">2023</SelectItem>
              </SelectContent>
            </Select>
          </div>
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

      {/* Table */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <FileText className="h-5 w-5" />
            <span>Assessments</span>
          </CardTitle>
          <CardDescription>
            Manage and track NOM-035 assessments for your organization
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <SortableHeader field="staff_name">Staff Member</SortableHeader>
                  <SortableHeader field="period">Period</SortableHeader>
                  <SortableHeader field="status">Status</SortableHeader>
                  <SortableHeader field="risk_level">Risk Level</SortableHeader>
                  <SortableHeader field="created_at">Created</SortableHeader>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  // Loading skeleton
                  Array.from({ length: 5 }).map((_, index) => (
                    <TableRow key={index}>
                      <TableCell>
                        <div className="h-4 bg-gray-200 rounded animate-pulse w-32" />
                      </TableCell>
                      <TableCell>
                        <div className="h-4 bg-gray-200 rounded animate-pulse w-16" />
                      </TableCell>
                      <TableCell>
                        <div className="h-5 bg-gray-200 rounded animate-pulse w-20" />
                      </TableCell>
                      <TableCell>
                        <div className="h-5 bg-gray-200 rounded animate-pulse w-16" />
                      </TableCell>
                      <TableCell>
                        <div className="h-4 bg-gray-200 rounded animate-pulse w-20" />
                      </TableCell>
                      <TableCell>
                        <div className="flex gap-2">
                          <div className="h-8 bg-gray-200 rounded animate-pulse w-16" />
                          <div className="h-8 bg-gray-200 rounded animate-pulse w-16" />
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                ) : assessmentData?.data.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center py-8">
                      <div className="flex flex-col items-center gap-2">
                        <FileText className="h-8 w-8 text-gray-400" />
                        <p className="text-gray-500">No assessments found</p>
                        <p className="text-sm text-gray-400">
                          Create your first assessment to get started
                        </p>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  assessmentData?.data.map((assessment) => (
                    <TableRow key={assessment.id}>
                      <TableCell className="font-medium">
                        {getStaffName(assessment)}
                      </TableCell>
                      <TableCell>{assessment.period}</TableCell>
                      <TableCell>{getStatusBadge(assessment.status)}</TableCell>
                      <TableCell>{getRiskBadge(assessment.risk_level)}</TableCell>
                      <TableCell className="text-sm text-gray-500">
                        {formatDate(assessment.created_at)}
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          {assessment.status === "pending" && (
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => onGenerateLink?.(assessment.id)}
                              className="flex items-center space-x-1"
                            >
                              <Link className="h-3 w-3" />
                              <span>Link</span>
                            </Button>
                          )}
                          {assessment.status === "completed" && (
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => onViewReport?.(assessment.id)}
                              className="flex items-center space-x-1"
                            >
                              <Eye className="h-3 w-3" />
                              <span>Report</span>
                            </Button>
                          )}
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>

          {/* Pagination */}
          {assessmentData && assessmentData.total > limit && (
            <div className="flex items-center justify-between px-2 py-4">
              <div className="flex items-center space-x-2 text-sm text-gray-500">
                <span>Show</span>
                <Select
                  value={limit.toString()}
                  onValueChange={(value) => onLimitChange(parseInt(value))}
                >
                  <SelectTrigger className="w-16">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="10">10</SelectItem>
                    <SelectItem value="25">25</SelectItem>
                    <SelectItem value="50">50</SelectItem>
                    <SelectItem value="100">100</SelectItem>
                  </SelectContent>
                </Select>
                <span>per page</span>
              </div>

              <Pagination>
                <PaginationContent>
                  <PaginationItem>
                    <PaginationPrevious
                      onClick={() => currentPage > 1 && onPageChange(currentPage - 1)}
                      className={
                        currentPage <= 1 ? "pointer-events-none opacity-50" : ""
                      }
                    />
                  </PaginationItem>

                  {Array.from(
                    { length: Math.ceil(assessmentData.total / limit) },
                    (_, i) => i + 1
                  )
                    .filter((page) => {
                      const totalPages = Math.ceil(assessmentData.total / limit);
                      if (totalPages <= 7) return true;
                      if (page === 1 || page === totalPages) return true;
                      if (Math.abs(page - currentPage) <= 1) return true;
                      return false;
                    })
                    .map((page, index, array) => (
                      <React.Fragment key={page}>
                        {index > 0 && array[index - 1] !== page - 1 && (
                          <PaginationItem>
                            <span className="px-2">...</span>
                          </PaginationItem>
                        )}
                        <PaginationItem>
                          <PaginationLink
                            onClick={() => onPageChange(page)}
                            isActive={page === currentPage}
                            className="cursor-pointer"
                          >
                            {page}
                          </PaginationLink>
                        </PaginationItem>
                      </React.Fragment>
                    ))}

                  <PaginationItem>
                    <PaginationNext
                      onClick={() =>
                        currentPage < Math.ceil(assessmentData.total / limit) &&
                        onPageChange(currentPage + 1)
                      }
                      className={
                        currentPage >= Math.ceil(assessmentData.total / limit)
                          ? "pointer-events-none opacity-50"
                          : ""
                      }
                    />
                  </PaginationItem>
                </PaginationContent>
              </Pagination>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
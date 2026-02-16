"use client";

import React, { useState, useEffect, useCallback, ReactNode } from "react";
import { ChevronUp, ChevronDown, Eye, Link, Plus, FileText, AlertTriangle, CheckCircle, Clock, Mail, MoreHorizontal, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
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
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { AssessmentWizard } from "./assessment-wizard";
import { AssessmentDeleteDialog } from "./assessment-delete-dialog";
import { getNom35RiskBadgeClass, isNom35RiskLevel } from "@/lib/nom35-risk";
import { getRiskLevelLabel } from "@/lib/translations";
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
  onSendEmail?: (assessmentId: string) => void;
}

type SortField = "created_at" | "staff_name" | "status" | "risk_level" | "period";
type SortDirection = "asc" | "desc";

interface SortableHeaderProps {
  field: SortField;
  children: ReactNode;
  activeField: SortField;
  direction: SortDirection;
  onSort: (field: SortField) => void;
}

function SortableHeader({ field, children, activeField, direction, onSort }: SortableHeaderProps) {
  return (
    <TableHead>
      <Button
        variant="ghost"
        className="h-auto p-0 font-medium hover:bg-transparent"
        onClick={() => onSort(field)}
      >
        {children}
        {activeField === field && (
          <span className="ml-1">
            {direction === "asc" ? (
              <ChevronUp className="h-3 w-3" />
            ) : (
              <ChevronDown className="h-3 w-3" />
            )}
          </span>
        )}
      </Button>
    </TableHead>
  );
}

const statusConfig = {
  pending: { label: "Pendiente", color: "bg-amber-100 text-amber-800 border border-amber-200", icon: Clock },
  completed: { label: "Completado", color: "risk-level-nulo border", icon: CheckCircle },
  cancelled: { label: "Cancelado", color: "bg-muted text-muted-foreground border border-border", icon: AlertTriangle },
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
  onSendEmail,
}: AssessmentTableProps) {
  const [sortField, setSortField] = useState<SortField>("created_at");
  const [sortDirection, setSortDirection] = useState<SortDirection>("desc");
  const [statusFilter, setStatusFilter] = useState<string>("all");
  const [periodFilter, setPeriodFilter] = useState<string>("all");
  const [deletingAssessment, setDeletingAssessment] = useState<Assessment | null>(null);

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

  const handleSort = useCallback(
    (field: SortField) => {
      if (sortField === field) {
        setSortDirection((prev) => (prev === "asc" ? "desc" : "asc"));
      } else {
        setSortField(field);
        setSortDirection("asc");
      }
    },
    [sortField]
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
      <Badge className={config.color}>
        <IconComponent className="h-3 w-3 mr-1" />
        {config.label}
      </Badge>
    );
  };

  const getRiskBadge = (riskLevel?: string) => {
    if (!riskLevel) return <Badge variant="outline">-</Badge>;

    if (!isNom35RiskLevel(riskLevel)) return <Badge variant="outline">{riskLevel}</Badge>;

    return (
      <Badge className={getNom35RiskBadgeClass(riskLevel)}>
        {getRiskLevelLabel(riskLevel)}
      </Badge>
    );
  };

  const getStaffName = (assessment: Assessment) => {
    if (assessment.staff) {
      return assessment.staff.full_name;
    }
    return `Empleado ${assessment.staff_id.slice(-4)}`;
  };

  const getStaffEmail = (assessment: Assessment) => {
    return assessment.staff?.email || null;
  };

  return (
    <div className="space-y-6">
      {/* Header with filters and actions */}
      <div className="flex flex-col sm:flex-row gap-4 justify-between items-start sm:items-center">
        <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center">
          <div className="flex items-center gap-2">
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-36">
                <SelectValue placeholder="Estado" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos los Estados</SelectItem>
                <SelectItem value="pending">Pendiente</SelectItem>
                <SelectItem value="completed">Completado</SelectItem>
                <SelectItem value="cancelled">Cancelado</SelectItem>
              </SelectContent>
            </Select>

            <Select value={periodFilter} onValueChange={setPeriodFilter}>
              <SelectTrigger className="w-28">
                <SelectValue placeholder="Periodo" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos</SelectItem>
                <SelectItem value="2026">2026</SelectItem>
                <SelectItem value="2025">2025</SelectItem>
                <SelectItem value="2024">2024</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        <AssessmentWizard
          trigger={
            <Button className="flex items-center space-x-2">
              <Plus className="h-4 w-4" />
              <span>Nueva Evaluacion</span>
            </Button>
          }
        />
      </div>

      {/* Table */}
      <Card className="portal-surface border-0">
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <FileText className="h-5 w-5" />
            <span>Evaluaciones</span>
          </CardTitle>
          <CardDescription>
            Gestiona y da seguimiento a las evaluaciones NOM-035 de tu organizacion
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <SortableHeader
                    field="staff_name"
                    activeField={sortField}
                    direction={sortDirection}
                    onSort={handleSort}
                  >
                    Empleado
                  </SortableHeader>
                  <SortableHeader
                    field="period"
                    activeField={sortField}
                    direction={sortDirection}
                    onSort={handleSort}
                  >
                    Periodo
                  </SortableHeader>
                  <SortableHeader
                    field="status"
                    activeField={sortField}
                    direction={sortDirection}
                    onSort={handleSort}
                  >
                    Estado
                  </SortableHeader>
                  <SortableHeader
                    field="risk_level"
                    activeField={sortField}
                    direction={sortDirection}
                    onSort={handleSort}
                  >
                    Nivel de Riesgo
                  </SortableHeader>
                  <SortableHeader
                    field="created_at"
                    activeField={sortField}
                    direction={sortDirection}
                    onSort={handleSort}
                  >
                    Creado
                  </SortableHeader>
                  <TableHead>Acciones</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  Array.from({ length: 5 }).map((_, index) => (
                    <TableRow key={index}>
                      <TableCell>
                        <div className="h-4 bg-muted rounded animate-pulse w-32" />
                      </TableCell>
                      <TableCell>
                        <div className="h-4 bg-muted rounded animate-pulse w-16" />
                      </TableCell>
                      <TableCell>
                        <div className="h-5 bg-muted rounded animate-pulse w-20" />
                      </TableCell>
                      <TableCell>
                        <div className="h-5 bg-muted rounded animate-pulse w-16" />
                      </TableCell>
                      <TableCell>
                        <div className="h-4 bg-muted rounded animate-pulse w-20" />
                      </TableCell>
                      <TableCell>
                        <div className="flex gap-2">
                          <div className="h-8 bg-muted rounded animate-pulse w-16" />
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                ) : assessmentData?.data && assessmentData.data.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center py-8">
                      <div className="flex flex-col items-center gap-2">
                        <FileText className="h-8 w-8 text-muted-foreground/70" />
                        <p className="text-muted-foreground">No hay evaluaciones encontradas</p>
                        <p className="text-sm text-muted-foreground/80">
                          Crea tu primera evaluacion para comenzar
                        </p>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  assessmentData?.data.map((assessment) => (
                    <TableRow key={assessment.id}>
                      <TableCell className="font-medium">
                        <div>
                          <p>{getStaffName(assessment)}</p>
                          {getStaffEmail(assessment) && (
                            <p className="text-xs text-muted-foreground">{getStaffEmail(assessment)}</p>
                          )}
                        </div>
                      </TableCell>
                      <TableCell>{assessment.period}</TableCell>
                      <TableCell>{getStatusBadge(assessment.status)}</TableCell>
                      <TableCell>{getRiskBadge(assessment.risk_level)}</TableCell>
                      <TableCell className="text-sm text-muted-foreground">
                        {formatDate(assessment.created_at)}
                      </TableCell>
                      <TableCell>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon" className="h-8 w-8">
                              <MoreHorizontal className="h-4 w-4" />
                              <span className="sr-only">Acciones</span>
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            {assessment.status === "pending" && (
                              <>
                                <DropdownMenuItem onClick={() => onGenerateLink?.(assessment.id)}>
                                  <Link className="h-4 w-4 mr-2" />
                                  Generar Enlace
                                </DropdownMenuItem>
                                {getStaffEmail(assessment) && onSendEmail && (
                                  <DropdownMenuItem 
                                    onClick={() => onSendEmail(assessment.id)}
                                  >
                                    <Mail className="h-4 w-4 mr-2" />
                                    Enviar por Correo
                                  </DropdownMenuItem>
                                )}
                                <DropdownMenuSeparator />
                              </>
                            )}
                            {assessment.status === "completed" && (
                              <DropdownMenuItem onClick={() => onViewReport?.(assessment.id)}>
                                <Eye className="h-4 w-4 mr-2" />
                                Ver Reporte
                              </DropdownMenuItem>
                            )}
                            <DropdownMenuItem
                              onClick={() => setDeletingAssessment(assessment)}
                              className="text-destructive focus:text-destructive"
                            >
                              <Trash2 className="h-4 w-4 mr-2" />
                              Eliminar
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
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
              <div className="flex items-center space-x-2 text-sm text-muted-foreground">
                <span>Mostrar</span>
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
                <span>por pagina</span>
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

      {/* Delete Dialog */}
      <AssessmentDeleteDialog
        assessment={deletingAssessment}
        open={!!deletingAssessment}
        onOpenChange={(open) => !open && setDeletingAssessment(null)}
        onAssessmentDeleted={() => {
          setDeletingAssessment(null);
          onCreateAssessment(); // Refresh the list
        }}
      />
    </div>
  );
}

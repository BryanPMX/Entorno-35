"use client";

import React, { useState, useCallback, ReactNode } from "react";
import { ChevronUp, ChevronDown, Upload, Users, Edit2, Trash2, MoreHorizontal } from "lucide-react";
import { Button } from "@/components/ui/button";
import { StaffCreateDialog } from "./staff-create-dialog";
import { StaffEditDialog } from "./staff-edit-dialog";
import { StaffDeleteDialog } from "./staff-delete-dialog";
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
import type { PaginatedResponse, Staff } from "@/types/backend";

// Helper function to extract string value from sql.NullString-like objects
const getEmployeeIdValue = (
  employeeId: string | null | { Valid: boolean; String: string } | undefined
): string | null => {
  if (!employeeId) return null;
  if (typeof employeeId === "object") {
    return employeeId.Valid ? employeeId.String : null;
  }
  return employeeId;
};

interface StaffTableProps {
  staffData: PaginatedResponse<Staff> | null;
  isLoading: boolean;
  onPageChange: (page: number) => void;
  onLimitChange: (limit: number) => void;
  onCsvUpload: () => void;
  onDataRefresh?: () => void;
  currentPage: number;
  limit: number;
}

type SortField = "full_name" | "email" | "created_at";
type SortDirection = "asc" | "desc";

interface SortButtonProps {
  field: SortField;
  children: ReactNode;
  activeField: SortField;
  direction: SortDirection;
  onSort: (field: SortField) => void;
}

function SortButton({ field, children, activeField, direction, onSort }: SortButtonProps) {
  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={() => onSort(field)}
      className="h-auto p-0 font-medium hover:bg-transparent"
    >
      {children}
      {activeField === field &&
        (direction === "asc" ? (
          <ChevronUp className="ml-1 h-4 w-4" />
        ) : (
          <ChevronDown className="ml-1 h-4 w-4" />
        ))}
    </Button>
  );
}

export function StaffTable({
  staffData,
  isLoading,
  onPageChange,
  onLimitChange,
  onCsvUpload,
  onDataRefresh,
  currentPage,
  limit,
}: StaffTableProps) {
  const [sortField, setSortField] = useState<SortField>("created_at");
  const [sortDirection, setSortDirection] = useState<SortDirection>("desc");
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingStaff, setEditingStaff] = useState<Staff | null>(null);
  const [deletingStaff, setDeletingStaff] = useState<Staff | null>(null);

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

  const totalPages = staffData ? Math.ceil(staffData.total / limit) : 0;
  const startItem = staffData ? (currentPage - 1) * limit + 1 : 0;
  const endItem = staffData ? Math.min(currentPage * limit, staffData.total) : 0;

  return (
    <div className="space-y-6">
      {/* Header with stats and buttons */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-2">
            <Users className="h-5 w-5 text-muted-foreground" />
            <div>
              <p className="text-sm font-medium">
                {staffData?.total || 0} Miembros del Personal
              </p>
              <p className="text-xs text-muted-foreground">
                Mostrando {startItem}-{endItem} de {staffData?.total || 0}
              </p>
            </div>
          </div>
        </div>
        <div className="flex space-x-2">
          <Button onClick={() => setIsCreateDialogOpen(true)} variant="default">
            <Users className="h-4 w-4 mr-2" />
            Agregar Personal
          </Button>
          <Button onClick={onCsvUpload} variant="outline" className="flex items-center space-x-2">
            <Upload className="h-4 w-4" />
            <span>Importar CSV</span>
          </Button>
        </div>
      </div>

      {/* Staff Table */}
      <Card>
        <CardHeader>
          <CardTitle>Lista de Personal</CardTitle>
          <CardDescription>
            Gestiona la informacion de los miembros del personal
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>
                    <SortButton
                      field="full_name"
                      activeField={sortField}
                      direction={sortDirection}
                      onSort={handleSort}
                    >
                      Nombre
                    </SortButton>
                  </TableHead>
                  <TableHead>Identificador</TableHead>
                  <TableHead>
                    <SortButton
                      field="email"
                      activeField={sortField}
                      direction={sortDirection}
                      onSort={handleSort}
                    >
                      Correo Electronico
                    </SortButton>
                  </TableHead>
                  <TableHead>Departamento</TableHead>
                  <TableHead>Puesto</TableHead>
                  <TableHead>
                    <SortButton
                      field="created_at"
                      activeField={sortField}
                      direction={sortDirection}
                      onSort={handleSort}
                    >
                      Fecha de Registro
                    </SortButton>
                  </TableHead>
                  <TableHead className="w-[70px]">Acciones</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center py-8">
                      <div className="flex items-center justify-center space-x-2">
                        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-primary"></div>
                        <span>Cargando personal...</span>
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
                        {staff.curp || getEmployeeIdValue(staff.employee_id) || "—"}
                      </TableCell>
                      <TableCell>{staff.email || "—"}</TableCell>
                      <TableCell>
                        {staff.demographics?.department || "—"}
                      </TableCell>
                      <TableCell>
                        {staff.demographics?.role || "—"}
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {new Date(staff.created_at).toLocaleDateString("es-MX")}
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
                            <DropdownMenuItem onClick={() => setEditingStaff(staff)}>
                              <Edit2 className="h-4 w-4 mr-2" />
                              Editar
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              onClick={() => setDeletingStaff(staff)}
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
                ) : (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center py-8 text-muted-foreground">
                      No hay personal registrado. Importa datos usando la funcion de carga CSV.
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
                <span className="text-sm text-muted-foreground">Mostrar</span>
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
                <span className="text-sm text-muted-foreground">por pagina</span>
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

      {/* Dialogs */}
      <StaffCreateDialog
        open={isCreateDialogOpen}
        onOpenChange={setIsCreateDialogOpen}
        onStaffCreated={() => {
          setIsCreateDialogOpen(false);
          onDataRefresh?.();
        }}
      />

      <StaffEditDialog
        staff={editingStaff}
        open={!!editingStaff}
        onOpenChange={(open) => !open && setEditingStaff(null)}
        onStaffUpdated={() => {
          setEditingStaff(null);
          onDataRefresh?.();
        }}
      />

      <StaffDeleteDialog
        staff={deletingStaff}
        open={!!deletingStaff}
        onOpenChange={(open) => !open && setDeletingStaff(null)}
        onStaffDeleted={() => {
          setDeletingStaff(null);
          onDataRefresh?.();
        }}
      />
    </div>
  );
}

"use client";

import { useState } from "react";
import { AlertTriangle } from "lucide-react";
import { toast } from "sonner";
import { useQueryClient } from "@tanstack/react-query";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { staffService } from "@/services/staff.service";
import type { Staff } from "@/types/backend";

interface StaffDeleteDialogProps {
  staff: Staff | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onStaffDeleted?: () => void;
}

/**
 * StaffDeleteDialog - Confirmation dialog for deleting staff members
 */
export function StaffDeleteDialog({ staff, open, onOpenChange, onStaffDeleted }: StaffDeleteDialogProps) {
  const [isDeleting, setIsDeleting] = useState(false);
  const queryClient = useQueryClient();

  const handleDelete = async () => {
    if (!staff) return;

    setIsDeleting(true);

    try {
      await staffService.delete(staff.id);

      // Invalidate all related queries to refresh dashboard and lists
      queryClient.invalidateQueries({ queryKey: ["general-report"] });
      queryClient.invalidateQueries({ queryKey: ["staff"] });
      queryClient.invalidateQueries({ queryKey: ["assessments"] });

      toast.success("Personal eliminado exitosamente");
      onOpenChange(false);
      onStaffDeleted?.();
    } catch (error: unknown) {
      const response = typeof error === "object" && error && "response" in error
        ? (error as { response?: { data?: { error?: string }; status?: number } }).response
        : undefined;
      const errorData = response?.data;
      const statusCode = response?.status;

      // Handle specific errors
      if (statusCode === 404) {
        toast.error("El registro ya no existe");
        onOpenChange(false);
        onStaffDeleted?.();
      } else {
        toast.error(errorData?.error || "Error al eliminar el personal");
      }
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-red-100">
              <AlertTriangle className="h-5 w-5 text-red-600" />
            </div>
            <div>
              <DialogTitle>Eliminar Personal</DialogTitle>
              <DialogDescription>
                Esta accion no se puede deshacer.
              </DialogDescription>
            </div>
          </div>
        </DialogHeader>

        <div className="py-4">
          <p className="text-sm text-muted-foreground">
            ¿Estas seguro de que deseas eliminar a <span className="font-semibold text-foreground">{staff?.full_name}</span>?
          </p>
          {staff?.email && (
            <p className="text-sm text-muted-foreground mt-1">
              Correo: {staff.email}
            </p>
          )}
          <div className="mt-4 rounded-md bg-amber-50 border border-amber-200 p-3">
            <p className="text-sm text-amber-800">
              <strong>Advertencia:</strong> Esta accion eliminara permanentemente al empleado y todas sus evaluaciones asociadas. 
              La accion sera registrada en el log de auditoria para cumplimiento normativo.
            </p>
          </div>
        </div>

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isDeleting}
          >
            Cancelar
          </Button>
          <Button
            type="button"
            variant="destructive"
            onClick={handleDelete}
            disabled={isDeleting}
          >
            {isDeleting ? "Eliminando..." : "Eliminar"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

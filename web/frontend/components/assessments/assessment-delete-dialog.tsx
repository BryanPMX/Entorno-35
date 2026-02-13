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
import { assessmentService } from "@/services/assessment.service";
import type { Assessment } from "@/types/backend";

interface AssessmentDeleteDialogProps {
  assessment: Assessment | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onAssessmentDeleted?: () => void;
}

/**
 * AssessmentDeleteDialog - Dialogo de confirmacion para eliminar evaluaciones
 */
export function AssessmentDeleteDialog({ 
  assessment, 
  open, 
  onOpenChange, 
  onAssessmentDeleted 
}: AssessmentDeleteDialogProps) {
  const [isDeleting, setIsDeleting] = useState(false);
  const queryClient = useQueryClient();

  const getStaffName = (assessment: Assessment) => {
    if (assessment.staff) {
      return assessment.staff.full_name;
    }
    return `Empleado ${assessment.staff_id.slice(-4)}`;
  };

  const handleDelete = async () => {
    if (!assessment) return;

    setIsDeleting(true);

    try {
      await assessmentService.delete(assessment.id);

      // Invalidate all related queries to refresh dashboard and lists
      queryClient.invalidateQueries({ queryKey: ["general-report"] });
      queryClient.invalidateQueries({ queryKey: ["assessments"] });
      queryClient.invalidateQueries({ queryKey: ["individual-report"] });

      toast.success("Evaluacion eliminada exitosamente");
      onOpenChange(false);
      onAssessmentDeleted?.();
    } catch (error: unknown) {
      const response = typeof error === "object" && error && "response" in error
        ? (error as { response?: { data?: { error?: string }; status?: number } }).response
        : undefined;
      const errorMessage = response?.data?.error || "Error al eliminar la evaluacion";

      if (response?.status === 404) {
        toast.error("La evaluacion ya no existe");
        onOpenChange(false);
        onAssessmentDeleted?.();
      } else if (response?.status === 400) {
        toast.error(errorMessage);
      } else {
        toast.error(errorMessage);
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
              <DialogTitle>Eliminar Evaluacion</DialogTitle>
              <DialogDescription>
                Esta accion no se puede deshacer.
              </DialogDescription>
            </div>
          </div>
        </DialogHeader>

        <div className="py-4">
          <p className="text-sm text-muted-foreground">
            ¿Estas seguro de que deseas eliminar la evaluacion de{" "}
            <span className="font-semibold text-foreground">
              {assessment ? getStaffName(assessment) : ""}
            </span>
            ?
          </p>
          <div className="mt-4 rounded-md bg-amber-50 border border-amber-200 p-3">
            <p className="text-sm text-amber-800">
              <strong>Advertencia:</strong> Esta accion eliminara permanentemente la evaluacion. 
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

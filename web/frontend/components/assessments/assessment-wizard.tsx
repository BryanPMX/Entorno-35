"use client";

import React, { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { CheckCircle, Users, Settings, Eye, ArrowLeft, ArrowRight } from "lucide-react";
import { staffService } from "@/services/staff.service";
import { assessmentService } from "@/services/assessment.service";
import { useAuthStore } from "@/lib/store/auth-store";
import type { Staff } from "@/types/backend";

interface AssessmentWizardProps {
  trigger: React.ReactNode;
}

type WizardStep = "select-staff" | "configure" | "review";

interface SelectedStaff {
  id: string;
  full_name: string;
  department?: string;
  role?: string;
}

export function AssessmentWizard({ trigger }: AssessmentWizardProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [currentStep, setCurrentStep] = useState<WizardStep>("select-staff");
  const [selectedStaff, setSelectedStaff] = useState<SelectedStaff[]>([]);
  const [period, setPeriod] = useState(new Date().getFullYear());

  const { user } = useAuthStore();
  const queryClient = useQueryClient();

  // Fetch staff data
  const { data: staffData, isLoading: staffLoading } = useQuery({
    queryKey: ["staff", { limit: 100 }],
    queryFn: () => staffService.getAll({ limit: 100 }),
    enabled: isOpen,
  });

  // Create assessment mutation
  const createAssessmentMutation = useMutation({
    mutationFn: async (staffId: string) => {
      return assessmentService.create({
        staff_id: staffId,
        period,
      });
    },
    onSuccess: () => {
      // Invalidate all related queries to refresh dashboard and lists
      queryClient.invalidateQueries({ queryKey: ["assessments"] });
      queryClient.invalidateQueries({ queryKey: ["general-report"] });
      setCurrentStep("review");
    },
  });

  const steps = [
    { id: "select-staff", title: "Seleccionar", icon: Users },
    { id: "configure", title: "Configurar", icon: Settings },
    { id: "review", title: "Revisar", icon: Eye },
  ];

  const currentStepIndex = steps.findIndex(step => step.id === currentStep);

  const handleStaffToggle = (staff: Staff) => {
    setSelectedStaff(prev => {
      const isSelected = prev.some(s => s.id === staff.id);
      if (isSelected) {
        return prev.filter(s => s.id !== staff.id);
      } else {
        return [...prev, {
          id: staff.id,
          full_name: staff.full_name,
          department: staff.demographics?.department,
          role: staff.demographics?.role,
        }];
      }
    });
  };

  const handleNext = () => {
    if (currentStep === "select-staff" && selectedStaff.length === 0) return;
    if (currentStep === "configure" && !period) return;

    const nextIndex = currentStepIndex + 1;
    if (nextIndex < steps.length) {
      setCurrentStep(steps[nextIndex].id as WizardStep);
    }
  };

  const handleBack = () => {
    const prevIndex = currentStepIndex - 1;
    if (prevIndex >= 0) {
      setCurrentStep(steps[prevIndex].id as WizardStep);
    }
  };

  const handleCreateAssessments = async () => {
    const promises = selectedStaff.map(staff =>
      createAssessmentMutation.mutateAsync(staff.id)
    );

    try {
      await Promise.all(promises);
    } catch (error) {
      console.error("Error al crear evaluaciones:", error);
    }
  };

  const resetWizard = () => {
    setCurrentStep("select-staff");
    setSelectedStaff([]);
    setPeriod(new Date().getFullYear());
  };

  const handleClose = () => {
    setIsOpen(false);
    resetWizard();
  };

  const renderStepContent = () => {
    switch (currentStep) {
      case "select-staff":
        return (
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-medium">Seleccionar Personal</h3>
              <Badge variant="secondary">
                {selectedStaff.length} seleccionados
              </Badge>
            </div>

            <div className="max-h-96 overflow-y-auto space-y-2">
              {staffLoading ? (
                <div className="space-y-2">
                  {[...Array(5)].map((_, i) => (
                    <div key={i} className="flex items-center space-x-3 p-3 border rounded animate-pulse">
                      <div className="w-4 h-4 bg-muted rounded"></div>
                      <div className="flex-1">
                        <div className="h-4 bg-muted rounded w-32 mb-1"></div>
                        <div className="h-3 bg-muted rounded w-24"></div>
                      </div>
                    </div>
                  ))}
                </div>
              ) : staffData?.data && staffData.data.length > 0 ? (
                staffData.data.map((staff) => (
                  <div
                    key={staff.id}
                    className="flex items-center space-x-3 p-3 border rounded hover:bg-muted/50 transition-colors"
                  >
                    <Checkbox
                      id={staff.id}
                      checked={selectedStaff.some(s => s.id === staff.id)}
                      onCheckedChange={() => handleStaffToggle(staff)}
                    />
                    <div className="flex-1">
                      <Label htmlFor={staff.id} className="font-medium cursor-pointer">
                        {staff.full_name}
                      </Label>
                      <p className="text-sm text-muted-foreground">
                        {staff.demographics?.department && `${staff.demographics.department}`}
                        {staff.demographics?.role && ` • ${staff.demographics.role}`}
                      </p>
                    </div>
                  </div>
                ))
              ) : (
                <div className="text-center py-8 text-muted-foreground">
                  <Users className="h-8 w-8 mx-auto mb-2" />
                  <p>No hay personal registrado</p>
                  <p className="text-sm">Primero debes agregar personal a tu organizacion</p>
                </div>
              )}
            </div>
          </div>
        );

      case "configure":
        return (
          <div className="space-y-6">
            <div>
              <Label htmlFor="period">Periodo de Evaluacion</Label>
              <Input
                id="period"
                type="number"
                value={period}
                onChange={(e) => setPeriod(parseInt(e.target.value))}
                className="mt-1"
                min={2020}
                max={2030}
              />
            </div>

            <div className="p-4 bg-blue-50 border border-blue-200 rounded-md">
              <p className="text-sm text-blue-800">
                <strong>Nota:</strong> El tipo de guia (Guia II o Guia III) se determina automaticamente 
                segun el numero de empleados de tu empresa. La Guia II se usa para empresas con 
                16-50 empleados, y la Guia III para empresas con mas de 50 empleados.
              </p>
            </div>

            <div className="rounded-md border p-4">
              <h4 className="font-medium mb-2">Resumen de Seleccion</h4>
              <p className="text-sm text-muted-foreground">
                Se crearan {selectedStaff.length} evaluacion{selectedStaff.length !== 1 ? 'es' : ''} para el periodo {period}
              </p>
            </div>
          </div>
        );

      case "review":
        return (
          <div className="space-y-6">
            <div className="text-center">
              <CheckCircle className="h-12 w-12 text-green-500 mx-auto mb-4" />
              <h3 className="text-lg font-medium mb-2">¡Evaluaciones Creadas Exitosamente!</h3>
              <p className="text-muted-foreground">
                {selectedStaff.length} evaluacion{selectedStaff.length !== 1 ? 'es' : ''} {selectedStaff.length !== 1 ? 'han' : 'ha'} sido creada{selectedStaff.length !== 1 ? 's' : ''} para el periodo {period}.
              </p>
            </div>

            <div className="space-y-2">
              <h4 className="font-medium">Evaluaciones Creadas:</h4>
              <div className="max-h-40 overflow-y-auto space-y-1">
                {selectedStaff.map((staff) => (
                  <div key={staff.id} className="flex items-center justify-between p-2 bg-muted/50 rounded">
                    <span className="text-sm">{staff.full_name}</span>
                    <Badge variant="outline">Periodo {period}</Badge>
                  </div>
                ))}
              </div>
            </div>

            <div className="p-4 bg-green-50 border border-green-200 rounded-md">
              <p className="text-sm text-green-800">
                <strong>Siguiente paso:</strong> Genera enlaces de evaluacion desde la tabla de evaluaciones 
                y comparte con tu personal para que completen su cuestionario NOM-035.
              </p>
            </div>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>
        {trigger}
      </DialogTrigger>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <Users className="h-5 w-5" />
            <span>Crear Ciclo de Evaluacion</span>
          </DialogTitle>
          <DialogDescription>
            Crea evaluaciones para multiples empleados en unos pocos pasos.
          </DialogDescription>
        </DialogHeader>

        {/* Stepper */}
        <div className="flex items-center justify-center mb-8">
          {steps.map((step, index) => {
            const Icon = step.icon;
            const isActive = index === currentStepIndex;
            const isCompleted = index < currentStepIndex;

            return (
              <React.Fragment key={step.id}>
                <div className="flex flex-col items-center">
                  <div
                    className={`w-10 h-10 rounded-full flex items-center justify-center border-2 transition-colors ${
                      isCompleted
                        ? 'bg-green-500 border-green-500 text-white'
                        : isActive
                        ? 'border-primary text-primary'
                        : 'border-muted-foreground text-muted-foreground'
                    }`}
                  >
                    {isCompleted ? (
                      <CheckCircle className="h-5 w-5" />
                    ) : (
                      <Icon className="h-5 w-5" />
                    )}
                  </div>
                  <span className={`text-xs mt-2 ${isActive ? 'text-primary font-medium' : 'text-muted-foreground'}`}>
                    {step.title}
                  </span>
                </div>
                {index < steps.length - 1 && (
                  <div className={`w-16 h-0.5 mx-4 mt-5 ${index < currentStepIndex ? 'bg-green-500' : 'bg-muted-foreground'}`} />
                )}
              </React.Fragment>
            );
          })}
        </div>

        {/* Step Content */}
        <div className="min-h-[300px]">
          {renderStepContent()}
        </div>

        {/* Navigation */}
        {currentStep !== "review" && (
          <div className="flex items-center justify-between pt-6 border-t">
            <Button
              variant="outline"
              onClick={handleBack}
              disabled={currentStepIndex === 0}
              className="flex items-center space-x-2"
            >
              <ArrowLeft className="h-4 w-4" />
              <span>Atras</span>
            </Button>

            <div className="flex space-x-2">
              <Button variant="outline" onClick={handleClose}>
                Cancelar
              </Button>

              {currentStep === "configure" ? (
                <Button
                  onClick={handleCreateAssessments}
                  disabled={createAssessmentMutation.isPending}
                  className="flex items-center space-x-2"
                >
                  {createAssessmentMutation.isPending ? (
                    <>
                      <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                      <span>Creando...</span>
                    </>
                  ) : (
                    <>
                      <span>Crear Evaluaciones</span>
                      <ArrowRight className="h-4 w-4" />
                    </>
                  )}
                </Button>
              ) : (
                <Button
                  onClick={handleNext}
                  disabled={
                    (String(currentStep) === "select-staff" && selectedStaff.length === 0) ||
                    (String(currentStep) === "configure" && !period)
                  }
                  className="flex items-center space-x-2"
                >
                  <span>Siguiente</span>
                  <ArrowRight className="h-4 w-4" />
                </Button>
              )}
            </div>
          </div>
        )}

        {currentStep === "review" && (
          <div className="flex justify-end pt-6 border-t">
            <Button onClick={handleClose} className="flex items-center space-x-2">
              <CheckCircle className="h-4 w-4" />
              <span>Listo</span>
            </Button>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}

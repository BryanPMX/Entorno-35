"use client";

import React, { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Progress } from "@/components/ui/progress";
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
  const [guideType, setGuideType] = useState<"II" | "III">("II");

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
      queryClient.invalidateQueries({ queryKey: ["assessments"] });
      setCurrentStep("review");
    },
  });

  const steps = [
    { id: "select-staff", title: "Select Staff", icon: Users },
    { id: "configure", title: "Configure", icon: Settings },
    { id: "review", title: "Review & Create", icon: Eye },
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
    // Create assessments for all selected staff
    const promises = selectedStaff.map(staff =>
      createAssessmentMutation.mutateAsync(staff.id)
    );

    try {
      await Promise.all(promises);
      // Success is handled by the mutation's onSuccess
    } catch (error) {
      console.error("Failed to create assessments:", error);
    }
  };

  const resetWizard = () => {
    setCurrentStep("select-staff");
    setSelectedStaff([]);
    setPeriod(new Date().getFullYear());
    setGuideType("II");
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
              <h3 className="text-lg font-medium">Select Staff Members</h3>
              <Badge variant="secondary">
                {selectedStaff.length} selected
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
              ) : (
                staffData?.data.map((staff) => (
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
              )}
            </div>
          </div>
        );

      case "configure":
        return (
          <div className="space-y-6">
            <div>
              <Label htmlFor="period">Assessment Period</Label>
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

            <div>
              <Label>Guide Type</Label>
              <p className="text-sm text-muted-foreground mb-3">
                Automatically determined based on company size
              </p>
              <Select value={guideType} onValueChange={(value: "II" | "III") => setGuideType(value)}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="II">Guide II (16-50 employees)</SelectItem>
                  <SelectItem value="III">Guide III (50+ employees)</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        );

      case "review":
        return (
          <div className="space-y-6">
            <div className="text-center">
              <CheckCircle className="h-12 w-12 text-green-500 mx-auto mb-4" />
              <h3 className="text-lg font-medium mb-2">Assessments Created Successfully!</h3>
              <p className="text-muted-foreground">
                {selectedStaff.length} assessment{selectedStaff.length !== 1 ? 's' : ''} have been created for the {period} period.
              </p>
            </div>

            <div className="space-y-2">
              <h4 className="font-medium">Created Assessments:</h4>
              <div className="max-h-40 overflow-y-auto space-y-1">
                {selectedStaff.map((staff) => (
                  <div key={staff.id} className="flex items-center justify-between p-2 bg-muted/50 rounded">
                    <span className="text-sm">{staff.full_name}</span>
                    <Badge variant="outline">Period {period}</Badge>
                  </div>
                ))}
              </div>
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
            <span>Create Assessment Cycle</span>
          </DialogTitle>
          <DialogDescription>
            Create assessments for multiple staff members in a few simple steps.
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
              <span>Back</span>
            </Button>

            <div className="flex space-x-2">
              <Button variant="outline" onClick={handleClose}>
                Cancel
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
                      <span>Creating...</span>
                    </>
                  ) : (
                    <>
                      <span>Create Assessments</span>
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
                  <span>Next</span>
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
              <span>Done</span>
            </Button>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
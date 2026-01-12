"use client";

import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

import { staffService, type UpdateStaffRequest } from "@/services/staff.service";
import type { Staff } from "@/types/backend";

// Special value for "no selection" since Radix UI doesn't allow empty strings
const NONE_VALUE = "__none__";

// Form validation schema
const staffSchema = z.object({
  full_name: z.string().min(2, "El nombre debe tener al menos 2 caracteres"),
  email: z.string().email("Correo electronico invalido").optional().or(z.literal("")),
  department: z.string().optional(),
  role: z.string().optional(),
  gender: z.string().optional(),
  age_range: z.string().optional(),
  marital_status: z.string().optional(),
  education_level: z.string().optional(),
  shift_type: z.string().optional(),
  time_in_position: z.string().optional(),
  total_work_experience: z.string().optional(),
});

type StaffFormValues = z.infer<typeof staffSchema>;

interface StaffEditDialogProps {
  staff: Staff | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onStaffUpdated?: () => void;
}

// Helper to convert empty string to NONE_VALUE for Select
const toSelectValue = (value: string | undefined | null): string => {
  return value || NONE_VALUE;
};

// Helper to convert NONE_VALUE back to undefined for API
const fromSelectValue = (value: string | undefined): string | undefined => {
  return value === NONE_VALUE ? undefined : value;
};

/**
 * StaffEditDialog - Modal dialog for editing staff members
 */
export function StaffEditDialog({ staff, open, onOpenChange, onStaffUpdated }: StaffEditDialogProps) {
  const [isSubmitting, setIsSubmitting] = useState(false);

  const form = useForm<StaffFormValues>({
    resolver: zodResolver(staffSchema),
    defaultValues: {
      full_name: "",
      email: "",
      department: "",
      role: "",
      gender: NONE_VALUE,
      age_range: NONE_VALUE,
      marital_status: NONE_VALUE,
      education_level: NONE_VALUE,
      shift_type: NONE_VALUE,
      time_in_position: NONE_VALUE,
      total_work_experience: NONE_VALUE,
    },
  });

  // Populate form when staff changes
  useEffect(() => {
    if (staff) {
      form.reset({
        full_name: staff.full_name || "",
        email: staff.email || "",
        department: staff.demographics?.department || "",
        role: staff.demographics?.role || "",
        gender: toSelectValue(staff.demographics?.gender),
        age_range: toSelectValue(staff.demographics?.age_range),
        marital_status: toSelectValue(staff.demographics?.marital_status),
        education_level: toSelectValue(staff.demographics?.education_level),
        shift_type: toSelectValue(staff.demographics?.shift_type),
        time_in_position: toSelectValue(staff.demographics?.time_in_position),
        total_work_experience: toSelectValue(staff.demographics?.total_work_experience),
      });
    }
  }, [staff, form]);

  const onSubmit = async (values: StaffFormValues) => {
    if (!staff) return;

    setIsSubmitting(true);

    try {
      const updateData: UpdateStaffRequest = {
        full_name: values.full_name,
        email: values.email || undefined,
        demographics: {
          department: values.department || undefined,
          role: values.role || undefined,
          gender: fromSelectValue(values.gender),
          age_range: fromSelectValue(values.age_range),
          marital_status: fromSelectValue(values.marital_status),
          education_level: fromSelectValue(values.education_level),
          shift_type: fromSelectValue(values.shift_type),
          time_in_position: fromSelectValue(values.time_in_position),
          total_work_experience: fromSelectValue(values.total_work_experience),
        },
      };

      await staffService.update(staff.id, updateData);

      toast.success("Personal actualizado exitosamente");
      onOpenChange(false);
      onStaffUpdated?.();
    } catch (error: any) {
      const errorMessage = error?.response?.data?.error || "Error al actualizar el personal";
      
      if (error?.response?.status === 409) {
        toast.error("El registro fue modificado por otro usuario. Por favor recarga la pagina.");
      } else {
        toast.error(errorMessage);
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px] max-h-[85vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Editar Personal</DialogTitle>
          <DialogDescription>
            Modifica la informacion del miembro del personal.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            {/* Informacion Basica */}
            <div className="space-y-4">
              <h3 className="text-lg font-medium">Informacion Basica</h3>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <FormField
                  control={form.control}
                  name="full_name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Nombre Completo *</FormLabel>
                      <FormControl>
                        <Input placeholder="Juan Perez Garcia" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Correo Electronico</FormLabel>
                      <FormControl>
                        <Input type="email" placeholder="juan.perez@empresa.com" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              {/* CURP is read-only */}
              {staff?.curp && (
                <div className="rounded-md bg-muted px-3 py-2">
                  <span className="text-sm text-muted-foreground">CURP: </span>
                  <span className="text-sm font-mono">{staff.curp}</span>
                </div>
              )}
            </div>

            {/* Informacion Laboral */}
            <div className="space-y-4">
              <h3 className="text-lg font-medium">Informacion Laboral</h3>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <FormField
                  control={form.control}
                  name="department"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Departamento</FormLabel>
                      <FormControl>
                        <Input placeholder="TI, RH, Ventas..." {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="role"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Puesto</FormLabel>
                      <FormControl>
                        <Input placeholder="Desarrollador, Gerente..." {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            </div>

            {/* Datos Demograficos */}
            <div className="space-y-4">
              <h3 className="text-lg font-medium">Datos Demograficos</h3>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <FormField
                  control={form.control}
                  name="gender"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Genero</FormLabel>
                      <Select onValueChange={field.onChange} value={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Seleccionar genero" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value={NONE_VALUE}>Sin especificar</SelectItem>
                          <SelectItem value="masculino">Masculino</SelectItem>
                          <SelectItem value="femenino">Femenino</SelectItem>
                          <SelectItem value="otro">Otro</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="age_range"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Rango de Edad</FormLabel>
                      <Select onValueChange={field.onChange} value={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Seleccionar rango" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value={NONE_VALUE}>Sin especificar</SelectItem>
                          <SelectItem value="18-25">18-25 años</SelectItem>
                          <SelectItem value="26-35">26-35 años</SelectItem>
                          <SelectItem value="36-45">36-45 años</SelectItem>
                          <SelectItem value="46-55">46-55 años</SelectItem>
                          <SelectItem value="56+">56+ años</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="marital_status"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Estado Civil</FormLabel>
                      <Select onValueChange={field.onChange} value={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Seleccionar estado" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value={NONE_VALUE}>Sin especificar</SelectItem>
                          <SelectItem value="soltero">Soltero(a)</SelectItem>
                          <SelectItem value="casado">Casado(a)</SelectItem>
                          <SelectItem value="divorciado">Divorciado(a)</SelectItem>
                          <SelectItem value="viudo">Viudo(a)</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="education_level"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Nivel Educativo</FormLabel>
                      <Select onValueChange={field.onChange} value={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Seleccionar nivel" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value={NONE_VALUE}>Sin especificar</SelectItem>
                          <SelectItem value="secundaria">Secundaria</SelectItem>
                          <SelectItem value="preparatoria">Preparatoria</SelectItem>
                          <SelectItem value="universidad">Universidad</SelectItem>
                          <SelectItem value="postgrado">Postgrado</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="shift_type"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Tipo de Turno</FormLabel>
                      <Select onValueChange={field.onChange} value={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Seleccionar turno" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value={NONE_VALUE}>Sin especificar</SelectItem>
                          <SelectItem value="diurno">Diurno</SelectItem>
                          <SelectItem value="nocturno">Nocturno</SelectItem>
                          <SelectItem value="mixto">Mixto</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="time_in_position"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Tiempo en el Puesto</FormLabel>
                      <Select onValueChange={field.onChange} value={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Seleccionar tiempo" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value={NONE_VALUE}>Sin especificar</SelectItem>
                          <SelectItem value="<1 año">Menos de 1 año</SelectItem>
                          <SelectItem value="1-3 años">1-3 años</SelectItem>
                          <SelectItem value="3-5 años">3-5 años</SelectItem>
                          <SelectItem value="5+ años">Mas de 5 años</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="total_work_experience"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Experiencia Laboral Total</FormLabel>
                      <Select onValueChange={field.onChange} value={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Seleccionar experiencia" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value={NONE_VALUE}>Sin especificar</SelectItem>
                          <SelectItem value="<1 año">Menos de 1 año</SelectItem>
                          <SelectItem value="1-5 años">1-5 años</SelectItem>
                          <SelectItem value="5-10 años">5-10 años</SelectItem>
                          <SelectItem value="10+ años">Mas de 10 años</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={isSubmitting}
              >
                Cancelar
              </Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? "Guardando..." : "Guardar Cambios"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

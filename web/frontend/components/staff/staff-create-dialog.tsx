"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Plus } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
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

import { staffService, type CreateStaffRequest } from "@/services/staff.service";

// Form validation schema
const staffSchema = z.object({
  full_name: z.string().min(2, "El nombre debe tener al menos 2 caracteres"),
  email: z.string().email("Correo electronico invalido").optional().or(z.literal("")),
  curp: z.string().length(18, "El CURP debe tener exactamente 18 caracteres").optional().or(z.literal("")),
  department: z.string().optional(),
  role: z.string().optional(),
  gender: z.enum(["masculino", "femenino", "otro"]).optional(),
  age_range: z.enum(["18-25", "26-35", "36-45", "46-55", "56+"]).optional(),
  marital_status: z.enum(["soltero", "casado", "divorciado", "viudo"]).optional(),
  education_level: z.enum(["secundaria", "preparatoria", "universidad", "postgrado"]).optional(),
  shift_type: z.enum(["diurno", "nocturno", "mixto"]).optional(),
  time_in_position: z.enum(["<1 año", "1-3 años", "3-5 años", "5+ años"]).optional(),
  total_work_experience: z.enum(["<1 año", "1-5 años", "5-10 años", "10+ años"]).optional(),
});

type StaffFormValues = z.infer<typeof staffSchema>;

interface StaffCreateDialogProps {
  onStaffCreated?: () => void;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}

/**
 * StaffCreateDialog - Modal para crear miembros del personal manualmente
 */
export function StaffCreateDialog({ onStaffCreated, open: controlledOpen, onOpenChange }: StaffCreateDialogProps) {
  const [internalOpen, setInternalOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const open = controlledOpen !== undefined ? controlledOpen : internalOpen;
  const setOpen = onOpenChange || setInternalOpen;

  const form = useForm<StaffFormValues>({
    resolver: zodResolver(staffSchema),
    defaultValues: {
      full_name: "",
      email: "",
      curp: "",
      department: "",
      role: "",
      gender: undefined,
      age_range: undefined,
      marital_status: undefined,
      education_level: undefined,
      shift_type: undefined,
      time_in_position: undefined,
      total_work_experience: undefined,
    },
  });

  const onSubmit = async (values: StaffFormValues) => {
    setIsSubmitting(true);

    try {
      const staffData: CreateStaffRequest = {
        full_name: values.full_name,
        email: values.email || undefined,
        curp: values.curp || undefined,
        demographics: {
          department: values.department || undefined,
          role: values.role || undefined,
          gender: values.gender,
          age_range: values.age_range,
          marital_status: values.marital_status,
          education_level: values.education_level,
          shift_type: values.shift_type,
          time_in_position: values.time_in_position,
          total_work_experience: values.total_work_experience,
        },
      };

      await staffService.create(staffData);

      toast.success("Personal creado exitosamente");
      form.reset();
      setOpen(false);
      onStaffCreated?.();
    } catch (error: any) {
      const errorMessage = error?.response?.data?.error || "Error al crear el personal";
      toast.error(errorMessage);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      {controlledOpen === undefined && (
        <DialogTrigger asChild>
          <Button>
            <Plus className="h-4 w-4 mr-2" />
            Agregar Personal
          </Button>
        </DialogTrigger>
      )}

      <DialogContent className="sm:max-w-[600px] max-h-[85vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Agregar Nuevo Personal</DialogTitle>
          <DialogDescription>
            Crea un miembro del personal con su informacion demografica.
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

              <FormField
                control={form.control}
                name="curp"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>CURP</FormLabel>
                    <FormControl>
                      <Input placeholder="PEGJ900101HDFRRR01" maxLength={18} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
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
                      <Select onValueChange={field.onChange} defaultValue={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Seleccionar genero" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
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
                      <Select onValueChange={field.onChange} defaultValue={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Seleccionar rango" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
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
                      <Select onValueChange={field.onChange} defaultValue={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Seleccionar estado" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
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
                      <Select onValueChange={field.onChange} defaultValue={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Seleccionar nivel" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
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
                      <Select onValueChange={field.onChange} defaultValue={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Seleccionar turno" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
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
                      <Select onValueChange={field.onChange} defaultValue={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Seleccionar tiempo" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
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
                      <Select onValueChange={field.onChange} defaultValue={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Seleccionar experiencia" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
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
                onClick={() => setOpen(false)}
                disabled={isSubmitting}
              >
                Cancelar
              </Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? "Creando..." : "Crear Personal"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

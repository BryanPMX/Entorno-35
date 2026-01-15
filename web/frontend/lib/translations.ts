/**
 * Spanish Translations for Entorno35 UI
 * 
 * This file contains all Spanish translations for the application.
 * The system uses Spanish as the primary language for NOM-035 compliance.
 */

export const translations = {
  // Common
  common: {
    loading: "Cargando...",
    error: "Error",
    success: "Exito",
    cancel: "Cancelar",
    save: "Guardar",
    delete: "Eliminar",
    edit: "Editar",
    create: "Crear",
    search: "Buscar",
    filter: "Filtrar",
    export: "Exportar",
    import: "Importar",
    download: "Descargar",
    upload: "Subir",
    back: "Volver",
    next: "Siguiente",
    previous: "Anterior",
    submit: "Enviar",
    confirm: "Confirmar",
    close: "Cerrar",
    yes: "Si",
    no: "No",
    all: "Todos",
    none: "Ninguno",
    actions: "Acciones",
    status: "Estado",
    date: "Fecha",
    name: "Nombre",
    email: "Correo electronico",
    password: "Contrasena",
    noData: "No hay datos disponibles",
  },

  // Navigation
  nav: {
    dashboard: "Panel de Control",
    staff: "Personal",
    assessments: "Evaluaciones",
    reports: "Reportes",
    settings: "Configuracion",
    logout: "Cerrar Sesion",
    login: "Iniciar Sesion",
    logoutSuccess: "Sesion cerrada exitosamente",
  },

  // Dashboard
  dashboard: {
    title: "Panel de Control",
    subtitle: "Resumen de Cumplimiento NOM-035",
    welcome: "Bienvenido a Entorno 35 - Plataforma de Cumplimiento NOM-035",
    totalStaff: "Personal Total",
    completedAssessments: "Evaluaciones Completadas",
    participationRate: "Tasa de Participacion",
    riskDistribution: "Distribucion de Riesgo",
    highRiskCases: "casos de alto riesgo",
    fromLastPeriod: "del periodo anterior",
    getStarted: "Comienza con las Evaluaciones NOM-035",
    getStartedDesc: "Crea tu primer ciclo de evaluacion y comienza a monitorear factores de riesgo psicosocial.",
    createFirstAssessment: "Crear Primera Evaluacion",
    createAssessment: "Crear Evaluacion",
    quickActions: "Acciones Rapidas",
    quickActionsDesc: "Tareas comunes para gestionar tu programa de cumplimiento NOM-035",
    manageStaff: "Gestionar Personal",
    viewReports: "Ver Reportes",
    newAssessmentCycle: "Nuevo Ciclo de Evaluacion",
    demographicAnalysis: "Analisis Demografico",
    demographicDesc: "Distribucion del personal por caracteristicas demograficas",
    age: "Edad",
    ageDesc: "Distribucion por rango de edad",
    maritalStatus: "Estado Civil",
    maritalStatusDesc: "Distribucion por estado civil",
    shiftType: "Turno",
    shiftTypeDesc: "Distribucion por tipo de turno",
    experience: "Experiencia",
    experienceDesc: "Distribucion por experiencia laboral",
    demographicDistribution: "Distribucion Demografica",
    riskByDemographics: "Riesgo Demografico",
    riskByAge: "Riesgo por Rango de Edad",
    riskByAgeDesc: "Distribucion de niveles de riesgo segun la edad del personal",
    riskByShift: "Riesgo por Tipo de Turno",
    riskByShiftDesc: "Distribucion de niveles de riesgo segun el turno laboral",
  },

  // Staff Management
  staff: {
    title: "Gestion de Personal",
    subtitle: "Administra los empleados de tu organizacion",
    addStaff: "Agregar Personal",
    importCSV: "Importar CSV",
    downloadTemplate: "Descargar Plantilla",
    fullName: "Nombre Completo",
    email: "Correo Electronico",
    curp: "CURP",
    department: "Departamento",
    position: "Puesto",
    createdAt: "Fecha de Registro",
    noStaff: "No hay personal registrado",
    importSuccess: "Personal importado exitosamente",
    importError: "Error al importar personal",
  },

  // Assessments
  assessments: {
    title: "Evaluaciones",
    subtitle: "Gestiona las evaluaciones NOM-035",
    createNew: "Nueva Evaluacion",
    period: "Periodo",
    guideType: "Tipo de Guia",
    status: "Estado",
    staffName: "Nombre del Empleado",
    completedAt: "Completado",
    viewReport: "Ver Reporte",
    sendReminder: "Enviar Recordatorio",
    pending: "Pendiente",
    inProgress: "En Progreso",
    completed: "Completado",
    expired: "Expirado",
    noAssessments: "No hay evaluaciones",
    guideAutomatic: "El tipo de guia se determina automaticamente segun el numero de empleados",
  },

  // Reports
  reports: {
    title: "Reporte de Evaluacion",
    subtitle: "Resultados de la Evaluacion de Riesgo Psicosocial NOM-035",
    staffInfo: "Informacion del Empleado",
    riskAssessment: "Evaluacion de Riesgo",
    totalScore: "Puntuacion Total",
    riskLevel: "Nivel de Riesgo",
    categoryAnalysis: "Analisis por Categoria",
    domainAnalysis: "Analisis por Dominio",
    recommendations: "Recomendaciones",
    downloadPDF: "Descargar PDF",
    refresh: "Actualizar",
    backToAssessments: "Volver a Evaluaciones",
    assessmentTimeline: "Linea de Tiempo",
    assessmentPeriod: "Periodo de Evaluacion",
    importantNote: "Importante: Este reporte refleja los resultados de la evaluacion de riesgo psicosocial solo para el periodo {period}.",
    multipleAssessmentsNote: "Si este empleado tiene multiples evaluaciones, cada reporte es independiente y especifico para su periodo de evaluacion.",
    medicalAttention: "Esta evaluacion indica necesidad de atencion medica. Por favor consulte con servicios de salud ocupacional.",
  },

  // Risk Levels
  riskLevels: {
    nulo: "Nulo",
    bajo: "Bajo",
    medio: "Medio",
    alto: "Alto",
    muy_alto: "Muy Alto",
    nuloDesc: "Sin factores de riesgo psicosocial significativos detectados",
    bajoDesc: "Factores de riesgo psicosocial minimos presentes",
    medioDesc: "Factores de riesgo psicosocial moderados que requieren atencion",
    altoDesc: "Factores de riesgo psicosocial altos que requieren intervencion",
    muyAltoDesc: "Factores de riesgo psicosocial severos que requieren accion inmediata",
  },

  // Login
  login: {
    title: "Iniciar Sesion",
    subtitle: "Accede a tu cuenta de Entorno 35",
    companyLogin: "Acceso Empresa",
    staffLogin: "Acceso Personal",
    companyId: "ID de Empresa",
    employeeId: "ID de Empleado",
    password: "Contrasena",
    rememberMe: "Recordarme",
    forgotPassword: "Olvide mi contrasena",
    loginButton: "Iniciar Sesion",
    loginError: "Credenciales invalidas",
    loginSuccess: "Sesion iniciada exitosamente",
  },

  // Validation Messages
  validation: {
    required: "Este campo es obligatorio",
    invalidEmail: "Correo electronico invalido",
    invalidCURP: "CURP invalido",
    minLength: "Minimo {min} caracteres",
    maxLength: "Maximo {max} caracteres",
    passwordMismatch: "Las contrasenas no coinciden",
  },

  // Charts
  charts: {
    riskDistribution: "Distribucion de Riesgo",
    riskDistributionDesc: "Distribucion de evaluaciones por nivel de riesgo",
    departmentHeatmap: "Riesgo Por Departamento",
    departmentHeatmapDesc: "Riesgo promedio por departamento",
    noDataAvailable: "No hay datos disponibles",
    riskScore: "Puntuacion de Riesgo",
    totalAssessments: "Evaluaciones Totales",
  },

  // Onboarding Cards
  onboarding: {
    staffManagement: "Gestion de Personal",
    staffManagementDesc: "Importa y gestiona a tus empleados",
    staffManagementHint: "Comienza importando tu lista de personal mediante carga CSV para una gestion masiva sencilla.",
    assessmentCreation: "Creacion de Evaluaciones",
    assessmentCreationDesc: "Genera enlaces de evaluacion seguros",
    assessmentCreationHint: "Crea ciclos de evaluacion y genera automaticamente enlaces seguros para la participacion del personal.",
    complianceReports: "Reportes de Cumplimiento",
    complianceReportsDesc: "Genera reportes de cumplimiento NOM-035",
    complianceReportsHint: "Accede a evaluaciones individuales y reportes de cumplimiento a nivel empresa con recomendaciones accionables.",
  },
} as const;

export type TranslationKey = keyof typeof translations;
export type Translations = typeof translations;

/**
 * Get a translation value by key path
 * @param path - Dot-separated path to translation (e.g., "dashboard.title")
 * @returns The translated string
 */
export function t(path: string): string {
  const keys = path.split(".");
  let value: unknown = translations;
  
  for (const key of keys) {
    if (value && typeof value === "object" && key in value) {
      value = (value as Record<string, unknown>)[key];
    } else {
      console.warn(`Translation not found: ${path}`);
      return path;
    }
  }
  
  return typeof value === "string" ? value : path;
}

/**
 * Get risk level translation
 */
export function getRiskLevelLabel(riskLevel: string): string {
  const levels: Record<string, string> = {
    nulo: translations.riskLevels.nulo,
    bajo: translations.riskLevels.bajo,
    medio: translations.riskLevels.medio,
    alto: translations.riskLevels.alto,
    muy_alto: translations.riskLevels.muy_alto,
  };
  return levels[riskLevel] || riskLevel;
}

/**
 * Get risk level description
 */
export function getRiskLevelDescription(riskLevel: string): string {
  const descriptions: Record<string, string> = {
    nulo: translations.riskLevels.nuloDesc,
    bajo: translations.riskLevels.bajoDesc,
    medio: translations.riskLevels.medioDesc,
    alto: translations.riskLevels.altoDesc,
    muy_alto: translations.riskLevels.muyAltoDesc,
  };
  return descriptions[riskLevel] || "Evaluacion de nivel de riesgo";
}

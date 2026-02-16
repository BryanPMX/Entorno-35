"use client";

import React, { useState, useEffect, useCallback } from "react";
import { useParams } from "next/navigation";
import { useQuery, useMutation } from "@tanstack/react-query";
import { motion, AnimatePresence } from "framer-motion";
import axios from "axios";
import { ChevronLeft, ChevronRight, CheckCircle, AlertCircle, Loader2, Cloud } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { assessmentService } from "@/services/assessment.service";
import { toast } from "sonner";

function resolveAssessmentErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const apiError = (error.response?.data as { error?: string } | undefined)?.error;
    if (typeof apiError === "string" && apiError.trim().length > 0) {
      return apiError;
    }

    if (typeof error.message === "string" && error.message.trim().length > 0) {
      return error.message;
    }
  }

  if (error instanceof Error && error.message.trim().length > 0) {
    return error.message;
  }

  return "No se pudo cargar la evaluacion. El enlace puede ser invalido o haber expirado.";
}

const LIKERT_OPTIONS = [
  { value: 0, label: "Siempre", description: "Todo el tiempo", key: "1" },
  { value: 1, label: "Casi siempre", description: "Casi todo el tiempo", key: "2" },
  { value: 2, label: "Algunas veces", description: "De forma ocasional", key: "3" },
  { value: 3, label: "Casi nunca", description: "Rara vez", key: "4" },
  { value: 4, label: "Nunca", description: "En ningun momento", key: "5" },
] as const;

/**
 * Public Assessment Exam Page
 *
 * Distraction-free interface for staff to complete NOM-035 assessments.
 * Features progress tracking, question navigation, and seamless submission.
 */
export default function AssessmentExamPage() {
  const params = useParams();
  const tokenParam = params.token;
  const token = Array.isArray(tokenParam) ? tokenParam[0] : tokenParam;

  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [responses, setResponses] = useState<Record<number, number>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitSuccess, setSubmitSuccess] = useState(false);
  const [selectedOption, setSelectedOption] = useState<number | null>(null);
  const [pressedKey, setPressedKey] = useState<string | null>(null);
  const [saveStatus, setSaveStatus] = useState<'saved' | 'saving' | 'error'>('saved');

  // Fetch assessment data and questions
  const {
    data: assessmentData,
    isLoading: isLoadingAssessment,
    error: assessmentError,
  } = useQuery({
    queryKey: ["public-assessment", token],
    queryFn: () => assessmentService.fetchByToken(token || ""),
    enabled: Boolean(token),
  });

  // Submit responses mutation
  const submitMutation = useMutation({
    mutationFn: (responses: { question_id: number; value: number }[]) => {
      if (!token) {
        throw new Error("Enlace de evaluacion invalido. Falta el token.");
      }
      return assessmentService.submitResponses(token, responses);
    },
    onSuccess: () => {
      setSubmitSuccess(true);
      setIsSubmitting(false);
    },
    onError: () => {
      setIsSubmitting(false);
      // Error handling is done in the UI
    },
  });

  const assessment = assessmentData?.assessment;
  const questions = assessmentData?.questions || [];
  const currentQuestion = questions[currentQuestionIndex];
  const progress = questions.length > 0 ? ((currentQuestionIndex + 1) / questions.length) * 100 : 0;

  const handlePrevious = useCallback(() => {
    setCurrentQuestionIndex(prev => Math.max(0, prev - 1));
  }, []);

  const handleNext = useCallback(() => {
    setCurrentQuestionIndex(prev => Math.min(questions.length - 1, prev + 1));
  }, [questions.length]);

  const handleAnswerSelect = useCallback((value: number) => {
    if (!currentQuestion) return;

    setResponses(prev => ({
      ...prev,
      [currentQuestion.id]: value,
    }));

    setSaveStatus('saving');
    setTimeout(() => setSaveStatus('saved'), 500);

    setTimeout(() => {
      setCurrentQuestionIndex(prev => {
        if (prev < questions.length - 1) {
          return prev + 1;
        }
        return prev;
      });
    }, 300);
  }, [currentQuestion, questions.length]);

  // Keyboard navigation
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (isSubmitting || submitSuccess) return;

      const key = event.key;
      setPressedKey(key);

      // Clear pressed key after animation
      setTimeout(() => setPressedKey(null), 100);

      // Arrow key navigation
      if (key === 'ArrowLeft' && currentQuestionIndex > 0) {
        event.preventDefault();
        handlePrevious();
      } else if (key === 'ArrowRight' && currentQuestionIndex < questions.length - 1) {
        event.preventDefault();
        handleNext();
      }

      // Number keys for direct selection (1-5)
      const numValue = parseInt(key);
      if (numValue >= 1 && numValue <= 5) {
        event.preventDefault();
        const likertValue = numValue - 1; // Convert 1-5 to 0-4
        setSelectedOption(likertValue);
        setTimeout(() => {
          handleAnswerSelect(likertValue);
          setSelectedOption(null);
        }, 100);
      }

      // Enter key to confirm current selection
      if (key === 'Enter' && selectedOption !== null) {
        event.preventDefault();
        handleAnswerSelect(selectedOption);
        setSelectedOption(null);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [currentQuestionIndex, questions.length, isSubmitting, submitSuccess, selectedOption, handlePrevious, handleNext, handleAnswerSelect]);

  const handleSubmit = async () => {
    if (!assessment || !questions.length) return;

    // Validate that all questions are answered
    const unansweredQuestions = questions.filter(q => responses[q.id] === undefined);
    if (unansweredQuestions.length > 0) {
      // Show error and navigate to first unanswered question
      const firstUnansweredIndex = questions.findIndex(q => responses[q.id] === undefined);
      setCurrentQuestionIndex(firstUnansweredIndex);
      const unansweredLabel = unansweredQuestions.length === 1 ? "pregunta" : "preguntas";

      toast.error("Responde todas las preguntas antes de enviar", {
        description: `Tienes ${unansweredQuestions.length} ${unansweredLabel} sin responder. Te llevamos a la primera pendiente.`,
        duration: 5000,
      });
      return;
    }

    setIsSubmitting(true);

    // Convert responses to the expected format
    const formattedResponses = questions.map(question => ({
      question_id: question.id,
      value: responses[question.id] ?? 0, // Should not be needed now, but keeping as safety
    }));

    try {
      await submitMutation.mutateAsync(formattedResponses);
    } catch {
      setIsSubmitting(false); // Reset submitting state on error
      // Error is handled by the mutation
    }
  };

  // Check if current question is answered
  const isCurrentAnswered = currentQuestion ? responses[currentQuestion.id] !== undefined : false;
  const isLastQuestion = currentQuestionIndex === questions.length - 1;
  const canProceed = isCurrentAnswered || isLastQuestion;

  // Check if all questions are answered
  const allQuestionsAnswered = questions.length > 0 && questions.every(q => responses[q.id] !== undefined);
  const answeredCount = Object.keys(responses).length;
  const totalQuestions = questions.length;
  const answeredLabel = answeredCount === 1 ? "respondida" : "respondidas";
  const assessmentErrorMessage = resolveAssessmentErrorMessage(assessmentError);

  // Loading state
  if (isLoadingAssessment) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <Loader2 className="h-8 w-8 animate-spin mx-auto mb-4 text-muted-foreground" />
          <p className="text-muted-foreground">Cargando evaluacion...</p>
        </div>
      </div>
    );
  }

  // Error state
  if (assessmentError) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Card className="portal-surface-strong w-full max-w-md border-0">
          <CardContent className="pt-6">
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>
                {assessmentErrorMessage}
              </AlertDescription>
            </Alert>
          </CardContent>
        </Card>
      </div>
    );
  }

  // Success state after submission
  if (submitSuccess) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Card className="portal-surface-strong w-full max-w-md border-0 text-center">
          <CardContent className="pt-8 pb-8">
            <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-[var(--gradient-start)] to-[var(--gradient-end)] shadow-lg">
              <CheckCircle className="h-9 w-9 text-white" />
            </div>
            <h2 className="text-2xl font-semibold mb-2">Evaluacion completada</h2>
            <p className="text-muted-foreground mb-4">
              Gracias por completar la evaluacion NOM-035.
              Tus respuestas se enviaron correctamente.
            </p>
            <p className="text-sm text-muted-foreground">
              Ahora puedes cerrar esta ventana.
            </p>
          </CardContent>
        </Card>
      </div>
    );
  }

  // Invalid token in URL
  if (!token) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Card className="portal-surface-strong w-full max-w-md border-0">
          <CardContent className="pt-6">
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>
                Enlace de evaluacion invalido. Usa la URL completa enviada en tu invitacion.
              </AlertDescription>
            </Alert>
          </CardContent>
        </Card>
      </div>
    );
  }

  // No assessment data
  if (!assessment || !questions.length) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Card className="portal-surface-strong w-full max-w-md border-0">
          <CardContent className="pt-6">
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>
                No se pudo cargar la evaluacion. Solicita un nuevo enlace de evaluacion.
              </AlertDescription>
            </Alert>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen py-10">
      <div className="mx-auto w-full max-w-4xl space-y-8">
        {/* Progress Bar */}
        <div className="mx-auto max-w-3xl">
          <div className="assessment-progress-shell sticky top-2 z-10 rounded-xl px-5 py-4">
            <div className="mb-2 flex items-center justify-between">
              <span className="text-sm font-medium text-foreground">
                Pregunta {currentQuestionIndex + 1} de {questions.length}
              </span>
              <div className="flex items-center space-x-3">
                <div className="flex w-20 justify-end">
                  {saveStatus === "saving" && (
                    <div className="flex items-center space-x-1 text-primary">
                      <Loader2 className="h-3 w-3 animate-spin" />
                      <span className="text-xs">Guardando...</span>
                    </div>
                  )}
                  {saveStatus === "saved" && (
                    <div className="flex items-center space-x-1 text-emerald-600 dark:text-emerald-400">
                      <Cloud className="h-3 w-3" />
                      <span className="text-xs">Guardado</span>
                    </div>
                  )}
                  {saveStatus === "error" && (
                    <div className="flex items-center space-x-1 text-destructive">
                      <AlertCircle className="h-3 w-3" />
                      <span className="text-xs">Sin conexion</span>
                    </div>
                  )}
                </div>
                <span className={`text-sm ${allQuestionsAnswered ? "font-medium text-emerald-600 dark:text-emerald-400" : "text-muted-foreground"}`}>
                  {allQuestionsAnswered ? "Completado" : `${Math.round(progress)}% completado`}
                </span>
              </div>
            </div>
            <Progress value={progress} className="h-2.5 bg-secondary/70" />
          </div>
        </div>

        {/* Question Card */}
        <div className="mx-auto max-w-3xl">
          <AnimatePresence mode="wait">
            <motion.div
              key={currentQuestionIndex}
              initial={{ x: 50, opacity: 0 }}
              animate={{ x: 0, opacity: 1 }}
              exit={{ x: -50, opacity: 0 }}
              transition={{ duration: 0.3, ease: "easeInOut" }}
            >
              <Card className="portal-surface-strong border-0 shadow-sm">
                <CardContent className="py-8">
                  <div className="space-y-8 text-center">
                    {/* Category Display */}
                    <AnimatePresence>
                      {currentQuestion?.category?.name && (
                        <motion.div
                          key={currentQuestion.category.name}
                          initial={{ opacity: 0, y: -10 }}
                          animate={{ opacity: 1, y: 0 }}
                          exit={{ opacity: 0, y: -10 }}
                          transition={{ duration: 0.2 }}
                          className="mb-6 flex justify-center"
                        >
                          <Badge variant="outline" className="border-primary/20 bg-primary/5 text-xs font-medium uppercase tracking-wider text-primary">
                            {currentQuestion.category.name}
                          </Badge>
                        </motion.div>
                      )}
                    </AnimatePresence>

                    {/* Question Text */}
                    <div className="space-y-4 text-center">
                      <motion.h1
                        className="mb-8 text-2xl font-medium leading-tight text-foreground"
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        transition={{ delay: 0.1, duration: 0.3 }}
                      >
                        {currentQuestion?.text}
                      </motion.h1>
                      {currentQuestion?.subsection && (
                        <motion.p
                          className="text-sm text-muted-foreground"
                          initial={{ opacity: 0 }}
                          animate={{ opacity: 1 }}
                          transition={{ delay: 0.2, duration: 0.3 }}
                        >
                          {currentQuestion.subsection}
                        </motion.p>
                      )}
                    </div>

                    {/* Likert Scale Options */}
                    <motion.div
                      className="space-y-4"
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: 0.3, duration: 0.3 }}
                    >
                      <p className="text-sm font-medium text-foreground/90">
                        Selecciona tu respuesta (teclado: 1-5 o flechas)
                      </p>
                      <div className="mx-auto grid max-w-4xl auto-rows-fr grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-5">
                        {LIKERT_OPTIONS.map((option) => {
                          const isSelected = responses[currentQuestion?.id] === option.value;
                          const isPressed = pressedKey === option.key || selectedOption === option.value;

                          return (
                            <motion.div
                              key={option.value}
                              whileTap={{ scale: 0.95 }}
                              transition={{ type: "spring", stiffness: 300, damping: 20 }}
                            >
                              <Button
                                variant={isSelected ? "default" : "outline"}
                                className={`relative h-full min-h-[132px] w-full flex-col items-start justify-start gap-2 whitespace-normal rounded-lg px-3 py-3 text-left transition-all hover:shadow-md ${
                                  isSelected
                                    ? "ring-2 ring-primary ring-offset-2"
                                    : isPressed
                                      ? "ring-2 ring-primary/70 ring-offset-2"
                                      : ""
                                }`}
                                onClick={() => handleAnswerSelect(option.value)}
                                disabled={isSubmitting}
                              >
                                <span className="absolute right-2 top-2 inline-flex h-5 w-5 items-center justify-center rounded-md border border-border/70 bg-background/70 font-mono text-[10px] leading-none text-muted-foreground/80">
                                  {option.key}
                                </span>
                                <div className="w-full space-y-1 pr-7">
                                  <span className="block break-words text-sm font-semibold leading-snug">{option.label}</span>
                                  <span className={`block break-words text-xs leading-snug ${isSelected ? "text-primary-foreground/80" : "text-muted-foreground"}`}>
                                    {option.description}
                                  </span>
                                </div>
                              </Button>
                            </motion.div>
                          );
                        })}
                      </div>
                    </motion.div>
                  </div>
                </CardContent>
              </Card>
            </motion.div>
          </AnimatePresence>

          {/* Navigation */}
          <motion.div
            className="mx-auto mt-6 flex max-w-3xl items-center justify-between"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.4, duration: 0.3 }}
          >
            <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
              <Button
                variant="outline"
                onClick={handlePrevious}
                disabled={currentQuestionIndex === 0 || isSubmitting}
                className="space-x-2 border-primary/20 bg-background/60"
              >
                <ChevronLeft className="h-4 w-4" />
                <span>Anterior</span>
              </Button>
            </motion.div>

            <div className={`text-center text-sm ${allQuestionsAnswered ? "font-medium text-emerald-600 dark:text-emerald-400" : "text-muted-foreground"}`}>
              {answeredCount} de {totalQuestions} {answeredLabel}
              {allQuestionsAnswered && <span className="ml-1">✓</span>}
            </div>

            {isLastQuestion ? (
              <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
                <Button
                  onClick={handleSubmit}
                  disabled={!allQuestionsAnswered || isSubmitting}
                  className="space-x-2 bg-gradient-to-r from-[var(--gradient-start)] to-[var(--gradient-end)] text-white"
                >
                  {isSubmitting ? (
                    <>
                      <Loader2 className="h-4 w-4 animate-spin" />
                      <span>Enviando...</span>
                    </>
                  ) : (
                    <>
                      <span>Enviar evaluacion</span>
                      <CheckCircle className="h-4 w-4" />
                    </>
                  )}
                </Button>
              </motion.div>
            ) : (
              <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
                <Button
                  variant="outline"
                  onClick={handleNext}
                  disabled={!canProceed || isSubmitting}
                  className="space-x-2 border-primary/20 bg-background/60"
                >
                  <span>Siguiente</span>
                  <ChevronRight className="h-4 w-4" />
                </Button>
              </motion.div>
            )}
          </motion.div>

          {/* Keyboard Legend */}
          <div className="mt-6 flex items-center justify-center space-x-6 text-xs text-muted-foreground">
            <div className="flex items-center space-x-1">
              <ChevronLeft className="h-3 w-3" />
              <ChevronRight className="h-3 w-3" />
              <span>Navegar</span>
            </div>
            <div className="flex items-center space-x-1">
              <span className="font-mono">1-5</span>
              <span>Seleccionar</span>
            </div>
            <div className="flex items-center space-x-1">
              <span className="font-mono">Intro</span>
              <span>Confirmar</span>
            </div>
          </div>
        </div>

        {/* Submit Error */}
        {submitMutation.isError && (
          <div className="mx-auto max-w-2xl">
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>
                No se pudo enviar la evaluacion. Intentalo de nuevo.
              </AlertDescription>
            </Alert>
          </div>
        )}
      </div>
    </div>
  );
}

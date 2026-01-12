"use client";

import React, { useState, useEffect, useCallback } from "react";
import { useParams, useRouter } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { motion, AnimatePresence } from "framer-motion";
import { ChevronLeft, ChevronRight, CheckCircle, AlertCircle, Loader2, Cloud } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { assessmentService } from "@/services/assessment.service";
import { toast } from "sonner";
import type { Question } from "@/types/backend";

/**
 * Public Assessment Exam Page
 *
 * Distraction-free interface for staff to complete NOM-035 assessments.
 * Features progress tracking, question navigation, and seamless submission.
 */
export default function AssessmentExamPage() {
  const params = useParams();
  const router = useRouter();
  const token = params.token as string;

  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [responses, setResponses] = useState<Record<number, number>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitSuccess, setSubmitSuccess] = useState(false);
  const [selectedOption, setSelectedOption] = useState<number | null>(null);
  const [pressedKey, setPressedKey] = useState<string | null>(null);
  const [saveStatus, setSaveStatus] = useState<'saved' | 'saving' | 'error'>('saved');
  const [currentCategory, setCurrentCategory] = useState<string>('');

  // Fetch assessment data and questions
  const {
    data: assessmentData,
    isLoading: isLoadingAssessment,
    error: assessmentError,
  } = useQuery({
    queryKey: ["public-assessment", token],
    queryFn: () => assessmentService.fetchByToken(token),
    enabled: !!token,
  });

  // Submit responses mutation
  const submitMutation = useMutation({
    mutationFn: (responses: { question_id: number; value: number }[]) =>
      assessmentService.submitResponses(token, responses),
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

  // Update category when question changes
  useEffect(() => {
    if (currentQuestion?.category?.name) {
      setCurrentCategory(currentQuestion.category.name);
    }
  }, [currentQuestion?.category?.name]);

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
  }, [currentQuestionIndex, questions.length, isSubmitting, submitSuccess, selectedOption]);

  // Auto-advance to next question after selection (with delay)
  const handleAnswerSelect = useCallback((value: number) => {
    if (!currentQuestion) return;

    setResponses(prev => ({
      ...prev,
      [currentQuestion.id]: value,
    }));

    // Update save status
    setSaveStatus('saving');
    setTimeout(() => setSaveStatus('saved'), 500);

    // Auto-advance after a brief delay
    setTimeout(() => {
      if (currentQuestionIndex < questions.length - 1) {
        setCurrentQuestionIndex(prev => prev + 1);
      }
    }, 300);
  }, [currentQuestion, currentQuestionIndex, questions.length]);

  const handlePrevious = () => {
    if (currentQuestionIndex > 0) {
      setCurrentQuestionIndex(prev => prev - 1);
    }
  };

  const handleNext = () => {
    if (currentQuestionIndex < questions.length - 1) {
      setCurrentQuestionIndex(prev => prev + 1);
    }
  };

  const handleSubmit = async () => {
    if (!assessment || !questions.length) return;

    // Validate that all questions are answered
    const unansweredQuestions = questions.filter(q => responses[q.id] === undefined);
    if (unansweredQuestions.length > 0) {
      // Show error and navigate to first unanswered question
      const firstUnansweredIndex = questions.findIndex(q => responses[q.id] === undefined);
      setCurrentQuestionIndex(firstUnansweredIndex);

      toast.error("Please answer all questions before submitting", {
        description: `You have ${unansweredQuestions.length} unanswered question(s). We've navigated to the first one.`,
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
    } catch (error) {
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

  // Loading state
  if (isLoadingAssessment) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <Loader2 className="h-8 w-8 animate-spin mx-auto mb-4 text-muted-foreground" />
          <p className="text-muted-foreground">Loading assessment...</p>
        </div>
      </div>
    );
  }

  // Error state
  if (assessmentError) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Card className="w-full max-w-md">
          <CardContent className="pt-6">
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>
                {assessmentError.message || "Failed to load assessment. The link may be invalid or expired."}
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
        <Card className="w-full max-w-md text-center">
          <CardContent className="pt-8 pb-8">
            <CheckCircle className="h-16 w-16 text-green-500 mx-auto mb-6" />
            <h2 className="text-2xl font-semibold mb-2">Assessment Completed</h2>
            <p className="text-muted-foreground mb-4">
              Thank you for completing the NOM-035 assessment.
              Your responses have been submitted successfully.
            </p>
            <p className="text-sm text-muted-foreground">
              You can now close this window.
            </p>
          </CardContent>
        </Card>
      </div>
    );
  }

  // No assessment data
  if (!assessment || !questions.length) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Card className="w-full max-w-md">
          <CardContent className="pt-6">
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>
                Assessment not found or no questions available.
              </AlertDescription>
            </Alert>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col justify-center py-12">
      <div className="max-w-4xl mx-auto w-full space-y-8">
        {/* Progress Bar */}
        <div className="max-w-2xl mx-auto">
          <div className="sticky top-0 bg-white/95 backdrop-blur-sm border-b border-slate-200 -mx-6 px-6 py-4">
        <div className="flex items-center justify-between mb-2">
          <span className="text-sm font-medium text-slate-700">
            Question {currentQuestionIndex + 1} of {questions.length}
          </span>
          <div className="flex items-center space-x-3">
            <div className="w-20 flex justify-end">
              {saveStatus === 'saving' && (
                <div className="flex items-center space-x-1 text-blue-600">
                  <Loader2 className="h-3 w-3 animate-spin" />
                  <span className="text-xs">Saving...</span>
                </div>
              )}
              {saveStatus === 'saved' && (
                <div className="flex items-center space-x-1 text-green-600">
                  <Cloud className="h-3 w-3" />
                  <span className="text-xs">Saved</span>
                </div>
              )}
              {saveStatus === 'error' && (
                <div className="flex items-center space-x-1 text-red-600">
                  <AlertCircle className="h-3 w-3" />
                  <span className="text-xs">Offline</span>
                </div>
              )}
            </div>
            <span className={`text-sm ${allQuestionsAnswered ? 'text-green-600 font-medium' : 'text-slate-500'}`}>
              {allQuestionsAnswered ? 'Complete!' : `${Math.round(progress)}% Complete`}
            </span>
          </div>
        </div>
            <Progress value={progress} className="h-2" />
          </div>
        </div>

        {/* Question Card */}
        <div className="max-w-2xl mx-auto">
      <AnimatePresence mode="wait">
        <motion.div
          key={currentQuestionIndex}
          initial={{ x: 50, opacity: 0 }}
          animate={{ x: 0, opacity: 1 }}
          exit={{ x: -50, opacity: 0 }}
          transition={{ duration: 0.3, ease: "easeInOut" }}
        >
          <Card className="shadow-sm">
            <CardContent className="pt-8 pb-8">
              <div className="text-center space-y-8">
                {/* Category Display */}
                <AnimatePresence>
                  {currentCategory && (
                    <motion.div
                      key={currentCategory}
                      initial={{ opacity: 0, y: -10 }}
                      animate={{ opacity: 1, y: 0 }}
                      exit={{ opacity: 0, y: -10 }}
                      transition={{ duration: 0.2 }}
                      className="flex justify-center mb-6"
                    >
                      <Badge variant="outline" className="text-xs font-medium uppercase tracking-wider">
                        {currentCategory}
                      </Badge>
                    </motion.div>
                  )}
                </AnimatePresence>

                {/* Question Text */}
                <div className="text-center space-y-4">
                  <motion.h1
                    className="text-2xl font-medium text-slate-900 leading-tight mb-8"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ delay: 0.1, duration: 0.3 }}
                  >
                    {currentQuestion?.text}
                  </motion.h1>
                  {currentQuestion?.subsection && (
                    <motion.p
                      className="text-sm text-slate-600"
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
                  <p className="text-sm font-medium text-slate-700">
                    Select your response: (Use keyboard: 1-5 or arrow keys)
                  </p>
                  <div className="grid grid-cols-1 sm:grid-cols-3 md:grid-cols-5 gap-4 max-w-2xl mx-auto">
                    {[
                      { value: 0, label: "Siempre", description: "Always", key: "1" },
                      { value: 1, label: "Casi Siempre", description: "Almost Always", key: "2" },
                      { value: 2, label: "Algunas Veces", description: "Sometimes", key: "3" },
                      { value: 3, label: "Casi Nunca", description: "Almost Never", key: "4" },
                      { value: 4, label: "Nunca", description: "Never", key: "5" },
                    ].map((option) => {
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
                            className={`h-auto py-6 px-3 flex flex-col items-center justify-center space-y-2 text-center hover:shadow-md transition-all relative min-h-[110px] w-full ${
                              isSelected ? 'ring-2 ring-primary ring-offset-2' : isPressed ? 'ring-2 ring-blue-500 ring-offset-2' : ''
                            }`}
                            onClick={() => handleAnswerSelect(option.value)}
                            disabled={isSubmitting}
                          >
                            <span className="font-medium text-sm leading-tight break-words">{option.label}</span>
                            <span className="text-xs opacity-75 leading-tight break-words">{option.description}</span>
                            <span className="absolute top-2 right-2 text-[10px] font-mono text-muted-foreground/50">
                              {option.key}
                            </span>
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
          className="flex items-center justify-between max-w-2xl mx-auto"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4, duration: 0.3 }}
        >
          <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
            <Button
              variant="outline"
              onClick={handlePrevious}
              disabled={currentQuestionIndex === 0 || isSubmitting}
              className="flex items-center space-x-2"
            >
              <ChevronLeft className="h-4 w-4" />
              <span>Previous</span>
            </Button>
          </motion.div>

          <div className={`text-center text-sm ${allQuestionsAnswered ? 'text-green-600 font-medium' : 'text-slate-500'}`}>
            {answeredCount} of {totalQuestions} answered
            {allQuestionsAnswered && <span className="ml-1">✓</span>}
          </div>

          {isLastQuestion ? (
            <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
              <Button
                onClick={handleSubmit}
                disabled={!allQuestionsAnswered || isSubmitting}
                className="flex items-center space-x-2"
              >
                {isSubmitting ? (
                  <>
                    <Loader2 className="h-4 w-4 animate-spin" />
                    <span>Submitting...</span>
                  </>
                ) : (
                  <>
                    <span>Submit Assessment</span>
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
                className="flex items-center space-x-2"
              >
                <span>Next</span>
                <ChevronRight className="h-4 w-4" />
              </Button>
            </motion.div>
          )}
        </motion.div>

        {/* Keyboard Legend */}
        <div className="flex items-center justify-center space-x-6 mt-6 text-xs text-muted-foreground">
          <div className="flex items-center space-x-1">
            <ChevronLeft className="w-3 h-3" />
            <ChevronRight className="w-3 h-3" />
            <span>Navigate</span>
          </div>
          <div className="flex items-center space-x-1">
            <span className="font-mono">1-5</span>
            <span>Select</span>
          </div>
          <div className="flex items-center space-x-1">
            <span className="font-mono">Enter</span>
            <span>Confirm</span>
          </div>
        </div>
      </div>
    </div>

      {/* Submit Error */}
      {submitMutation.isError && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>
            Failed to submit assessment. Please try again.
          </AlertDescription>
        </Alert>
      )}
    </div>
  );
}
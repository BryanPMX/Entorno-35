"use client";

import React, { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { ChevronLeft, ChevronRight, CheckCircle, AlertCircle, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { assessmentService } from "@/services/assessment.service";
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

  // Auto-advance to next question after selection (with delay)
  const handleAnswerSelect = (value: number) => {
    if (!currentQuestion) return;

    setResponses(prev => ({
      ...prev,
      [currentQuestion.id]: value,
    }));

    // Auto-advance after a brief delay
    setTimeout(() => {
      if (currentQuestionIndex < questions.length - 1) {
        setCurrentQuestionIndex(prev => prev + 1);
      }
    }, 300);
  };

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

    setIsSubmitting(true);

    // Convert responses to the expected format
    const formattedResponses = questions.map(question => ({
      question_id: question.id,
      value: responses[question.id] ?? 0, // Default to 0 if not answered
    }));

    try {
      await submitMutation.mutateAsync(formattedResponses);
    } catch (error) {
      // Error is handled by the mutation
    }
  };

  // Check if current question is answered
  const isCurrentAnswered = currentQuestion ? responses[currentQuestion.id] !== undefined : false;
  const isLastQuestion = currentQuestionIndex === questions.length - 1;
  const canProceed = isCurrentAnswered || isLastQuestion;

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
    <div className="max-w-4xl mx-auto space-y-8">
      {/* Progress Bar */}
      <div className="sticky top-0 bg-white/95 backdrop-blur-sm border-b border-slate-200 -mx-6 px-6 py-4">
        <div className="flex items-center justify-between mb-2">
          <span className="text-sm font-medium text-slate-700">
            Question {currentQuestionIndex + 1} of {questions.length}
          </span>
          <span className="text-sm text-slate-500">
            {Math.round(progress)}% Complete
          </span>
        </div>
        <Progress value={progress} className="h-2" />
      </div>

      {/* Question Card */}
      <Card className="shadow-sm">
        <CardContent className="pt-8 pb-8">
          <div className="text-center space-y-8">
            {/* Question Text */}
            <div className="space-y-4">
              <h1 className="text-xl font-medium text-slate-900 leading-relaxed">
                {currentQuestion?.text}
              </h1>
              {currentQuestion?.subsection && (
                <p className="text-sm text-slate-600">
                  {currentQuestion.subsection}
                </p>
              )}
            </div>

            {/* Likert Scale Options */}
            <div className="space-y-4">
              <p className="text-sm font-medium text-slate-700">Select your response:</p>
              <div className="grid grid-cols-1 sm:grid-cols-5 gap-3 max-w-2xl mx-auto">
                {[
                  { value: 0, label: "Siempre", description: "Always" },
                  { value: 1, label: "Casi Siempre", description: "Almost Always" },
                  { value: 2, label: "Algunas Veces", description: "Sometimes" },
                  { value: 3, label: "Casi Nunca", description: "Almost Never" },
                  { value: 4, label: "Nunca", description: "Never" },
                ].map((option) => (
                  <Button
                    key={option.value}
                    variant={responses[currentQuestion?.id] === option.value ? "default" : "outline"}
                    className="h-auto py-4 px-3 flex flex-col items-center space-y-1 text-center hover:shadow-md transition-all"
                    onClick={() => handleAnswerSelect(option.value)}
                    disabled={isSubmitting}
                  >
                    <span className="font-medium text-sm">{option.label}</span>
                    <span className="text-xs opacity-75">{option.description}</span>
                  </Button>
                ))}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Navigation */}
      <div className="flex items-center justify-between">
        <Button
          variant="outline"
          onClick={handlePrevious}
          disabled={currentQuestionIndex === 0 || isSubmitting}
          className="flex items-center space-x-2"
        >
          <ChevronLeft className="h-4 w-4" />
          <span>Previous</span>
        </Button>

        <div className="text-sm text-slate-500">
          {Object.keys(responses).length} of {questions.length} answered
        </div>

        {isLastQuestion ? (
          <Button
            onClick={handleSubmit}
            disabled={!canProceed || isSubmitting}
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
        ) : (
          <Button
            variant="outline"
            onClick={handleNext}
            disabled={!canProceed || isSubmitting}
            className="flex items-center space-x-2"
          >
            <span>Next</span>
            <ChevronRight className="h-4 w-4" />
          </Button>
        )}
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
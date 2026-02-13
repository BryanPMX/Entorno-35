"use client";

import React, { useState, useCallback } from "react";
import { Upload, X, FileText, AlertCircle, CheckCircle, Download } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { staffService } from "@/services/staff.service";
import type { ImportResult } from "@/types/backend";

interface CsvUploadModalProps {
  isOpen: boolean;
  onClose: () => void;
  onUploadComplete: () => void;
}

export function CsvUploadModal({ isOpen, onClose, onUploadComplete }: CsvUploadModalProps) {
  const [dragActive, setDragActive] = useState(false);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [uploadStep, setUploadStep] = useState<'select' | 'upload' | 'complete'>('select');
  const [uploadProgress, setUploadProgress] = useState(0);
  const [uploadResult, setUploadResult] = useState<ImportResult | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleDrag = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true);
    } else if (e.type === "dragleave") {
      setDragActive(false);
    }
  }, []);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);

    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      const file = e.dataTransfer.files[0];
      if (file.type === "text/csv" || file.name.endsWith('.csv')) {
        setSelectedFile(file);
        setError(null);
      } else {
        setError("Please select a valid CSV file.");
      }
    }
  }, []);

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      if (file.type === "text/csv" || file.name.endsWith('.csv')) {
        setSelectedFile(file);
        setError(null);
      } else {
        setError("Please select a valid CSV file.");
      }
    }
  };

  const uploadCsv = async () => {
    if (!selectedFile) return;

    try {
      setUploadStep('upload');
      setUploadProgress(0);

      // Simulate progress
      const progressInterval = setInterval(() => {
        setUploadProgress(prev => {
          if (prev >= 90) {
            clearInterval(progressInterval);
            return prev;
          }
          return prev + 10;
        });
      }, 200);

      const result = await staffService.uploadCSV(selectedFile);
      setUploadProgress(100);
      setUploadResult(result);
      setUploadStep('complete');
      onUploadComplete();
    } catch (err) {
      setError("Failed to upload CSV file. Please try again.");
      console.error("CSV upload error:", err);
      setUploadStep('select');
    }
  };

  const resetModal = () => {
    setSelectedFile(null);
    setUploadStep('select');
    setUploadProgress(0);
    setUploadResult(null);
    setError(null);
  };

  const handleClose = () => {
    resetModal();
    onClose();
  };

  const downloadTemplate = () => {
    // Create a sample CSV template
    const csvContent = "Name,CURP,Email,Area,Job,Shift,Gender\nJuan Pérez García,ABCD123456HIJKLM01,juan.perez@example.com,Producción,Operador,Diurno,Masculino\nMaría González López,EFGH567890MNOPQR02,maria.gonzalez@example.com,Recursos Humanos,Analista,Diurno,Femenino\n";
    const blob = new Blob([csvContent], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'staff_import_template.csv';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
  };

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <Upload className="h-5 w-5" />
            <span>Import Staff from CSV</span>
          </DialogTitle>
          <DialogDescription>
            Upload a CSV file to bulk import staff members. Download the template for the correct format.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          {/* Template Download */}
          <div className="flex justify-between items-center">
            <p className="text-sm text-muted-foreground">
              Need a template? Download the CSV template with the correct column format.
            </p>
            <Button variant="outline" onClick={downloadTemplate} className="flex items-center space-x-2">
              <Download className="h-4 w-4" />
              <span>Download Template</span>
            </Button>
          </div>

          {/* Error Display */}
          {error && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          {/* File Selection Step */}
          {uploadStep === 'select' && (
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Select CSV File</CardTitle>
                <CardDescription>
                  Choose a CSV file containing your staff data. The file should include columns for Name, CURP, Email, Area, Job, Shift, and Gender.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div
                  className={`border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
                    dragActive ? 'border-primary bg-primary/5' : 'border-muted-foreground/25'
                  }`}
                  onDragEnter={handleDrag}
                  onDragLeave={handleDrag}
                  onDragOver={handleDrag}
                  onDrop={handleDrop}
                >
                  <FileText className="mx-auto h-12 w-12 text-muted-foreground mb-4" />
                  <div className="space-y-2">
                    <p className="text-lg font-medium">
                      {selectedFile ? selectedFile.name : "Drop your CSV file here"}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      {selectedFile
                        ? `${(selectedFile.size / 1024).toFixed(1)} KB`
                        : "or click to browse files"
                      }
                    </p>
                  </div>
                  {!selectedFile && (
                    <div className="mt-4">
                      <input
                        type="file"
                        accept=".csv"
                        onChange={handleFileSelect}
                        className="hidden"
                        id="csv-file-input"
                      />
                      <label htmlFor="csv-file-input">
                        <Button variant="outline" className="cursor-pointer">
                          Choose File
                        </Button>
                      </label>
                    </div>
                  )}
                  {selectedFile && (
                    <div className="mt-4 flex justify-center space-x-2">
                      <Button onClick={() => setSelectedFile(null)} variant="outline">
                        <X className="h-4 w-4 mr-2" />
                        Remove
                      </Button>
                      <Button onClick={uploadCsv}>
                        <Upload className="h-4 w-4 mr-2" />
                        Upload & Import
                      </Button>
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          )}

          {/* Upload Progress Step */}
          {uploadStep === 'upload' && (
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Uploading and Processing</CardTitle>
                <CardDescription>
                  Please wait while we process your CSV file and import the staff data.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Processing CSV data...</span>
                      <span>{uploadProgress}%</span>
                    </div>
                    <Progress value={uploadProgress} className="w-full" />
                  </div>
                  <p className="text-sm text-muted-foreground">
                    This may take a few moments depending on the size of your file.
                  </p>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Completion Step */}
          {uploadStep === 'complete' && uploadResult && (
            <Card>
              <CardHeader>
                <CardTitle className="text-lg flex items-center space-x-2">
                  <CheckCircle className="h-5 w-5 text-green-500" />
                  <span>Import Complete</span>
                </CardTitle>
                <CardDescription>
                  Your CSV file has been successfully processed.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-3 gap-4 mb-6">
                  <div className="text-center p-4 bg-green-50 rounded-lg">
                    <div className="text-2xl font-bold text-green-600">
                      {uploadResult.success_count}
                    </div>
                    <div className="text-sm text-green-600">Successfully Imported</div>
                  </div>
                  <div className="text-center p-4 bg-yellow-50 rounded-lg">
                    <div className="text-2xl font-bold text-yellow-600">
                      {uploadResult.skipped_count}
                    </div>
                    <div className="text-sm text-yellow-600">Skipped (Duplicates)</div>
                  </div>
                  <div className="text-center p-4 bg-red-50 rounded-lg">
                    <div className="text-2xl font-bold text-red-600">
                      {uploadResult.errors.length}
                    </div>
                    <div className="text-sm text-red-600">Errors</div>
                  </div>
                </div>

                {uploadResult.errors.length > 0 && (
                  <div className="space-y-2">
                    <h4 className="font-medium text-red-600">Errors Encountered:</h4>
                    <div className="max-h-32 overflow-y-auto bg-red-50 p-3 rounded border">
                      {uploadResult.errors.map((error, index) => (
                        <div key={index} className="text-sm text-red-600">
                          {error}
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                <div className="flex justify-end space-x-2 mt-6">
                  <Button variant="outline" onClick={handleClose}>
                    Close
                  </Button>
                  <Button onClick={resetModal}>
                    Import Another File
                  </Button>
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}

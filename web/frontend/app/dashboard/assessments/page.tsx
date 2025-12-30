"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

/**
 * Assessments Page
 * 
 * Placeholder page for assessments functionality.
 * Will be implemented in Phase 5.4
 */
export default function AssessmentsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Assessments</h1>
        <p className="text-gray-600 mt-2">
          Create and manage NOM-035 assessments
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Assessment List</CardTitle>
          <CardDescription>
            View and manage your assessments
          </CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-gray-600">
            Assessment management functionality will be implemented in Phase 5.4
          </p>
        </CardContent>
      </Card>
    </div>
  );
}


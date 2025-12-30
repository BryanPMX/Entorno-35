"use client";

import { useAuthStore } from "@/lib/store/auth-store";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

/**
 * Dashboard Page
 * 
 * Main dashboard page for authenticated users.
 * Displays overview and navigation to key features.
 */
export default function DashboardPage() {
  const { user } = useAuthStore();

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-gray-600 mt-2">
          Welcome to Entorno 35 - NOM-035 Compliance Platform
        </p>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle>Staff Management</CardTitle>
            <CardDescription>
              Manage your staff members and import via CSV
            </CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-600">
              View, add, and import staff members for assessments.
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Assessments</CardTitle>
            <CardDescription>
              Create and manage NOM-035 assessments
            </CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-600">
              Create assessment links and track completion status.
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Reports</CardTitle>
            <CardDescription>
              View individual and general reports
            </CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-600">
              Generate compliance reports and recommendations.
            </p>
          </CardContent>
        </Card>
      </div>

      {user && (
        <Card>
          <CardHeader>
            <CardTitle>Account Information</CardTitle>
          </CardHeader>
          <CardContent>
            <dl className="space-y-2">
              <div>
                <dt className="text-sm font-medium text-gray-500">Role</dt>
                <dd className="text-sm text-gray-900 capitalize">{user.role}</dd>
              </div>
              {user.email && (
                <div>
                  <dt className="text-sm font-medium text-gray-500">Email</dt>
                  <dd className="text-sm text-gray-900">{user.email}</dd>
                </div>
              )}
              <div>
                <dt className="text-sm font-medium text-gray-500">Company ID</dt>
                <dd className="text-sm text-gray-900 font-mono">{user.company_id}</dd>
              </div>
              {user.staff_id && (
                <div>
                  <dt className="text-sm font-medium text-gray-500">Staff ID</dt>
                  <dd className="text-sm text-gray-900 font-mono">{user.staff_id}</dd>
                </div>
              )}
            </dl>
          </CardContent>
        </Card>
      )}
    </div>
  );
}


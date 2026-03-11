package http

import "testing"

func TestIsStripeSubscriptionActive(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		status   string
		expected bool
	}{
		{name: "active", status: "active", expected: true},
		{name: "active uppercase", status: "ACTIVE", expected: true},
		{name: "trialing", status: "trialing", expected: false},
		{name: "incomplete", status: "incomplete", expected: false},
		{name: "past due", status: "past_due", expected: false},
		{name: "empty", status: "", expected: false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := isStripeSubscriptionActive(tc.status)
			if actual != tc.expected {
				t.Fatalf("expected %v, got %v for status %q", tc.expected, actual, tc.status)
			}
		})
	}
}

func TestCanTransitionRefundRequestStatus(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		current  string
		target   string
		expected bool
	}{
		{name: "requested to approved", current: "requested", target: "approved", expected: true},
		{name: "requested to rejected", current: "requested", target: "rejected", expected: true},
		{name: "requested to refunded", current: "requested", target: "refunded", expected: true},
		{name: "approved to refunded", current: "approved", target: "refunded", expected: true},
		{name: "approved to rejected", current: "approved", target: "rejected", expected: true},
		{name: "approved to approved", current: "approved", target: "approved", expected: false},
		{name: "rejected to approved", current: "rejected", target: "approved", expected: false},
		{name: "refunded to rejected", current: "refunded", target: "rejected", expected: false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := canTransitionRefundRequestStatus(tc.current, tc.target)
			if actual != tc.expected {
				t.Fatalf("expected %v, got %v for %s -> %s", tc.expected, actual, tc.current, tc.target)
			}
		})
	}
}

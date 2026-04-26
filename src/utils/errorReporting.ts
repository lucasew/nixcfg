/**
 * Centralized error reporting utility.
 *
 * As described by Martin Fowler in "Refactoring", centralizing error handling
 * extracts the error reporting responsibility (Single Responsibility Principle)
 * into a single module. This reduces code duplication and ensures that all errors
 * are processed consistently across the codebase.
 */

export function reportError(error: Error, context?: Record<string, unknown>): void {
  // Funnel all unexpected errors through this single function.
  // This call site does not know or care which backend is active.

  const payload = {
    message: error.message,
    stack: error.stack,
    context: context || {},
    timestamp: new Date().toISOString(),
  };

  // Log the error with enough context. If Sentry or another backend
  // is added later, it will be integrated here without changing call sites.
  console.error(JSON.stringify(payload));
}

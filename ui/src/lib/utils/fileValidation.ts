/**
 * Maximum file size in bytes (10MB)
 */
export const MAX_FILE_SIZE = 30 * 1024 * 1024;

/**
 * Validates file size and returns error message if invalid
 * @param file The file to validate
 * @returns Error message string if invalid, empty string if valid
 */
export function validateFileSize(file: File | null): string {
  if (!file) return "";

  if (file.size > MAX_FILE_SIZE) {
    return `File size exceeds the limit of ${formatFileSize(MAX_FILE_SIZE)}`;
  }

  return "";
}

/**
 * Formats file size in bytes to a human-readable string
 * @param bytes Size in bytes
 * @returns Formatted string (e.g., "10 MB")
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${Number.parseFloat((bytes / (k ** i)).toFixed(2))} ${sizes[i]}`;
}

/**
 * Returns the maximum file size as a human-readable string
 */
export function getMaxFileSizeDisplay(): string {
  return formatFileSize(MAX_FILE_SIZE);
}

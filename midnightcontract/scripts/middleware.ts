export function validateRequest<T>(schema: { safeParse: (data: any) => { success: boolean; data?: T; error?: any } }, data: any) {
  const parsed = schema.safeParse(data);
  if (!parsed.success) throw new Error('Validation failed');
  return parsed.data as T;
}

export function transformFormat<T>(payload: T): T {
  return payload;
}

export function normalizeResponse<T>(payload: T): T {
  return payload;
}


import type { Solution } from '../types/solution';

type HttpResponse<T> = {
  status: number;
  data: T;
};

async function parseResponse(res: Response): Promise<any> {
  const contentType = res.headers.get('content-type') || '';
  if (contentType.includes('application/json')) {
    return res.json();
  }

  try {
    return await res.json();
  } catch {
    return res.text();
  }
}

async function request<T>(url: string, init?: RequestInit): Promise<HttpResponse<T>> {
  const res = await fetch(url, init);
  const data = await parseResponse(res);
  return { status: res.status, data };
}

export async function getHealth(): Promise<string> {
  const { data } = await request<{ status: string }>('/api/health');
  return data.status;
}

export async function getSolutions(): Promise<Solution[]> {
  const { data } = await request<Solution[]>('/api/solutions');
  return data;
}

export async function requestSolutionEndpoint(args: {
  endpoint: string;
  method: string;
  body?: string;
}): Promise<HttpResponse<any>> {
  const init: RequestInit = {
    method: args.method,
    headers: { 'Content-Type': 'application/json' },
  };

  if (args.method !== 'GET' && args.body !== undefined) {
    init.body = args.body;
  }

  return request<any>(args.endpoint, init);
}

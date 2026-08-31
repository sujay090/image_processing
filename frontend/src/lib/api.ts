// A simple fetch wrapper that automatically handles the JWT token

export const getAuthToken = () => localStorage.getItem('token');
export const setAuthToken = (token: string) => localStorage.setItem('token', token);
export const removeAuthToken = () => localStorage.removeItem('token');

interface FetchOptions extends RequestInit {
  data?: any;
}

export async function apiFetch(endpoint: string, options: FetchOptions = {}) {
  const { data, headers: customHeaders, ...customConfig } = options;
  const token = getAuthToken();

  const headers: Record<string, string> = {
    ...(customHeaders as Record<string, string>),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  // Only set Content-Type to application/json if data is provided and it's not a FormData
  if (data && !(data instanceof FormData)) {
    headers['Content-Type'] = 'application/json';
  }

  const config: RequestInit = {
    ...customConfig,
    headers,
  };

  if (data) {
    config.body = data instanceof FormData ? data : JSON.stringify(data);
  }

  const response = await fetch(`/api/v1${endpoint}`, config);

  if (response.status === 401) {
    removeAuthToken();
    window.location.href = '/login';
    return Promise.reject('Unauthorized');
  }

  const isJson = response.headers.get('content-type')?.includes('application/json');
  const responseData = isJson ? await response.json() : await response.text();

  if (response.ok) {
    return responseData;
  } else {
    return Promise.reject(responseData);
  }
}

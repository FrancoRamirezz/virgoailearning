export type LoginResponse = {
  message?: string;
  token?: string;
  user?: unknown;
};
export type AuthResult = { message?: string };


export async function login(role: string,email: string, password: string, rememberMe: boolean): Promise<LoginResponse> {
  const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL;
  const res = await fetch(`${API_BASE}/api/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include", // important if backend uses cookies
    body: JSON.stringify({ role, email, password, rememberMe }),
  });

  const data = (await res.json()) as LoginResponse;

  if (!res.ok) {
    throw new Error(data.message ?? "Invalid credentials");
  }

  return data;
}

export async function signup(params: {
  role: string;
  email: string;
  password: string;
  name?: string;
  rememberMe?: boolean;
}): Promise<AuthResult> {
  const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL;

  const res = await fetch(`${API_BASE}/api/auth/signup`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include", // cookie-based session
    body: JSON.stringify(params),
  });

  const data = (await res.json()) as AuthResult;

  if (!res.ok) {
    throw new Error(data.message ?? "Unable to create account");
  }

  return data;
}


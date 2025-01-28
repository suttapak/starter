import { NextRequest, NextResponse } from "next/server";
import { cookies } from "next/headers";
import { jwtDecode } from "jwt-decode";
const isValidToken = (accessToken: string) => {
  if (!accessToken) {
    return false;
  }
  const decodedToken = jwtDecode(accessToken);
  const currentTime = Date.now() / 1000;

  if (!decodedToken.exp) {
    return false;
  }

  return decodedToken.exp > currentTime;
};
// 1. Specify protected and public routes
const protectedRoutes = ["/stock-service"];
const publicRoutes = ["/login", "/signup", "/"];

export default async function middleware(req: NextRequest) {
  // 2. Check if the current route is protected or public
  const path = req.nextUrl.pathname;

  const isPublicRoute = publicRoutes.includes(path);

  const isProtectedRoute = protectedRoutes.some((route) =>
    path.startsWith(route),
  );

  // 3. Decrypt the session from the cookie
  const session = (await cookies()).get("session")?.value;
  const isValid = isValidToken(session || "");

  // 4. Redirect to /login if the user is not authenticated
  if (isProtectedRoute && !isValid) {
    const loginUrl = new URL("/login", req.nextUrl);

    loginUrl.searchParams.set("redirect", path); // Add current path to query params

    return NextResponse.redirect(loginUrl);
  }

  // 5. Redirect to /dashboard if the user is authenticated
  if (
    isPublicRoute &&
    isValid &&
    !req.nextUrl.pathname.startsWith("/stock-service")
  ) {
    return NextResponse.redirect(new URL("/stock-service", req.nextUrl));
  }

  return NextResponse.next();
}

// Routes Middleware should not run on
export const config = {
  matcher: ["/((?!api|_next/static|_next/image|.*\\.png$).*)"],
};

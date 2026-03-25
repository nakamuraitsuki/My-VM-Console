import { createBrowserRouter, Navigate } from "react-router"
import type { RouteObject } from "react-router"
import { MainLayout } from "./MainLayout"

/**
 * ルートオブジェクトの拡張タイプ
 * pages ディレクトリの個々で定義して、ここで集約するイメージ感
 */
export type AppRouteMeta = {
  name: string
  requiresAuth?: boolean
}

export type AppRoute = RouteObject & {
  meta?: AppRouteMeta
}

type PageRouteModule = {
  route: AppRoute
}

// 型ガード
const isAppRoute = (value: unknown): value is AppRoute => {
  if (!value || typeof value !== "object") {
    return false
  }

  const route = value as Partial<AppRoute>
  return typeof route.path === "string" && route.path.length > 0 && Boolean(route.element)
}

// cf. https://ja.vite.dev/guide/features.html#glob-%E3%81%AE%E3%82%A4%E3%83%B3%E3%83%9B%E3%82%9A%E3%83%BC%E3%83%88
// pages ディレクトリ以下の page.tsx から route を一括インポートして集約してしまおうというたくらみ…
const pageModules = import.meta.glob<PageRouteModule>("./pages/**/page.tsx", { eager: true })

const appChildRoutes: AppRoute[] = [
  ...Object.values(pageModules)
    .map((mod) => mod.route)
    .filter(isAppRoute),
  {
    path: "*",
    element: <Navigate to="/" replace />,
    meta: { name: "fallback" },
  },
]

export const router = createBrowserRouter([
  {
    path: "/",
    element: <MainLayout />,
    children: appChildRoutes,
  },
])

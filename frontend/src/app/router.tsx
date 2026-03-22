import { createBrowserRouter, Navigate } from "react-router"
import { MainLayout } from "./MainLayout"
import { HomePage, LoginPage, DashboardPage } from "./pages"

export const router = createBrowserRouter([
  {
    path: "/",
    element: <MainLayout />,
    children: [
      {
        path: "/",
        element: <HomePage />,
      },
      {
        path: "/login",
        element: <LoginPage />,
      },
      {
        path: "/dashboard",
        element: <DashboardPage />,
      },
      {
        path: "*",
        element: <Navigate to="/" replace />,
      },
    ]
  },
])

import { createBrowserRouter, Navigate } from "react-router"
import { MainLayout } from "./MainLayout"
import { HomePage, LoginPage, DashboardPage, InstanceDetailPage } from "./pages"

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
        path: "/dashboard/instances/:instanceId",
        element: <InstanceDetailPage />,
      },
      {
        path: "*",
        element: <Navigate to="/" replace />,
      },
    ]
  },
])

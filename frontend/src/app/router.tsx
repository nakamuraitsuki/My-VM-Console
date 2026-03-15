import { createBrowserRouter, Navigate } from "react-router"
import { MainLayout } from "./MainLayout"
import { HomePage, LoginPage } from "./routes"

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
        path: "*",
        element: <Navigate to="/" replace />,
      },
    ]
  },
])

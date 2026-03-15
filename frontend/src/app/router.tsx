import { createBrowserRouter, Navigate } from "react-router"
import { MainLayout } from "./MainLayout"
import { HomePage, LoginPage, MyPage, VideoPlayPage, VideoUploadPage, VideoManagePage } from "./routes"

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
        path: "/my-page",
        element: <MyPage />,
      },
      {
        path: "/upload",
        element: <VideoUploadPage />,
      },
      {
        path: "/video/:videoId",
        element: <VideoPlayPage />,
      },
      {
        path: "/manage/:videoId",
        element: <VideoManagePage />,
      },
      {
        path: "*",
        element: <Navigate to="/" replace />,
      },
    ]
  },
])

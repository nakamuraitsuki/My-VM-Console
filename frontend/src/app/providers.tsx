import { type ReactNode } from "react"
import { ServiceProvider } from "../context/ServiceContext"
import { AuthProvider } from "../context/AuthContext"

type Props = {
  children: ReactNode
}

export function AppProviders({ children }: Props) {
  return (
    <ServiceProvider>
      <AuthProvider>
        {children}
      </AuthProvider>
    </ServiceProvider>
  )
}

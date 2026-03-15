import React from "react";
import { createContext } from "react";
import type { IAuthRepository } from "../domain/auth/auth.repository";
import type { IUserRepository } from "../domain/user/user.repository";
import { AuthRepositoryImpl } from "../gateways/auth/auth.repository.impl";
import { UserRepositoryImpl } from "../gateways/user/user.repository.impl";
import { AuthRepositoryMock } from "../gateways/auth/auth.repository.mock";

interface ServiceContextType {
  authRepo: IAuthRepository;
  userRepo: IUserRepository;
}

const ServiceContext = createContext<ServiceContextType | undefined>(undefined);

export const ServiceProvider = ({ children }: { children: React.ReactNode }) => {
  const isMock = import.meta.env.VITE_USE_MOCK === "true";

  const services: ServiceContextType = {
    authRepo: isMock ? new AuthRepositoryMock() : new AuthRepositoryImpl(),
    userRepo: new UserRepositoryImpl(),
  };

  return (
    <ServiceContext.Provider value={services}>
      {children}
    </ServiceContext.Provider>
  )
};

// helper hook to use the services in custom hooks or components
export const useServices = (): ServiceContextType => {
  const context = React.useContext(ServiceContext);
  if (!context) {
    throw new Error("useServices must be used within a ServiceProvider");
  }
  return context;
}
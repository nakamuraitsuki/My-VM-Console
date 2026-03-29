import React from "react";
import { createContext } from "react";
import type { IAuthRepository } from "../domain/auth/auth.repository";
import { AuthRepositoryImpl } from "../gateways/auth/auth.repository.impl";
import type { IComputeRepository } from "../domain/compute/compute.repository";
import { ComputeGateway } from "../gateways/compute/compute.gateway";
import type { IKeyRepository } from "../domain/compute/key.repository";
import { KeyGateway } from "../gateways/compute/key.gateway";

interface ServiceContextType {
  authRepo: IAuthRepository;
  computeRepository: IComputeRepository;
  keyRepository: IKeyRepository;
}

const ServiceContext = createContext<ServiceContextType | undefined>(undefined);

export const ServiceProvider = ({ children }: { children: React.ReactNode }) => {
  const services: ServiceContextType = {
    authRepo: new AuthRepositoryImpl(),
    computeRepository: new ComputeGateway(),
    keyRepository: new KeyGateway(),
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
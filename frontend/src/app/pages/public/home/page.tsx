import type { AppRoute } from "@/app/router";

export const HomePage = () => {


  return (
    <div style={{ maxWidth: "1200px", margin: "0 auto", padding: "1rem" }}>
      <h2 style={{ textAlign: "center", marginBottom: "1rem" }}>Home</h2>
    </div>
  );
};

export const route: AppRoute = {
  path: "/",
  element: <HomePage />,
  meta: { name: "home" },
};

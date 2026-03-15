import { NavLink, Outlet } from "react-router"
import styles from "./MainLayout.module.css"
import { useAuth } from "../context/AuthContext"

export const MainLayout = () => {
  const { session } = useAuth();
  return (
    <div className={styles.wrapper}>
      <header className={styles.header}>
        <div className={styles.inner}>
          <NavLink to="/" className={styles.logo}>
            Infra Stream
          </NavLink>
          <nav className={styles.nav}>
            <NavLink
              to="/"
              className={({ isActive }) =>
                isActive
                  ? `${styles.navLink} ${styles.active}`
                  : styles.navLink
              }
            >
              Home
            </NavLink>
            {session.status === "authenticated" ? (
              <NavLink
                to="/my-page"
                className={({ isActive }) =>
                  isActive
                    ? `${styles.navLink} ${styles.active}`
                    : styles.navLink
                }
              >
                Library
              </NavLink>
            ) : session.status === "unauthenticated" ? (
              <NavLink
                to="/login"
                className={({ isActive }) =>
                  isActive
                    ? `${styles.navLink} ${styles.active}`
                    : styles.navLink
                }
              >
                Login
              </NavLink>
            ) : (
              <span className={styles.loading}>Loading...</span>
            )}
          </nav>
        </div>
      </header>
      <main>
        <Outlet /> {/* children routes render here */}
      </main>
    </div>
  )
}
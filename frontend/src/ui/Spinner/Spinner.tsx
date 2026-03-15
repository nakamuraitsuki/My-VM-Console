import styles from "./Spinner.module.css";

type SpinnerProps = {
  isLoading: boolean;
}

export const Spinner = (props: SpinnerProps) => {
  if (!props.isLoading) {
    return null;
  }
  return (
    <div className={styles.overlay}>
      <div className={styles.spinner}></div>
    </div>
  )
}
import style from "./IconButton.module.css";

type IconButtonProps = {
  type?: "button" | "submit" | "reset";
  icon: React.ReactNode;
  label?: string;
  onClick?: () => void;
  disabled?: boolean;
  ariaLabel?: string;
}

export const IconButton = ({
  icon,
  label,
  onClick,
  disabled = false,
  ariaLabel,
  type = "button",
}: IconButtonProps) => {
  const buttonStyle =  disabled ? style.buttonDisabled : style.button;
  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled}
      className={buttonStyle}
      aria-label={ariaLabel}
    >
      <span className={style.icon}>{icon}</span>
      {label && <span className={style.label}>{label}</span>}
    </button>
  );
}
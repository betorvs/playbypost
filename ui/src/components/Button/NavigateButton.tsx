import Button from "react-bootstrap/Button";
import { useNavigate } from "react-router-dom";

interface props {
  link: string;
  variant: string;
  children: string;
  disabled?: boolean;
}

const NavigateButton = ({ link, variant, children, disabled }: props) => {
  const navigate = useNavigate();
  const onClick = () => {
    navigate(link);
  };
  return (
    <Button variant={variant} disabled={disabled} onClick={() => onClick()}>
      {children}
    </Button>
  );
};

export default NavigateButton;

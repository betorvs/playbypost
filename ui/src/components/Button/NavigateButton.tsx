import Button from "react-bootstrap/Button";
import { useNavigate } from "react-router-dom";
import CheckSession from "../../context/CheckSession";

interface props {
  link: string;
  variant: string;
  children: string;
  disabled?: boolean;
}

const NavigateButton = ({ link, variant, children, disabled }: props) => {
  const navigate = useNavigate();
  const onClick = () => {
    if (CheckSession() === true) {
      navigate(link);
    } else {
      console.log("Navigating to: /login");
      navigate("/login");
    }
    
  };
  return (
    <Button variant={variant} disabled={disabled} onClick={() => onClick()}>
      {children}
    </Button>
  );
};

export default NavigateButton;

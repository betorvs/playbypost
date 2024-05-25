import Button from "react-bootstrap/Button";
import { useNavigate } from "react-router-dom";

interface props {
  link: string;
  variant: string;
  children: string;
}

const NavigateButton = ({ link, variant, children }: props) => {
  const navigate = useNavigate();
  const onClick = () => {
    navigate(link);
  };
  return (
    <Button variant={variant} onClick={() => onClick()}>
      {children}
    </Button>
  );
};

export default NavigateButton;

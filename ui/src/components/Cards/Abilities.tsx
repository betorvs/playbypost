
interface props {
    key: String;
    value: String
  }


const Ability = ({ key, value }: props) => {
    return (
      <>
      <ul>{key}: {value}</ul>
      </>
  );
};

export default Ability;
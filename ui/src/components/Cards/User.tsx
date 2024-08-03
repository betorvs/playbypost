import UsersCard from "../../types/UserCard";
import NavigateButton from "../Button/NavigateButton";

interface props {
  user: UsersCard;
}

const UserCards = ({ user }: props) => {
  return (
    <>
      <div className="col-md-6">
        <div className="card mb-4">
          <div className="card-body">
            <h5 className="card-title">Username: {user.username}</h5>
          </div>
          <ul className="list-group list-group-flush">
            <li className="list-group-item">UserID: {user.user_id}</li>
            <li className="list-group-item">Channel: {user.channel}</li>
          </ul>
          <div className="card-footer">
          <NavigateButton link={`/users/${user.user_id}`} variant="primary">
            Add as Storyteller
          </NavigateButton>{" "}
          <NavigateButton link={`/users/player/${user.user_id}`} variant="primary">
            Add as Player
          </NavigateButton>{" "}
          </div>
        </div>
      </div>
    </>
  );
};

export default UserCards;

import UsersCard from "../../types/UserCard";

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
            <li className="list-group-item">
              Stories: {user.stories.join(",") || " no stories created"}{" "}
            </li>
            <li className="list-group-item">
              Players: {user.players.join(",") || " no players created"}{" "}
            </li>
          </ul>
          <div className="card-footer">
            {/* <a href="#" className="btn btn-primary">
              Invite
            </a> */}
          </div>
        </div>
      </div>
    </>
  );
};

export default UserCards;

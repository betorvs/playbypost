import Task from "../../types/Task";

interface props {
  task: Task;
}

const TaskCard = ({ task }: props) => {
  return (
    <>
      <div className="col-md-6">
        <div className="card mb-4">
          <div className="card-body">
            <h5 className="card-title">ID: {task.id}</h5>
          </div>
          <ul className="list-group list-group-flush">
            <li className="list-group-item">Description: {task.description}</li>
            <li className="list-group-item">Ability: {task.ability}</li>
            <li className="list-group-item">Skill: {task.skill}</li>
          </ul>
          <div className="card-footer">
            <p>Target: {task.target}</p>
          </div>
        </div>
      </div>
    </>
  );
};

export default TaskCard;
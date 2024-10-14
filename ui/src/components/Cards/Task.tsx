import Task from "../../types/Task";
import { useTranslation } from "react-i18next";

interface props {
  task: Task;
}

const TaskCard = ({ task }: props) => {
  const { t } = useTranslation(['home', 'main']);
  return (
    <>
      <div className="col-md-6">
        <div className="card mb-4">
          <div className="card-body">
            <h5 className="card-title">ID: {task.id}</h5>
          </div>
          <ul className="list-group list-group-flush">
            <li className="list-group-item">{t("common.description", {ns: ['main', 'home']})}: {task.description}</li>
            <li className="list-group-item">{t("common.ability", {ns: ['main', 'home']})}: {task.ability}</li>
            <li className="list-group-item">{t("common.skill", {ns: ['main', 'home']})}: {task.skill}</li>
          </ul>
          <div className="card-footer">
            <p>{t("common.target", {ns: ['main', 'home']})}: {task.target}</p>
          </div>
        </div>
      </div>
    </>
  );
};

export default TaskCard;
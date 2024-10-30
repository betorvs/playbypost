import Task from "../../types/Task";
import { useTranslation } from "react-i18next";
import NavigateButton from "../Button/NavigateButton";
import { Button } from "react-bootstrap";
import { DeleteTaskByID } from "../../functions/Tasks";

interface props {
  task: Task;
}

const TaskCard = ({ task }: props) => {
  const { t } = useTranslation(['home', 'main']);

  const handleDelete = (id: number) => {
    console.log("Deleting task " + id);
    DeleteTaskByID(id);
  }


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
            <li className="list-group-item">{t("common.target", {ns: ['main', 'home']})}: {task.target}</li>
          </ul>
          <div className="card-footer">
            <NavigateButton link={`/tasks/${task.id}/edit`} variant="warning">
              {t("common.edit", {ns: ['main', 'home']})}
            </NavigateButton>{" "}
            <Button variant="danger" size="sm" onClick={() => handleDelete(task.id)}>{t("common.delete", {ns: ['main', 'home']})}</Button>
          </div>
        </div>
      </div>
    </>
  );
};

export default TaskCard;
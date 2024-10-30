import { useContext, useEffect, useState } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import { Form } from "react-bootstrap";
import Button from "react-bootstrap/Button";
import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import { useNavigate, useParams } from "react-router-dom";
import UseLocation from "../context/UseLocation";
import { useTranslation } from "react-i18next";
import Task from "../types/Task";
import { FetchTask } from "../functions/Tasks";

const EditTask = () => {
  const { Logoff } = useContext(AuthContext);
  const { id } = useParams();
  const [task, setTask] = useState<Task>();
  const navigate = useNavigate();
  const { t } = useTranslation(["home", "main"]);

  const safeID: string = id ?? "";

  const cancelButton = () => {
    navigate("/tasks");
  };

  const onChangeDescription = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (task) {
      setTask({ ...task, description: e.target.value });
    }
  }

  const onChangeAbility = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (task) {
      setTask({ ...task, ability: e.target.value });
    }
  }

  const onChangeSkill = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (task) {
      setTask({ ...task, skill: e.target.value });
    }
  }

  const onChangeTarget = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (task) {
      setTask({ ...task, target: Number(e.target.value) });
    }
  }

  useEffect(() => {
    FetchTask(safeID, setTask);
  }, []);

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const apiURL = UseLocation();
    const urlAPI = new URL("api/v1/task/" + id, apiURL);
    const response = await fetch(urlAPI, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "X-Username": GetUsername(),
        "X-Access-Token": GetToken(),
      },
      body: JSON.stringify({
        id: Number(id),
        description: task?.description,
        ability: task?.ability,
        skill: task?.skill,
        target: task?.target,
      }),
    });
    if (response.ok) {
      alert(t("alert.task-edit", {ns: ['main', 'home']}));
      navigate("/tasks");
    } else {
      alert(t("alert.task-edit-wrong", {ns: ['main', 'home']}));
    }
  }
  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        <h2>{t("task.edit-header", {ns: ['main', 'home']})}</h2>
        <hr />
      </div>
      <div className="container mt-3" key="2">
        <Form onSubmit={handleSubmit}>
          <Form.Group className="mb-3" controlId="formTitle">
            <Form.Label>{t("common.description", {ns: ['main', 'home']})}</Form.Label>
            <Form.Control
              type="text"
              placeholder="title"
              value={task?.description}
              onChange={onChangeDescription}
            />
            <Form.Text className="text-muted">
            {t("task.new-task-text1", {ns: ['main', 'home']})}
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formAbility">
            <Form.Label>{t("common.ability", {ns: ['main', 'home']})}</Form.Label>
            <Form.Control
              type="text"
              placeholder="ability"
              value={task?.ability}
              onChange={onChangeAbility}
            />
            <Form.Text className="text-muted">
            {t("task.new-task-text2", {ns: ['main', 'home']})}
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formSkill">
            <Form.Label>{t("common.skill", {ns: ['main', 'home']})}</Form.Label>
            <Form.Control
              type="text"
              placeholder="skill"
              value={task?.skill}
              onChange={onChangeSkill}
            />
            <Form.Text className="text-muted">
            {t("task.new-task-text3", {ns: ['main', 'home']})}
            </Form.Text>
          </Form.Group>
          <Form.Group className="mb-3" controlId="formTarget">
            <Form.Label>{t("common.target", {ns: ['main', 'home']})}</Form.Label>
            <Form.Control
              type="number"
              placeholder="target"
              value={task?.target}
              onChange={onChangeTarget}
            />
            <Form.Text className="text-muted">
            {t("task.new-task-text4", {ns: ['main', 'home']})}
            </Form.Text>
          </Form.Group>
          <Button variant="primary" type="submit">
          {t("common.submit", {ns: ['main', 'home']})}
          </Button>{" "}
          <Button variant="secondary" onClick={() => cancelButton()}>
          {t("common.cancel", {ns: ['main', 'home']})}
          </Button>{" "}
        </Form>
      </div>
    </>
  );
};

export default EditTask;

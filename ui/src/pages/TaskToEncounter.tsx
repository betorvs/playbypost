import { useContext, useEffect, useState } from "react";
import { AuthContext } from "../context/AuthContext";
import { useNavigate, useParams } from "react-router-dom";
import { Form } from "react-bootstrap";
import Layout from "../components/Layout";
import { Button } from "react-bootstrap";
import Task from "../types/Task";
import FetchTasks from "../functions/Tasks";
import UseLocation from "../context/UseLocation";
import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";

const TaskToEncounter = () => {
  const { Logoff } = useContext(AuthContext);
  const { id, story, encounterid, storyteller_id } = useParams();
  const [tasks, setTask] = useState<Task[]>([]);
  const [taskID, setTaskID] = useState(0);
  const [text, setText] = useState("");

  const navigate = useNavigate();

  useEffect(() => {
    FetchTasks(setTask);
  }, []);

  const cancelButton = () => {
    navigate(`/stages/${id}/story/${story}/encounter/${encounterid}`);
  };

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const apiURL = UseLocation();
        const urlAPI = new URL("api/v1/stage/encounter/task", apiURL);
        const response = await fetch(urlAPI, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "X-Username": GetUsername(),
            "X-Access-Token": GetToken(),
          },
          body: JSON.stringify({
            stage_id: Number(id),
            task_id: taskID,
            storyteller_id: Number(storyteller_id),
            encounter_id: Number(encounterid),
            text: text,
          }),
        });
        if (response.ok) {
          alert("Task added! Have a great session with your friends!");
          navigate(`/stages/${id}/story/${story}/encounter/${encounterid}`);
        } else {
          alert("Something goes wrong. Task was not added to encounter.");
        }
  };
  
  return (
      <>
        <div className="container mt-3" key="1">
            <Layout Logoff={Logoff} />
            <h2>Add Task to Encounter</h2>
        </div>
        <div className="container mt-3" key="2">
            <Form onSubmit={handleSubmit}>
                <Form.Group className="mb-3" controlId="formTaskEncounter">
                    <Form.Label>Task</Form.Label>
                    <Form.Select as="select"
                      value={taskID}
                      onChange={e => {
                        console.log("set e.target.value", e.target.value);
                        setTaskID(Number(e.target.value));
                      }}>
                        <option value="-1">Select a Task</option>
                      {
                        tasks != null ? (
                            tasks.map((task) => (
                            <option value={task.id}>{task.description}</option>
                          ))) : (
                            <option>tasks not found</option>
                          )
                      }
                    </Form.Select>
                    <Form.Text className="text-muted">
                      Choose a Task to be used in the encounter
                    </Form.Text>
                </Form.Group>
                <Form.Group className="mb-3" controlId="formName">
                    <Form.Label>Text</Form.Label>
                    <Form.Control
                      type="text"
                      placeholder="name"
                      value={text}
                      onChange={(e) => setText(e.target.value)}
                    />
                    <Form.Text className="text-muted">
                    Text that will be diplaid to the players
                    </Form.Text>
                </Form.Group>
                
                <Button variant="primary" type="submit">
                Submit
                </Button>{" "}
                <Button variant="secondary" onClick={() => cancelButton()}>
                Cancel
                </Button>{" "}
            </Form>
        </div>
      </>
  );
};
  
export default TaskToEncounter;
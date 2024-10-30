import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import UseLocation from "../context/UseLocation";
import Task from "../types/Task";
import CleanSession from "../context/CleanSession";

const FetchTasks = async (
  setTask: React.Dispatch<React.SetStateAction<Task[]>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/task", apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    console.log(data);
    setTask(JSON.parse(data));
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
  }
};

const FetchTask = async (
  id: string,
  setTask: React.Dispatch<React.SetStateAction<Task | undefined>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/task/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    setTask(JSON.parse(data));
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
  }
}

const DeleteTaskByID = async (id: number): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/task/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "DELETE",
    headers: requestHeaders,
  });
  if (response.ok) {
    console.log("Task deleted");
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
  }
}

export default FetchTasks;
export { FetchTask, DeleteTaskByID };
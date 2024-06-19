import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import UseLocation from "../context/UseLocation";
import UsersCard from "../types/UserCard";

const FetchTasks = async (
  setUser: React.Dispatch<React.SetStateAction<UsersCard[]>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/user/card", apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    console.log(data);
    setUser(JSON.parse(data));
  }
};
export default FetchTasks;
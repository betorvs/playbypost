import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import Story from "../types/Story";
import UseLocation from "../context/UseLocation";
import CleanSession from "../context/CleanSession";

const FetchStoriesByUserID = async (
  userID: number,
  setStory: React.Dispatch<React.SetStateAction<Story[]>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/story/writer/" + userID, apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    setStory(JSON.parse(data));
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
  }
};

const FetchStory = async (
  id: string,
  setStory: React.Dispatch<React.SetStateAction<Story | undefined>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/story/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    setStory(JSON.parse(data));
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
  }
};

export default FetchStory;
export { FetchStory, FetchStoriesByUserID };

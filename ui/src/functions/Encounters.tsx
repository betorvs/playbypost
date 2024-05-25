import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import Encounter from "../types/Encounter";
import UseLocation from "../context/UseLocation";

const FetchEncounters = async (
  id: string,
  setEncounters: React.Dispatch<React.SetStateAction<Encounter[]>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/encounter/story/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    setEncounters(JSON.parse(data));
  }
};

export default FetchEncounters;

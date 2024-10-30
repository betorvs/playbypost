import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import Encounter from "../types/Encounter";
import UseLocation from "../context/UseLocation";
import CleanSession from "../context/CleanSession";

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
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
  }
};

const FetchEncounter = async (
  id: string,
  setEncounter: React.Dispatch<React.SetStateAction<Encounter | undefined>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/encounter/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    setEncounter(JSON.parse(data));
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
  }
};

const DeleteEncounterByID = async (id: number): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/encounter/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "DELETE",
    headers: requestHeaders,
  });
  if (response.ok) {
    console.log("Encounter deleted");
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
  }
}

export default FetchEncounters;
export { FetchEncounter, DeleteEncounterByID };

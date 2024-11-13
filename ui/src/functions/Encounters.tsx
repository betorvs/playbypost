import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import Encounter from "../types/Encounter";
import UseLocation from "../context/UseLocation";
import CleanSession from "../context/CleanSession";

const FetchEncounters = async (
  id: string,
  cursor: string,
  encounters: Encounter[],
  setLoading: React.Dispatch<React.SetStateAction<boolean>>,
  setEncounters: React.Dispatch<React.SetStateAction<Encounter[]>>,
  setCursor: React.Dispatch<React.SetStateAction<string>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/encounter/story/" + id, apiURL);
  urlAPI.searchParams.append("limit", "10");
  if (cursor !== "") {
    urlAPI.searchParams.append("cursor", cursor);
  }
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    
    setEncounters([...encounters, ...JSON.parse(data)]);
    const header = response.headers.get("X-Cursor");
    console.log("Header Cursor: " + header);
    if (header !== null) {
      setCursor(header);
    } else {
      setLoading(true);
    }
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

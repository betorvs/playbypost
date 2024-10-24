import GetToken from "../context/GetToken";
import GetUsername from "../context/GetUsername";
import UseLocation from "../context/UseLocation";
import AutoPlay from "../types/AutoPlay";
import { EncounterList } from "../types/Next";



const FetchAutoPlay = async (
  setAutoPlay: React.Dispatch<React.SetStateAction<AutoPlay[]>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/autoplay", apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    setAutoPlay(JSON.parse(data));
  }
};

const FetchAutoPlayByID = async (
  id: string,
  setAutoPlay: React.Dispatch<React.SetStateAction<AutoPlay | undefined>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/autoplay/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    setAutoPlay(JSON.parse(data));
  }
};

const FetchEncounterListAutoPlay = async (
  id: string,
  setEncounters: React.Dispatch<React.SetStateAction<EncounterList | undefined>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/autoplay/encounter/story/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    setEncounters(JSON.parse(data));
  }
};

const DeleteAutoPlayNextEncounter = async (
  id: number,
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/autoplay/next/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "DELETE",
    headers: requestHeaders,
  });
  if (response.ok) {
    console.log("Deleted");
  }
}

export default FetchAutoPlay;
export { FetchAutoPlayByID, FetchEncounterListAutoPlay, DeleteAutoPlayNextEncounter };
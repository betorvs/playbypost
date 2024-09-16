import GetToken from "../context/GetToken";
import GetUsername from "../context/GetUsername";
import UseLocation from "../context/UseLocation";
import AutoPlay, { AutoPlayEncounterList } from "../types/AutoPlay";



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

const FetchEncountersAutoPlay = async (
  id: string,
  setEncounters: React.Dispatch<React.SetStateAction<AutoPlayEncounterList | undefined>>
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

export default FetchAutoPlay;
export { FetchAutoPlayByID, FetchEncountersAutoPlay };
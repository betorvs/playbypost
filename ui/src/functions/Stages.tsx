import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import Stage from "../types/Stage";
import UseLocation from "../context/UseLocation";
import StageAggregated from "../types/StageAggregated";
import Encounter from "../types/Encounter";
import { EncounterList } from "../types/Next";


const FetchStages = async (
  setStage: React.Dispatch<React.SetStateAction<Stage[]>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/stage", apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    setStage(JSON.parse(data));
  }
};

const FetchStage = async (
  id: string,
  setStage: React.Dispatch<React.SetStateAction<StageAggregated | undefined>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/stage/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    setStage(JSON.parse(data));
  }
};

const FetchStageByStoryID = async (
  id: string,
  setStage: React.Dispatch<React.SetStateAction<Stage[]>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/stage/story/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    setStage(JSON.parse(data));
  }
};

const FetchStageEncountersByID = async (
  id: string,
  setEncounters: React.Dispatch<React.SetStateAction<Encounter[]>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/stage/encounters/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    setEncounters(JSON.parse(data));
  }
};

const FetchStageEncounterByEncounterID = async (
  id: string,
  setEncounter: React.Dispatch<React.SetStateAction<Encounter | undefined>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/stage/encounter/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    setEncounter(JSON.parse(data));
  }
};

const FetchEncounterListStage = async (
  id: string,
  setEncounterList: React.Dispatch<React.SetStateAction<EncounterList | undefined>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/stage/encounter/story/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    setEncounterList(JSON.parse(data));
  }
};

const DeleteStageNextEncounter = async (
  id: number,
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/stage/encounter/next/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "DELETE",
    headers: requestHeaders,
  });
  if (response.ok) {
    console.log("Stage deleted");
  }
}

const CloseStage = async (
  id: number,
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/stage/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "PUT",
    headers: requestHeaders,
  });
  if (response.ok) {
    console.log("Stage closed");
  }
}

export default FetchStages;
export { FetchStage, FetchStageByStoryID, FetchStageEncountersByID, FetchStageEncounterByEncounterID, FetchEncounterListStage,DeleteStageNextEncounter, CloseStage };
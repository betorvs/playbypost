import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import Stage from "../types/Stage";
import UseLocation from "../context/UseLocation";
import StageAggregated from "../types/StageAggregated";
import Encounter from "../types/Encounter";
import { EncounterList } from "../types/Next";
import CleanSession from "../context/CleanSession";


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
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
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
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
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
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
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
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
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
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
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
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
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
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
  }
}

const DeleteStageEncounterByID = async (
  id: number,
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/stage/encounter/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "DELETE",
    headers: requestHeaders,
  });
  if (response.ok) {
    console.log("Stage Encounter deleted");
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
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
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
  }
}

export default FetchStages;
export { FetchStage, FetchStageByStoryID, FetchStageEncountersByID, FetchStageEncounterByEncounterID, FetchEncounterListStage,DeleteStageNextEncounter, DeleteStageEncounterByID, CloseStage };
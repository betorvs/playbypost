import { Session, SessionEvent } from '../types/Session';
import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import UseLocation from "../context/UseLocation";

export const fetchAllSessions = async (): Promise<Session[]> => {
  const requestHeaders: HeadersInit = new Headers();
    requestHeaders.set("Content-Type", "application/json");
    requestHeaders.set("X-Username", GetUsername());
    requestHeaders.set("X-Access-Token", GetToken());
    const apiURL = UseLocation();
    const urlAPI = new URL("api/v1/session", apiURL);
    const response = await fetch(urlAPI, {
      method: "GET",
      headers: requestHeaders,
    });
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  // filter response by username
  const data = await response.json();
  return data.filter((session: Session) => session.Username === GetUsername());
};

export const fetchSessionEvents = async (): Promise<SessionEvent[]> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/session/events", apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  return response.json();
};

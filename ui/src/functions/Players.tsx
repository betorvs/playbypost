import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import Players from "../types/Players";
import UseLocation from "../context/UseLocation";
import CleanSession from "../context/CleanSession";

const FetchPlayers = async (
  id: string,
  setPlayer: React.Dispatch<React.SetStateAction<Players[]>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/stage/player/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    setPlayer(JSON.parse(data));
  } else if (response.status === 403) {
    console.log("Not authorized");
    CleanSession();
  }
};
export default FetchPlayers;

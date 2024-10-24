import GetToken from "../context/GetToken";
import GetUsername from "../context/GetUsername";
import UseLocation from "../context/UseLocation";
import Activities from "../types/Activities";

const FetchActivities = async (
    id: string,
    setActivities: React.Dispatch<React.SetStateAction<Activities[]>>
) => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/stage/encounter/activities/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    setActivities(JSON.parse(data));
  }
};
export default FetchActivities;
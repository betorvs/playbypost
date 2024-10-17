import GetToken from "../context/GetToken";
import GetUsername from "../context/GetUsername";
import UseLocation from "../context/UseLocation";
import Validator from "../types/validator";


const ValidatorPut = async (
  id: number,
  kind: string,
//   allowed kind values: stage, story, autoplay
): Promise<void> => {
  if (kind !== "stage" && kind !== "story" && kind !== "autoplay") {
    throw new Error("Invalid kind value");
  }
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/validator/"+ kind + "/" + id, apiURL);
  const response = await fetch(urlAPI, {
    method: "PUT",
    headers: requestHeaders,
  });
  if (response.ok) {
    console.log("Request to validator done");
  }
}

const FetchValidator = async (
  setValidator: React.Dispatch<React.SetStateAction<Validator[]>>
): Promise<void> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL("api/v1/validator", apiURL);
  const response = await fetch(urlAPI, {
    method: "GET",
    headers: requestHeaders,
  });
  if (response.ok) {
    const data = await response.text();
    setValidator(JSON.parse(data));
  }
};

const FetchValidatorByIDKind = async (
    id: number,
    kind: string,
    setValidator: React.Dispatch<React.SetStateAction<Validator | undefined>>,
  ): Promise<void> => {
    const requestHeaders: HeadersInit = new Headers();
    requestHeaders.set("Content-Type", "application/json");
    requestHeaders.set("X-Username", GetUsername());
    requestHeaders.set("X-Access-Token", GetToken());
    const apiURL = UseLocation();
    const params = new URLSearchParams(
        'output=json',
    );
    const urlAPI = new URL("api/v1/validator/" + kind + "/" + id, apiURL);
    urlAPI.search = params.toString();
    const response = await fetch(urlAPI, {
      method: "GET",
      headers: requestHeaders,
    });
    if (response.ok) {
      const data = await response.text();
      setValidator(JSON.parse(data));
    }
  };


export { ValidatorPut, FetchValidatorByIDKind };
export default FetchValidator;
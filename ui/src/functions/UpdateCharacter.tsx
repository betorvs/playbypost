import GetUsername from "../context/GetUsername";
import GetToken from "../context/GetToken";
import UseLocation from "../context/UseLocation";
import CleanSession from "../context/CleanSession";
import Players from "../types/Players";

const UpdateCharacter = async (
  characterId: number,
  characterData: Partial<Players>
): Promise<boolean> => {
  const requestHeaders: HeadersInit = new Headers();
  requestHeaders.set("Content-Type", "application/json");
  requestHeaders.set("X-Username", GetUsername());
  requestHeaders.set("X-Access-Token", GetToken());
  const apiURL = UseLocation();
  const urlAPI = new URL(`api/v1/characters/${characterId}`, apiURL);

  try {
    const response = await fetch(urlAPI, {
      method: "PUT",
      headers: requestHeaders,
      body: JSON.stringify(characterData),
    });

    if (response.ok) {
      console.log(`Character ${characterId} updated successfully.`);
      return true;
    } else if (response.status === 403) {
      console.error("Not authorized to update character.");
      CleanSession();
      return false;
    } else {
      const errorData = await response.text();
      console.error(`Failed to update character ${characterId}: ${errorData}`);
      return false;
    }
  } catch (error) {
    console.error("Network error or other issue during character update:", error);
    return false;
  }
};

export default UpdateCharacter;

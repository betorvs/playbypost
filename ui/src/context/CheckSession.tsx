import CleanSession from "./CleanSession";

function CheckSession(): boolean {
  var expire = sessionStorage.getItem("expire") || "";
  const inputDate = new Date(expire);

  // Get the current date
  const currentDate = new Date();
  if (inputDate < currentDate) {
    CleanSession();
    return false;
  } else if (currentDate < inputDate) {
    // force to verify and return true
    return true;
  }
  //   always return false in case of any issue above
  return false;
}

export default CheckSession;

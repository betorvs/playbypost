// import { useLocation } from "react-router-dom";

// function UseLocation(): string {
//   const apiURL = useLocation();
//   return apiURL.state;
// }

function UseLocation(): string {
  let urlCur: string = window.location.origin;
  // return "http://192.168.1.210:3000/";
  if (window.location.origin === "http://localhost:5173") {
    urlCur = "http://localhost:3000"
  }
  return urlCur;
}

export default UseLocation;

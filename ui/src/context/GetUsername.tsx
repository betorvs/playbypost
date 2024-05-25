function GetUsername(): string {
  let token = sessionStorage.getItem("username") || "";
  return token;
}
export default GetUsername;

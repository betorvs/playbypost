function GetToken(): string {
  let token = sessionStorage.getItem("token") || "";
  return token;
}
export default GetToken;

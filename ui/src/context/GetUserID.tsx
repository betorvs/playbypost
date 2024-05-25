function GetUserID(): number {
  let user_id = sessionStorage.getItem("user_id") || "";
  if (user_id === "") {
    return -1;
  }
  return Number(user_id);
}
export default GetUserID;
